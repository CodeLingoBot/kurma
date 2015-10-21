// Copyright 2015 Apcera Inc. All rights reserved.

package api

import (
	"io"

	pb "github.com/apcera/kurma/stage1/client"
)

func (s *rpcServer) Enter(inStream pb.Kurma_EnterServer) error {
	s.log.Debug("Received enter request")

	// read the first chunk to make sure its real and get the container ID
	chunk, err := inStream.Recv()
	if err != nil {
		return err
	}

	// create the outbound stream
	outStream, err := s.client.Enter(inStream.Context())
	if err != nil {
		return err
	}

	// write the first byte to the backend so it is initialized
	if err := outStream.Send(chunk); err != nil {
		return err
	}

	// Create our inbound streams
	inWriter := pb.NewByteStreamWriter(inStream, chunk.StreamId)
	inReader := pb.NewByteStreamReader(pb.NewEnterRequestBrokerReader(inStream), nil)

	// Create our outbound streams
	outWriter := pb.NewByteStreamWriter(pb.NewEnterRequestBrokerWriter(outStream), chunk.StreamId)
	outReader := pb.NewByteStreamReader(outStream, nil)

	// stream between!
	go io.Copy(outWriter, inReader)
	io.Copy(inWriter, outReader)
	outStream.CloseSend()
	return nil
}
