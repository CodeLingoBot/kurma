// Copyright 2015-2016 Apcera Inc. All rights reserved.

package image

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/apcera/kurma/pkg/backend"
	"github.com/apcera/kurma/pkg/remote"
	"github.com/apcera/kurma/pkg/remote/aci"
	"github.com/apcera/kurma/pkg/remote/docker"
	"github.com/apcera/kurma/pkg/remote/http"

	"github.com/apcera/util/tempfile"

	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

// A FetchConfig contains configuration for an image fetch operation.
type FetchConfig struct {
	// ACILabels are labels that scope the image resolution request.
	ACILabels map[types.ACIdentifier]string `json:"aci_labels,omitempty"`

	// Insecure is an option that, if enabled, will fetch images insecurely.
	// For instance, this will disable signature verification, and will use
	// HTTP for fetching images rather than HTTPS, where applicable.
	Insecure bool `json:"insecure"`
}

// FetchAndLoad retrieves a container image and loads it for use within kurmad.
func (f *FetchConfig) FetchAndLoad(imageURI string, imageManager backend.ImageManager) (
	string, *schema.ImageManifest, error) {

	layers, err := f.Fetch(imageURI)
	if err != nil {
		return "", nil, err
	}
	for _, l := range layers {
		defer l.Close()
	}

	hash, manifest, err := loadFromFile(layers[0], imageManager)
	if err != nil {
		return "", nil, err
	}
	return hash, manifest, nil
}

// Fetch retrieves a container image. Images may be sourced from the local
// machine, or may be retrieved from a remote server.
func (f *FetchConfig) Fetch(imageURI string) ([]tempfile.ReadSeekCloser, error) {
	u, err := url.Parse(imageURI)
	if err != nil {
		return nil, err
	}

	var puller remote.Puller

	switch u.Scheme {
	case "file":
		filename := u.Path
		if u.Host != "" {
			filename = filepath.Join(u.Host, u.Path)
		}
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		t, err := tempfile.New(f)
		if err != nil {
			return nil, err
		}
		return []tempfile.ReadSeekCloser{t}, nil
	case "http", "https":
		puller = http.New()
	case "docker":
		puller = docker.New(f.Insecure)
	case "aci", "":
		puller = aci.New(f.Insecure, f.ACILabels)
	default:
		return nil, fmt.Errorf("%q scheme not supported", u.Scheme)
	}

	layers, err := puller.Pull(imageURI)
	if err != nil {
		return nil, err
	}

	wrappedLayers := make([]tempfile.ReadSeekCloser, len(layers))
	for j, layer := range layers {
		t, err := tempfile.New(layer)
		if err != nil {
			return nil, err
		}
		wrappedLayers[j] = t
	}

	return wrappedLayers, nil
}

// loadFromFile loads a file as an image for use within Kurma.
func loadFromFile(f tempfile.ReadSeekCloser, imageManager backend.ImageManager) (string, *schema.ImageManifest, error) {
	hash, manifest, err := imageManager.CreateImage(f)
	if err != nil {
		return "", nil, err
	}
	return hash, manifest, nil
}
