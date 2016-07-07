package core

import (
	"fmt"
	"testing"
)

type dummyStagerSetupOkRunner struct{}

func (s *dummyStagerSetupOkRunner) setupSignalHandling()       {}
func (s *dummyStagerSetupOkRunner) readManifest() error        { return nil }
func (s *dummyStagerSetupOkRunner) writeState() error          { return nil }
func (s *dummyStagerSetupOkRunner) populateState()             {}
func (s *dummyStagerSetupOkRunner) createFactory() error       { return nil }
func (s *dummyStagerSetupOkRunner) launchInit() error          { return nil }
func (s *dummyStagerSetupOkRunner) containerFilesystem() error { return nil }
func (s *dummyStagerSetupOkRunner) createContainers() error    { return nil }
func (s *dummyStagerSetupOkRunner) markRunning()               {}
func (s *dummyStagerSetupOkRunner) markShuttingDown()          {}
func (s *dummyStagerSetupOkRunner) signalReadyPipe()           {}
func (s *dummyStagerSetupOkRunner) stop()                      {}

func TestStagerRunOK(t *testing.T) {
	s := &dummyStagerSetupOkRunner{}
	err := run(s)
	if err != nil {
		t.Fatalf("expected no errors running the stager")
	}
}

type dummyStagerSetupRunner struct {
	state      string
	stopCalled bool
}

func (s *dummyStagerSetupRunner) setupSignalHandling() {}
func (s *dummyStagerSetupRunner) readManifest() error  { return nil }
func (s *dummyStagerSetupRunner) writeState() error    { return nil }
func (s *dummyStagerSetupRunner) populateState() {
	s.state = "setup"
}
func (s *dummyStagerSetupRunner) createFactory() error       { return nil }
func (s *dummyStagerSetupRunner) launchInit() error          { return nil }
func (s *dummyStagerSetupRunner) containerFilesystem() error { return fmt.Errorf("failed") }
func (s *dummyStagerSetupRunner) createContainers() error    { return nil }
func (s *dummyStagerSetupRunner) markRunning() {
	s.state = "running"
}
func (s *dummyStagerSetupRunner) markShuttingDown() {}
func (s *dummyStagerSetupRunner) signalReadyPipe()  {}
func (s *dummyStagerSetupRunner) stop() {
	s.stopCalled = true
}

func TestStagerRunFailsDuringSetup(t *testing.T) {
	s := &dummyStagerSetupRunner{}
	err := run(s)
	if err == nil {
		t.Fatalf("expected an error to come up during run")
	}
	got := s.state
	expected := "setup"
	if got != expected {
		t.Fatalf("expected for stager to be at state %q. got: %q", expected, got)
	}

	if !s.stopCalled {
		t.Fatalf("expected stop to have been called on failure")
	}
}
