// Copyright 2016 Apcera Inc. All rights reserved.

// +build ignore cli

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/apcera/kurma/stager/container/core"
	"github.com/apcera/kurma/stager/container/run"

	"github.com/opencontainers/runc/libcontainer"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
)

func init() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		factory, _ := libcontainer.New("")
		if err := factory.StartInitialization(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n\r", err)
			os.Exit(1)
			return
		}
		panic("--this line should have never been executed, congratulations--")
	}
}

func main() {
	execName := filepath.Base(os.Args[0])

	var execFunc func() error
	switch execName {
	case "stager":
		execFunc = core.Run
	case "run":
		execFunc = run.Run
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized command %q\n", execName)
		os.Exit(2)
	}

	if err := execFunc(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n\r", err)
		os.Exit(1)
	}
}
