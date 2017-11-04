package test

//go:generate mockgen -destination ../mocks/test/mocks.go -package test github.com/makii42/golcov/test Test,Outcome
import (
	"fmt"
	"io"
	"os/exec"
	"syscall"

	"github.com/makii42/golcov/osadapter"
)

var (
	tempfilePrefix = "golcov-coverage-"
)

type (
	// Test represents one `go test` run with exactly one package
	Test interface {
		Run() (Outcome, error)
	}
	// Outcome is a test run outcome, bundleing a pointer to a coverage file
	// as well the console output.
	Outcome interface {
		io.Reader
		ConsoleOutput() []byte
		CoverFile() osadapter.File
	}
	test struct {
		goBin   string
		osa     osadapter.OS
		pkg     string
		outcome outcome
	}
	outcome struct {
		consoleOutput []byte
		coverFile     osadapter.File
	}
	testError struct {
		rc       int
		pkg      string
		output   []byte
		original error
	}
)

func NewTest(goBin, pkg string, osa osadapter.OS) Test {
	return &test{
		goBin: goBin,
		pkg:   pkg,
		osa:   osa,
	}
}

func (t *test) Run() (Outcome, error) {
	f, err := t.osa.TempFile("", tempfilePrefix)
	if err != nil {
		return nil, err
	}
	args := []string{
		"test",
		"-cover",
		"-coverprofile",
		f.Name(),
		"-v",
		t.pkg,
	}
	cmd := t.osa.Command(t.goBin, args...)
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
			t.pkg,
			output,
			err,
		)
	}
	return &outcome{
		consoleOutput: output,
		coverFile:     f,
	}, nil
}

func (o *outcome) Read(p []byte) (n int, err error) {
	return o.coverFile.Read(p)
}

func (o *outcome) ConsoleOutput() []byte {
	return o.consoleOutput
}

func (o *outcome) CoverFile() osadapter.File {
	return o.coverFile
}

func newTestError(rc int, pkg string, out []byte, err error) *testError {
	return &testError{rc, pkg, out, err}
}

func (e *testError) Error() string {
	return fmt.Sprintf("error during test execution for pkg '%s'", e.pkg)
}
