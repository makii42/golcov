package runner

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"

	"github.com/makii42/golcov/osadapter"
)

type (
	TestRunner interface {
		Run() (io.Reader, error)
	}
	testRunner struct {
		osa      osadapter.OS
		goBinary string
		Out      io.Writer
		Packages []string
	}
	testError struct {
		rc       int
		pkg      string
		output   []byte
		original error
	}
	testOutcome struct {
		output    []byte
		coverFile osadapter.File
	}
)

var (
	tempfilePrefix = "golcov-coverage-"
)

func NewTestRunner(osa osadapter.OS, out io.Writer, pkgs ...string) (TestRunner, error) {
	goBinary, err := osa.LookPath("go")
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		pkgs = append(pkgs, ".")
	}

	return &testRunner{
		osa:      osa,
		goBinary: goBinary,
		Out:      out,
		Packages: pkgs,
	}, nil
}

func (tr *testRunner) Run() (io.Reader, error) {
	var covers []io.Reader
	for _, pkg := range tr.Packages {
		outcome, err := tr.oneTest(pkg)
		if err != nil {
			// BOOM! a test run blew up!
			return nil, err
		}
		covers = append(covers, outcome.coverFile)
	}
	return io.MultiReader(covers...), nil
}

func (tr *testRunner) oneTest(pkg string) (*testOutcome, error) {
	f, err := tr.osa.TempFile("", tempfilePrefix)
	if err != nil {
		return nil, err
	}
	args := []string{
		"test",
		"-cover",
		"-coverprofile",
		f.Name(),
		"-v",
		pkg,
	}
	cmd := tr.osa.Command(tr.goBinary, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		var rc int
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				rc = status.ExitStatus()
			} else if !exitErr.Success() {
				rc = 999
			}
		} else {
			rc = 998
		}
		return nil, newTestError(
			rc,
			pkg,
			output,
			err,
		)
	}
	return &testOutcome{
		output:    output,
		coverFile: f,
	}, nil
}

func newTestError(rc int, pkg string, out []byte, err error) *testError {
	return &testError{rc, pkg, out, err}
}

func (e *testError) Error() string {
	return fmt.Sprintf("error during test execution for pkg '%s'", e.pkg)
}

func (o *testOutcome) Read(p []byte) (n int, err error) {
	return o.coverFile.Read(p)
}
