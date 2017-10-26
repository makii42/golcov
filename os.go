package main

//go:generate mockgen -destination ./os_mocks.go -package main -source os.go
// go:generate mockgen -destination ./apiclient_mocks.go -package apiclient -source=apiclient.go

import (
	"io"
	"os/exec"
)

type (
	// OS is our indirection of the pieces from os and os/exec that we use,
	// so we can mock those functions in our tests.
	OS interface {
		Command(cmd string, args ...string) Command
		LookPath(file string) (string, error)
	}
	// OSExecCommand is the interface wrapping all functions on exec.Cmd
	OSExecCommand interface {
		CombinedOutput() ([]byte, error)
		Output() ([]byte, error)
		Run() error
		Start() error
		StderrPipe() (io.ReadCloser, error)
		StdinPipe() (io.WriteCloser, error)
		StdoutPipe() (io.ReadCloser, error)
		Wait() error
	}

	Command interface {
		OSExecCommand
		Path() string
		Dir() string
		Args() []string
	}

	// fake os interface
	osproxy struct{}

	cmd struct {
		*exec.Cmd
	}
)

// RealOS returns an OS that is backed by the os- and sub-package.
func RealOS() OS {
	return &osproxy{}
}

func (o *osproxy) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (o *osproxy) Command(name string, arg ...string) Command {
	realCmd := exec.Command(name, arg...)
	return &cmd{
		realCmd,
	}
}

func (c *cmd) Path() string {
	return c.Cmd.Path
}

func (c *cmd) Dir() string {
	return c.Cmd.Dir
}

func (c *cmd) Args() []string {
	return c.Cmd.Args
}
