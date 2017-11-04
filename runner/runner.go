package runner

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/makii42/golcov/osadapter"
	"github.com/makii42/golcov/test"
)

type (
	TestRunner interface {
		Run() (io.Reader, error)
	}
	testRunner struct {
		osa      osadapter.OS
		goBinary string
		Out      io.Writer
		tests    []test.Test
	}
)

// NewTestRunner creates a new runner that executes all tests in the packages specified.
// If not packages are specified, it will discover all packages containing go sources,
// excluding some, like `vendor`.
func NewTestRunner(goBinary string, osa osadapter.OS, out io.Writer, tests ...test.Test) (TestRunner, error) {
	// todo: verify gobin.
	if len(tests) == 0 {
		return nil, fmt.Errorf("no tests specified")
	}
	return &testRunner{
		osa:      osa,
		goBinary: goBinary,
		Out:      out,
		tests:    tests,
	}, nil
}

// Run runs the tests in this runner.
func (tr *testRunner) Run() (io.Reader, error) {
	var covers []io.Reader
	for _, test := range tr.tests {
		outcome, err := test.Run()
		if err != nil {
			// BOOM! a test run blew up!
			return nil, err
		}
		covers = append(covers, outcome)
	}
	return io.MultiReader(covers...), nil
}

func (tr *testRunner) DiscoverPkgs(p string) ([]string, error) {
	folders := make(map[string]int)
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			folders[path] = 0
		} else if strings.HasSuffix(path, ".go") {

		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	pkgs := []string{}
	for k, v := range folders {
		if v > 0 {
			pkgs = append(pkgs, k)
		}
	}
	return pkgs, nil
}
