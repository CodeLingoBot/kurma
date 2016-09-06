// Copyright 2015-2016 Apcera Inc. All rights reserved.

package docker

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/apcera/kurma/pkg/remote"

	"github.com/apcera/util/docker"

	docker2aci "github.com/appc/docker2aci/lib"
	docker2acicommon "github.com/appc/docker2aci/lib/common"
)

// A dockerPuller represents allows for fetching Docker images to run as Kurma
// containers.
type dockerPuller struct {
	// insecure, if true, will pull a Docker in an insecure manner, skipping
	// signature verification.
	insecure bool

	// convertToACI, if true, will convert Docker images into App Container
	// Images.
	convertToACI bool

	// squashLayers, if true, will squash together all of the layers of the
	// Docker image into a single tarball.
	squashLayers bool
}

// New creates a new dockerPuller to pull a remote Docker image.
func New(insecure bool) remote.Puller {
	puller := &dockerPuller{
		insecure:     insecure,
		convertToACI: true, // FIXME: should we support pull without conversion?
		squashLayers: true,
	}
	return puller
}

// Pull fetches a remote Docker image.
func (d *dockerPuller) Pull(dockerImageURI string) ([]io.ReadCloser, error) {
	if !strings.HasPrefix(dockerImageURI, "docker://") {
		return nil, fmt.Errorf("only 'docker://' scheme image URLs supported")
	}

	dockerURL, err := docker.ParseDockerRegistryURL(dockerImageURI)
	if err != nil {
		return nil, err
	}

	tmpdir, err := ioutil.TempDir(os.TempDir(), "docker2aci")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp path to handle Docker image conversion: %s", err)
	}
	defer os.RemoveAll(tmpdir)

	schemelessURL := strings.TrimPrefix(dockerURL.String(), "docker://")

	config := docker2aci.RemoteConfig{
		CommonConfig: docker2aci.CommonConfig{
			Squash:      d.squashLayers,
			OutputDir:   tmpdir, // FIXME(alex): this should go to the image path on disk
			TmpDir:      tmpdir,
			Compression: docker2acicommon.NoCompression,
		},
		Insecure: d.insecure,
	}

	if d.convertToACI {
		return d.pullAsACI(schemelessURL, config)
	}
	return nil, errors.New("only ACI conversion pull supported")
}

// pullAsACI fetches a Docker image and converts it into an ACI. Callers are
// responsible for closing the Reader(s).
func (d *dockerPuller) pullAsACI(dockerURL string, config docker2aci.RemoteConfig) ([]io.ReadCloser, error) {
	if !d.convertToACI {
		return nil, errors.New("not configured to convert to ACI")
	}

	acis, err := docker2aci.ConvertRemoteRepo(dockerURL, config)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Docker image: %s", err)
	}

	if d.squashLayers && len(acis) != 1 {
		return nil, fmt.Errorf("fetched %d layer(s), expected 1", len(acis))
	}

	files := make([]io.ReadCloser, len(acis))

	for j, aci := range acis {
		f, err := os.Open(aci)
		if err != nil {
			return nil, fmt.Errorf("failed to open converted Docker image %q: %s", aci, err)
		}
		files[j] = f
	}

	return files, nil
}
