// Copyright 2016 Apcera Inc. All rights reserved.

package docker

import (
	"fmt"
	"net/url"
	"testing"

	v2mock "github.com/apcera/util/dockertest/v2"
)

var dockerRegistryURL string

func init() {
	v2Registry := v2mock.RunMockRegistry()
	dockerRegistryURL = v2Registry.URL
}

func TestDockerPull_ImageNotFound(t *testing.T) {
	puller := New(true)

	imageURI := "docker://fake"

	_, err := puller.Pull(imageURI)
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

func TestDockerPull(t *testing.T) {
	puller := New(true)

	imageURI := fmt.Sprintf("%s/library/nats:latest", dockerRegistryURL)
	u, err := url.Parse(imageURI)
	if err != nil {
		t.Fatalf("Failed to parse %q as URL: %s", imageURI, err)
	}

	u.Scheme = "docker"
	_, err = puller.Pull(u.String())
	if err != nil {
		t.Fatalf("Failed to pull %q: %s", imageURI, err)
	}
}
