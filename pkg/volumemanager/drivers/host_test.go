// Copyright 2016 Apcera Inc. All rights reserved.

package drivers

import (
	"testing"

	"github.com/apcera/kurma/pkg/backend/mocks"
	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"

	tt "github.com/apcera/util/testtool"
)

func TestHostVolumesProvision(t *testing.T) {
	tt.StartTest(t)
	defer tt.FinishTest(t)

	driver := &hostVolumeDriver{
		baseDirectory: "/example",
	}

	trueVal := true
	podSpec := schema.PodManifest{
		Volumes: []types.Volume{
			types.Volume{
				Name:     "test1",
				Kind:     "host",
				Source:   "/somewhere/over/there",
				ReadOnly: &trueVal,
			},
		},
	}
	pod := mocks.Pod{
		PodManifestFunc: func() *schema.PodManifest { return podSpec },
	}

	mounts, err := driver.Provision(pod)
	tt.TestExpectSuccess(t, err)
	tt.TestEqual(t, len(mounts), 1)
}

func TestHostVolumesCantEscapeBase(t *testing.T) {
	tt.StartTest(t)
	defer tt.FinishTest(t)

	driver := &hostVolumeDriver{
		baseDirectory: "/example",
	}

	// Success cases

	p, err := driver.getHostPath("a/b")
	tt.TestExpectSuccess(t, err)
	tt.TestEqual(t, p, "/example/a/b")

	p, err = driver.getHostPath("/a/b")
	tt.TestExpectSuccess(t, err)
	tt.TestEqual(t, p, "/example/a/b")

	p, err = driver.getHostPath("/a/b/../c")
	tt.TestExpectSuccess(t, err)
	tt.TestEqual(t, p, "/example/a/c")

	p, err = driver.getHostPath("/")
	tt.TestExpectSuccess(t, err)
	tt.TestEqual(t, p, "/example")

	p, err = driver.getHostPath("")
	tt.TestExpectSuccess(t, err)
	tt.TestEqual(t, p, "/example")

	// Failure cases

	_, err = driver.getHostPath("../a/b")
	tt.TestExpectError(t, err)
	tt.TestEqual(t, err.Error(), "the provided source path is escaping the base volume path")

	_, err = driver.getHostPath("/../a/b")
	tt.TestExpectError(t, err)
	tt.TestEqual(t, err.Error(), "the provided source path is escaping the base volume path")

	_, err = driver.getHostPath("/a/../../c")
	tt.TestExpectError(t, err)
	tt.TestEqual(t, err.Error(), "the provided source path is escaping the base volume path")
}
