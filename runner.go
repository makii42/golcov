package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/afero"
)

type (
	TestRunner interface {
		Run() error
	}
	testRunner struct {
		fs       afero.Fs
		os       OS
		goBinary string
		Out      io.Writer
		Packages []string
	}
)

func NewTestRunner(os OS, fs afero.Fs, out io.Writer, pkgs ...string) (TestRunner, error) {
	goBinary, err := os.LookPath("go")
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		pkgs = append(pkgs, ".")
	}

	return &testRunner{
		os:       os,
		fs:       fs,
		goBinary: goBinary,
		Out:      out,
		Packages: pkgs,
	}, nil
}

func (tr *testRunner) Run() error {
	rc, output, err := tr.oneTest(".")
	if err != nil {
		return err
	}
	if rc != 0 {
		if n, err := os.Stderr.Write(output); err != nil {
			return err
		} else if n != len(output) {
			return fmt.Errorf("could not write all %d bytes to stderr", len(output))
		}
		os.Exit(rc)
	}
	return nil
}

func (tr *testRunner) oneTest(pkg string) (int, []byte, error) {
	f, err := afero.TempFile(tr.fs, "", "golcov-coverage-")
	if err != nil {
		return -1, nil, err
	}
	args := []string{
		"test",
		"-cover",
		"-coverprofile",
		f.Name(),
		"-v",
	}
	cmd := tr.os.Command(tr.goBinary, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return -1, nil, err
	}
	return 0, output, nil
}
