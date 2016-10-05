// Copyright 2016 Apcera Inc. All rights reserved.

package drivers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/apcera/logray"
	"github.com/opencontainers/runc/libcontainer/configs"

	"github.com/apcera/kurma/pkg/backend"
)

type hostVolumeDriver struct {
	log *logray.Logger

	baseDirectory string
}

// NewHostVolumeDriver handles creating a new volume driver that supports
// handling "host" volume types on pods.
func NewHostVolumeDriver(baseDirectory string) (backend.VolumeDriver, error) {
	// Attempt to create the directory, in the case it doesn't yet exist
	if err := os.MkdirAll(baseDirectory, os.FileMode(0750)); err != nil {
		return nil, fmt.Errorf("failed to create \"host\" volume path: %v", err)
	}
	return &hostVolumeDriver{baseDirectory: baseDirectory}, nil
}

// Kind returns the kind of volumes this implementation handles.
func (d *hostVolumeDriver) Kind() string {
	return "host"
}

// SetLog sets the logger to be used by the driver.
func (d *hostVolumeDriver) SetLog(log *logray.Logger) {
	d.log = log
}

// Validate handles validating the "host" volumes on the Pod are configured
// properly.
func (d *hostVolumeDriver) Validate(pod backend.Pod) error {
	return nil
}

// Provision handles setting up the "host" volumes on the pod.
func (d *hostVolumeDriver) Provision(pod backend.Pod) ([]*configs.Mount, error) {
	mounts := []*configs.Mount{}

	for _, volume := range pod.PodManifest().Volumes {
		if volume.Kind != "host" {
			continue
		}

		// Determine and validate the source path
		p, err := d.getHostPath(volume.Source)
		if err != nil {
			return nil, err
		}

		// populate the libcontainer mount
		m := &configs.Mount{
			Source:      p,
			Destination: filepath.Join("/volumes", volume.Name.String()),
			Device:      "bind",
			Flags:       syscall.MS_BIND,
		}

		if volume.ReadOnly != nil && *volume.ReadOnly {
			m.Flags |= syscall.MS_RDONLY
		}
		// TODO: Recusrive -> MS_REC

		mounts = append(mounts, m)
	}

	return mounts, nil
}

// Deprovision handles cleaning any "host" volumes. This is essentially a no-op.
func (d *hostVolumeDriver) Deprovision(pod backend.Pod) error {
	return nil
}

func (d *hostVolumeDriver) getHostPath(source string) (string, error) {
	p := filepath.Join(d.baseDirectory, source)

	// Ensure the joined path is still relative to the base path. filepath.Join
	// also calls filepath.Clean, so it should be cleaned up.
	if !strings.HasPrefix(p, fmt.Sprintf("%s%c", d.baseDirectory, filepath.Separator)) && p != d.baseDirectory {
		return "", fmt.Errorf("the provided source path is escaping the base volume path")
	}

	return p, nil
}
