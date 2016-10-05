// Copyright 2016 Apcera Inc. All rights reserved.

package drivers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/apcera/logray"
	"github.com/opencontainers/runc/libcontainer/configs"

	"github.com/apcera/kurma/pkg/backend"
)

type emptyVolumeDriver struct {
	log *logray.Logger

	volumeDirectory string
}

// NewEmptyVolumeDriver handles creating a new volume driver that supports
// handling "empty" volume types on pods.
func NewEmptyVolumeDriver(volumeDirectory string) (backend.VolumeDriver, error) {
	// Attempt to create the directory, in the case it doesn't yet exist
	if err := os.MkdirAll(volumeDirectory, os.FileMode(0750)); err != nil {
		return nil, fmt.Errorf("failed to create \"host\" volume path: %v", err)
	}
	return &emptyVolumeDriver{volumeDirectory: volumeDirectory}, nil
}

// Kind returns the kind of volumes this implementation handles.
func (d *emptyVolumeDriver) Kind() string {
	return "empty"
}

// SetLog sets the logger to be used by the driver.
func (d *emptyVolumeDriver) SetLog(log *logray.Logger) {
	d.log = log
}

// Validate handles validating the "empty" volumes on the Pod are configured
// properly.
func (d *emptyVolumeDriver) Validate(pod backend.Pod) error {
	return nil
}

// Provision handles setting up the "empty" volumes on the pod.
func (d *emptyVolumeDriver) Provision(pod backend.Pod) ([]*configs.Mount, error) {
	mounts := []*configs.Mount{}

	for _, volume := range pod.PodManifest().Volumes {
		if volume.Kind != "empty" {
			continue
		}

		mode := os.FileMode(0755)
		if volume.Mode != nil && *volume.Mode != "" {
			var err error
			mode, err = parseFileMode(*volume.Mode)
			if err != nil {
				return nil, err
			}
		}

		// Handle directory creation and mode
		p := filepath.Join(d.volumeDirectory, pod.UUID(), volume.Name.String())
		if err := os.MkdirAll(p, mode); err != nil {
			return nil, err
		}

		// Handle the owners
		uid, gid := 0, 0
		if volume.UID != nil {
			uid = *volume.UID
		}
		if volume.GID != nil {
			gid = *volume.GID
		}
		if err := os.Chown(p, uid, gid); err != nil {
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

		mounts = append(mounts, m)
	}

	return mounts, nil
}

// Deprovision handles removing any "empty" volumes that were created for a Pod
// after it has been shutdown.
func (d *emptyVolumeDriver) Deprovision(pod backend.Pod) error {
	p := filepath.Join(d.volumeDirectory, pod.UUID())
	if err := os.RemoveAll(p); err != nil {
		return fmt.Errorf("failed to remove 'empty' volumes on pod shutdown: %v", err)
	}
	return nil
}

func parseFileMode(mode string) (os.FileMode, error) {
	m, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse mode %q: %v", mode, err)
	}
	return os.FileMode(m), nil
}
