package main

import (
	"log"
	"os"

	"github.com/makii42/golcov/osadapter"
	"github.com/makii42/golcov/runner"

	"github.com/urfave/cli"
)

func main() {
	app := cli.App{
		Name: "golcov",
		Commands: []cli.Command{
			testCmd,
		},
		Before: setup,
		Action: testAction,
	}
	app.Run(os.Args)
}

var (
	osa osadapter.OS
)

var testCmd = cli.Command{
	Name:   "test",
	Usage:  "Runs go test and collects coverage information for them. If the tests fail, test output is printed, otherwise coverage data.",
	Action: testAction,
}

func setup(c *cli.Context) error {
	osa = osadapter.RealOS()
	return nil
}

func testAction(c *cli.Context) {
	args := c.Args()
	r, err := runner.NewTestRunner(osa, nil, args...)
	if err != nil {
		log.Printf("cannot create testrunner: %s", err.Error())
		os.Exit(1)
	}
	coverSource, err := r.Run()
	if err != nil {
		log.Printf("error running tests: %s", err.Error())
		os.Exit(2)
	}
	osa.Copy(os.Stdout, coverSource)
}
