// Copyright 2016 Apcera Inc. All rights reserved.

package mocks

import (
	"io"
	"os"
	"time"

	"github.com/apcera/kurma/pkg/backend"
	"github.com/appc/spec/schema"

	ntypes "github.com/apcera/kurma/pkg/networkmanager/types"
	kschema "github.com/apcera/kurma/schema"
)

type Pod struct {
	UUIDFunc         func() string
	NameFunc         func() string
	PodManifestFunc  func() *schema.PodManifest
	NetworksFunc     func() []*ntypes.IPResult
	StateFunc        func() backend.PodState
	StopFunc         func() error
	EnterFunc        func(appName string, app *kschema.RunApp, stdin io.Reader, stdout, stderr io.Writer, postStart func()) (*os.Process, error)
	WaitForStateFunc func(timeout time.Duration, states ...backend.PodState) error
	WaitFunc         func()
}

func (p *Pod) UUID() string {
	return p.UUIDFunc()
}

func (p *Pod) Name() string {
	return p.NameFunc()
}

func (p *Pod) PodManifest() *schema.PodManifest {
	return p.PodManifestFunc()
}

func (p *Pod) Networks() []*ntypes.IPResult {
	return p.NetworksFunc()
}

func (p *Pod) State() backend.PodState {
	return p.StateFunc()
}

func (p *Pod) Stop() error {
	return p.StopFunc()
}

func (p *Pod) Enter(appName string, app *kschema.RunApp, stdin io.Reader, stdout, stderr io.Writer, postStart func()) (*os.Process, error) {
	return p.EnterFunc(appName, app, stdin, stdout, stderr, postStart)
}

func (p *Pod) WaitForState(timeout time.Duration, states ...backend.PodState) error {
	return p.WaitForStateFunc(timeout, states...)
}

func (p *Pod) Wait() {
	p.WaitFunc()
}
