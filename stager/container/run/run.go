// Copyright 2016 Apcera Inc. All rights reserved.

package run

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/apcera/kurma/schema"
	"github.com/opencontainers/runc/libcontainer"
)

func Run() error {
	// Read in the app configuration
	var app *schema.RunApp
	f := os.NewFile(3, "app.json")
	if err := json.NewDecoder(f).Decode(&app); err != nil {
		return err
	}
	f.Close()

	// Load the container with libcontainer
	factory, err := libcontainer.New("/containers")
	if err != nil {
		return err
	}
	container, err := factory.Load(os.Args[1])
	if err != nil {
		return err
	}

	// Allocate a wait group which is primarily used when a tty is requested, to
	// ensure all content is written before returning.
	wg := sync.WaitGroup{}

	// Setup the process
	workingDirectory := app.WorkingDirectory
	if workingDirectory == "" {
		workingDirectory = "/"
	}
	process := &libcontainer.Process{
		Cwd:    workingDirectory,
		User:   app.User,
		Args:   app.Exec,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	for _, env := range app.Environment {
		process.Env = append(process.Env, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}

	// Create a tty for the process if the caller wants it
	if app.Tty {
		console, err := process.NewConsole(os.Getuid())
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			io.Copy(os.Stdout, console)
			wg.Done()
		}()
		go io.Copy(console, os.Stdin)
	}

	// Run it!
	if err := container.Start(process); err != nil {
		return err
	}

	// Get the process status. Ignore the error, it'll always have an error if the
	// process exited non-zero.
	ps, _ := process.Wait()

	// Wait for other routines to finish up and flush output. We time box this to
	// 100ms. This is because the process can exit before we've finished
	// reading from its stdout and writing it out our stdout.
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		break
	case <-time.After(time.Millisecond * 100):
		break
	}

	// We'll explicitly exit here so we can propagate up the exit code
	if status, ok := ps.Sys().(syscall.WaitStatus); ok {
		if code := status.ExitStatus(); code != 0 {
			os.Exit(code)
		}
	}

	return nil
}
