// Copyright 2015-2016 Apcera Inc. All rights reserved.

package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/apcera/kurma/pkg/apiclient"
	"github.com/appc/spec/schema/types"
)

type ImageService struct {
	server *Server
}

func (s *Server) imageCreateRequest(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	hash, manifest, err := s.options.ImageManager.CreateImage(req.Body)
	if err != nil {
		s.log.Errorf("Failed create image: %v", err)
		http.Error(w, "Failed to create image", 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := &apiclient.ImageResponse{Image: &apiclient.Image{Hash: hash, Manifest: manifest}}
	json.NewEncoder(w).Encode(resp)
}

// imageFetchRequest is a handler for requests instructing the daemon to fetch
// and create a particular image.
func (s *Server) imageFetchRequest(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// TODO: should clients be able to specify insecure and/or labels? Or are they config-time?
	var imageFetchRequest *apiclient.ImageFetchRequest

	if err := json.NewDecoder(req.Body).Decode(&imageFetchRequest); err != nil {
		s.log.Errorf("Failed to unmarshal request body: %s", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	if imageFetchRequest.FetchConfig.ACILabels == nil {
		imageFetchRequest.FetchConfig.ACILabels = make(map[types.ACIdentifier]string)
	}

	if imageFetchRequest.FetchConfig.ACILabels[types.ACIdentifier("os")] == "" {
		imageFetchRequest.FetchConfig.ACILabels[types.ACIdentifier("os")] = runtime.GOOS
	}

	if imageFetchRequest.FetchConfig.ACILabels[types.ACIdentifier("arch")] == "" {
		imageFetchRequest.FetchConfig.ACILabels[types.ACIdentifier("arch")] = runtime.GOARCH
	}

	hash, manifest, err := s.options.ImageManager.FetchImage(imageFetchRequest.ImageURI)
	if err != nil {
		s.log.Errorf("Failed create image: %s", err)
		http.Error(w, "Failed to create image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := &apiclient.ImageResponse{Image: &apiclient.Image{Hash: hash, Manifest: manifest}}
	json.NewEncoder(w).Encode(resp)
}

func (s *ImageService) List(r *http.Request, args *apiclient.None, resp *apiclient.ImageListResponse) error {
	images := s.server.options.ImageManager.ListImages()
	resp.Images = make([]*apiclient.Image, 0, len(images))
	for hash, image := range images {
		imageSize, err := s.server.options.ImageManager.GetImageSize(hash)
		if err != nil {
			s.server.log.Warnf("Failed to get image size %s: %v", hash, err)
			continue
		}
		resp.Images = append(resp.Images, &apiclient.Image{Hash: hash, Manifest: image, Size: imageSize})
	}
	return nil
}

func (s *ImageService) Get(r *http.Request, hash *string, resp *apiclient.ImageResponse) error {
	if hash == nil {
		return fmt.Errorf("no image hash was specified")
	}
	image := s.server.options.ImageManager.GetImage(*hash)
	if image == nil {
		return fmt.Errorf("specified image not found")
	}
	imageSize, err := s.server.options.ImageManager.GetImageSize(*hash)
	if err != nil {
		return err
	}
	resp.Image = &apiclient.Image{Hash: *hash, Manifest: image, Size: imageSize}
	return nil
}

func (s *ImageService) Delete(r *http.Request, hash *string, resp *apiclient.ImageResponse) error {
	if hash == nil {
		return fmt.Errorf("no image hash was specified")
	}
	return s.server.options.ImageManager.DeleteImage(*hash)
}
