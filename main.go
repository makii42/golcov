package main

import (
	"fmt"
	"log"
	"os"

	"github.com/makii42/golcov/osadapter"
	"github.com/makii42/golcov/runner"
	"github.com/makii42/golcov/test"

	"github.com/urfave/cli"
)

func main() {
	app := cli.App{
		Name:  "golcov",
		Usage: "Runs go test with default coverage options for each package and writes it to standard out",
		Commands: []cli.Command{
			testCmd,
			versionCmd,
		},
		Before: setup,
		Action: testAction,
		Flags: []cli.Flag{
			cli.BoolTFlag{
				Name:  "vendor",
				Usage: "not sure about this yet",
			},
		},
	}
	app.Run(os.Args)
}

var (
	osa      osadapter.OS
	goBinary string
	version  string = "dev"
	revision string = "dev"
)

var testCmd = cli.Command{
	Name:   "test",
	Usage:  "Runs go test and collects coverage information for them. If the tests fail, test output is printed, otherwise coverage data.",
	Action: testAction,
}

func setup(c *cli.Context) error {
	osa = osadapter.RealOS()
	bin, err := osa.LookPath("go")
	if err != nil {
		return err
	}
	goBinary = bin
	return nil
}

func testAction(c *cli.Context) {
	args := c.Args()
	tests := createTests(args...)
	r, err := runner.NewTestRunner(osa, tests...)
	if err != nil {
		log.Printf("cannot create testrunner: %s", err.Error())
		os.Exit(1)
	}
	coverSource, err := r.Run()
	if err != nil {
		tf, ok := test.IsTestFailure(err)
		if ok {
			fmt.Fprintf(os.Stderr, "error running tests in pkg %s\n", tf.Pkg)
			fmt.Fprintf(os.Stderr, "Output:\n--------------------------------\n%s\n--------------------------------\n", string(tf.Output))
			os.Exit(2)
		}
		log.Printf("error running tests: %s", err.Error())
		os.Exit(3)
	}
	osa.Copy(os.Stdout, coverSource)
}

func createTests(pkgs ...string) (tests []test.Test) {
	tests = []test.Test{}
	for _, pkg := range pkgs {
		tests = append(tests, test.NewTest(goBinary, pkg, osa))
	}
	return
}

var versionCmd = cli.Command{
	Name:  "version",
	Usage: "displays the version",
	Action: func(c *cli.Context) {
		fmt.Printf("%s version: %s (%s)\n", c.App.Name, version, revision)
	},
}
