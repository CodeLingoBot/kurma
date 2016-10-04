// Copyright 2015-2016 Apcera Inc. All rights reserved.

package imagestore

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFetch_LocalFile(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "localACi")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer f.Close()

	uri := "file://" + f.Name()

	m := &Manager{
		Options: &Options{
			FetchConfig: &FetchConfig{},
		},
	}

	readers, err := fetch(uri, m)
	if err != nil {
		t.Fatalf("Expected no error retrieving %s; got %s", uri, err)
	}

	for _, r := range readers {
		r.Close()
	}
}

func TestFetch_UnsupportedScheme(t *testing.T) {
	uri := "fakescheme://google.com"

	m := &Manager{
		Options: &Options{
			FetchConfig: &FetchConfig{},
		},
	}

	_, err := fetch(uri, m)
	if err == nil {
		t.Fatalf("Expected error with URI %q, got none", uri)
	}
}
