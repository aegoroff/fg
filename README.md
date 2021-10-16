fg
==
[![CI](https://github.com/aegoroff/fg/actions/workflows/ci.yml/badge.svg)](https://github.com/aegoroff/fg/actions/workflows/ci.yml) [![codecov](https://codecov.io/gh/aegoroff/fg/branch/master/graph/badge.svg)](https://codecov.io/gh/aegoroff/fg) [![Go Report Card](https://goreportcard.com/badge/github.com/aegoroff/fg)](https://goreportcard.com/report/github.com/aegoroff/fg)

A small commandline app written in Go that allows you to easily group
all files in the dir specified into several child subdirectories.
Grouping uses file modification time or file name and may be converted into day, month
or year in the form yyyy-dd-mm or yyyy-mm or yyyy. Also you can group files by file extension and first 3 letters of the file name. This can be specified
using command line option.

## Install the pre-compiled binary

**homebrew** (only on macOS and Linux for now):

Add my tap (do it once):
```sh
brew tap aegoroff/tap
```
And then install fg:
```sh
brew install fg
```
Update fg if already installed:
```sh
brew upgrade fg
```

**scoop**:

```sh
scoop bucket add aegoroff https://github.com/aegoroff/scoop-bucket.git
scoop install fg
```

**manually**:

Download the pre-compiled binaries from the [releases](https://github.com/aegoroff/fg/releases) and
copy to the desired location.