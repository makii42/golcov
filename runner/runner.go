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
	Discoverer interface {
		DiscoverPkgs(p string) ([]string, error)
	}
	testRunner struct {
		osa     osadapter.OS
		tests   []test.Test
		ignored []string
	}
	discoverer struct {
		osa osadapter.OS
	}
)

// NewTestRunner creates a new runner that executes all tests in the packages specified.
// If not packages are specified, it will discover all packages containing go sources,
// excluding some, like `vendor`.
func NewTestRunner(osa osadapter.OS, tests ...test.Test) (TestRunner, error) {
	if len(tests) == 0 {
		return nil, fmt.Errorf("no tests specified")
	}
	return &testRunner{
		osa:   osa,
		tests: tests,
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

func NewDiscoverer(osa osadapter.OS) Discoverer {
	return &discoverer{
		osa: osa,
	}
}
func (d *discoverer) DiscoverPkgs(p string) ([]string, error) {
	folders := make(map[string]int)
	cwd, err := d.osa.Getwd()
	if err != nil {
		return nil, err
	}
	err = d.osa.Walk(p, discoverFunc(cwd, folders))
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

func discoverFunc(cwd string, folders map[string]int) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			folders[path] = 0
		} else if strings.HasSuffix(path, ".go") {
			c := folders[path]
			c++
			folders[path] = c
		}
		return nil
	}
}
