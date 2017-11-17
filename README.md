Golang Coverage
===============

[Golang][go] provides great tools for unit tests, and capturing as well as displaying coverage information.

With the advent of specialized tools to display those statistics as well as tracking them by branch and over time, `golcov` closes the gap between [lcov-server][lcov-gh] and golang.

To use it, grab the [latest release][rel] for your platform, or install the master HEAD with `go get`:

    go get github.com/makii42/golcov

As part of your CI Job, run your golang tests like this:

    golcov test $(go list ./...) | lcov-server --upload https://${YOUR_LCOV_SERVER}>/

This should take of uploading your coverage information.

What it does
------------

`golcov` runs `go test` under the hood, storing the coverage information in temporary files. Due to limitations in `go test` we have to run tests per package, as storing detailled information using the `-coverprofile` does not work for multiple packages at once.

If all tests are successful, the coverage information files are written to `os.Stdout`, where `lcov-server` can pick it up. 

If tests in one package fail, the test output is written to `os.Stderr`, so you can see the actual failure.

[go]: https://golang.org/
[lcov-gh]: https://github.com/gabrielcsapo/lcov-server
[rel]: https://github.com/makii42/golcov/releases