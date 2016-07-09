// Copyright 2016 Apcera Inc. All rights reserved.

package core

import (
	"runtime"

	"github.com/apcera/logray"
	"github.com/opencontainers/runc/libcontainer"
)

func Run() error {
	logray.AddDefaultOutput("stdout://", logray.ALL)
	cs := &containerSetup{
		log:           logray.New(),
		stagerConfig:  defaultStagerConfig,
		appContainers: make(map[string]libcontainer.Container),
		appProcesses:  make(map[string]*libcontainer.Process),
		appWaitch:     make(map[string]chan struct{}),
	}

	if err := run(cs); err != nil {
		cs.log.Errorf("Startup function errored: %s", err)
		cs.log.Flush()
		return err
	}

	runtime.Goexit()
	return nil
}
