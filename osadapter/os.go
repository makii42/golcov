package osadapter

//go:generate mockgen -destination ../mocks/osadapter_mocks.go -package mocks github.com/makii42/golcov/osadapter OS,Command,File

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type (
	// OS is our indirection of the pieces from os and os/exec that we use,
	// so we can mock those functions in our tests.
	OS interface {
		Command(cmd string, args ...string) Command
		LookPath(file string) (string, error)
		TempFile(dir, prefix string) (f File, err error)
		Copy(dst io.Writer, src io.Reader) (written int64, err error)
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
		GetPath() string
		GetDir() string
		GetArgs() []string
		String() string
	}

	File interface {
		Chdir() error
		Chmod(mode os.FileMode) error
		Chown(uid, gid int) error
		Close() error
		Fd() uintptr
		Name() string
		Read(b []byte) (n int, err error)
		ReadAt(b []byte, off int64) (n int, err error)
		Readdir(n int) ([]os.FileInfo, error)
		Readdirnames(n int) (names []string, err error)
		Seek(offset int64, whence int) (ret int64, err error)
		Stat() (os.FileInfo, error)
		Sync() error
		Truncate(size int64) error
		Write(b []byte) (n int, err error)
		WriteAt(b []byte, off int64) (n int, err error)
		WriteString(s string) (n int, err error)
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
func (o *osproxy) TempFile(dir, prefix string) (f File, err error) {
	return ioutil.TempFile(dir, prefix)
}

func (o *osproxy) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

func (o *osproxy) Command(name string, arg ...string) Command {
	realCmd := exec.Command(name, arg...)
	return &cmd{
		realCmd,
	}
}

func (c *cmd) GetPath() string {
	return c.Cmd.Path
}

func (c *cmd) GetDir() string {
	return c.Cmd.Dir
}

func (c *cmd) GetArgs() []string {
	return c.Cmd.Args
}

func (c *cmd) String() string {
	return fmt.Sprintf("%#v", *c.Cmd)
}
