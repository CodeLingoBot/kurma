// Copyright 2015-2016 Apcera Inc. All rights reserved.

package image

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/apcera/kurma/pkg/remote"
	"github.com/apcera/kurma/pkg/remote/aci"
	"github.com/apcera/kurma/pkg/remote/docker"
	"github.com/apcera/kurma/pkg/remote/http"

	"github.com/apcera/util/tempfile"

	"github.com/appc/spec/schema/types"
)

// A FetchConfig contains configuration for an image fetch operation.
type FetchConfig struct {
	// ACILabels are labels that scope the image resolution request.
	ACILabels map[types.ACIdentifier]string `json:"aciLabels,omitempty"`

	// Insecure is an option that, if enabled, will fetch images insecurely.
	// For instance, this will disable signature verification, and will use
	// HTTP for fetching images rather than HTTPS, where applicable.
	Insecure bool `json:"insecure"`
}

// Fetch retrieves a container image. Images may be sourced from the local
// machine, or may be retrieved from a remote server.
func Fetch(imageURI string, cfg *FetchConfiug) ([]tempfile.ReadSeekCloser, error) {
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
		puller = docker.New(cfg.Insecure)
	case "aci", "":
		puller = aci.New(cfg.Insecure, cfg.ACILabels)
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
