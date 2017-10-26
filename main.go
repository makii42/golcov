package main

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
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
	fs   afero.Fs
	osif OS
)

var testCmd = cli.Command{
	Name:   "test",
	Usage:  "Runs go test and collects coverage information for them. If the tests fail, test output is printed, otherwise coverage data.",
	Action: testAction,
}

func setup(c *cli.Context) error {
	fs = afero.NewOsFs()
	osif = RealOS()
	return nil
}

func testAction(c *cli.Context) {
	fmt.Println("Hello World")
}
