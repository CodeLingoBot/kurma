// Copyright 2016 Apcera Inc. All rights reserved.

package aci

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/apcera/kurma/pkg/remote"
	"github.com/apcera/kurma/pkg/remote/aci/server"
	"github.com/appc/spec/schema/types"

	tt "github.com/apcera/util/testtool"
)

func TestACIPull_ImageNotFound(t *testing.T) {
	puller := buildACIPuller(true)

	imageURI := "aci://fake"

	_, err := puller.Pull(imageURI)
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

func TestACIPull_InsecureHTTPDiscover(t *testing.T) {
	img := "etcd"

	server, imageName := bootstrapACIServerWithImage(t, img, true)
	defer server.Close()

	puller := buildACIPuller(true)

	acis, err := puller.Pull(imageName)
	if err != nil {
		t.Fatalf("Failed to pull %q: %s", imageName, err)
	}

	if len(acis) != 1 {
		t.Fatalf("Expected 1 ACI, got %d", len(acis))
	}

	for _, a := range acis {
		a.Close()
	}
}

func TestACIPull_SecureHTTPSDiscover(t *testing.T) {
	img := "etcd"

	server, imageName := bootstrapACIServerWithImage(t, img, false)
	defer server.Close()

	puller := buildACIPuller(true)

	_, err := puller.Pull(imageName)
	if err == nil {
		t.Fatal("Expected error pulling image with HTTPS enabled, got none")
	}
}

func buildACIPuller(insecure bool) remote.Puller {
	labels := make(map[types.ACIdentifier]string)
	labels[types.ACIdentifier("os")] = runtime.GOOS
	labels[types.ACIdentifier("arch")] = runtime.GOARCH

	return New(insecure, labels)
}

func bootstrapACIServerWithImage(t *testing.T, img string, insecure bool) (*server.Server, string) {
	// Binding ports like 80, 443 requires root.
	tt.TestRequiresRoot(t)

	setup := server.GetDefaultServerSetup()
	if insecure {
		setup.Protocol = server.ProtocolHttp
	}
	server := server.NewServer(setup)

	if testing.Verbose() {
		go logMessages(t, server)
	}

	imageName := fmt.Sprintf("localhost/%s", img)

	tmpDir := os.TempDir()

	image := writeTestACI(t, tmpDir, imageName)
	images := map[string]string{
		fmt.Sprintf("%s.aci", img): image,
	}

	if err := server.UpdateFileSet(images); err != nil {
		t.Fatalf("Failed to set images on remote server: %s", err)
	}
	return server, imageName
}

func writeTestACI(t *testing.T, tmpDir, imageName string) string {
	manifestTemplate := `{"acKind":"ImageManifest","acVersion":"0.7.4","name":"IMG_NAME","labels":[{"name":"version","value":"1.2.0"},{"name":"arch","value":"GOARCH"},{"name":"os","value":"GOOS"}]}`
	manifest := strings.Replace(manifestTemplate, "IMG_NAME", imageName, -1)
	manifest = strings.Replace(manifest, "GOARCH", runtime.GOARCH, -1)
	manifest = strings.Replace(manifest, "GOOS", runtime.GOOS, -1)

	tmpManifest, err := ioutil.TempFile(tmpDir, "manifest")
	if err != nil {
		t.Fatalf("Failed to create temp file for manifest: %s", err)
	}
	defer os.Remove(tmpManifest.Name())
	if err := ioutil.WriteFile(tmpManifest.Name(), []byte(manifest), 0600); err != nil {
		t.Fatalf("Failed to write temp manifest file: %s", err)
	}

	imageFileName := fmt.Sprintf("%s.aci", filepath.Base(imageName))

	image, err := filepath.Abs(filepath.Join(tmpDir, imageFileName))
	if err != nil {
		t.Fatalf("Failed to get absolute path: %s", err)
	}

	if err := os.Rename(tmpManifest.Name(), image); err != nil {
		t.Fatalf("Failed to rename %q to %q: %s", tmpManifest, imageFileName, err)
	}
	return image
}

func logMessages(t *testing.T, server *server.Server) {
	for {
		select {
		case msg, ok := <-server.Msg:
			if ok {
				t.Logf("--- ACI server msg: %v", msg)
			} else {
				return
			}
		}
	}
}
