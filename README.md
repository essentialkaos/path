<p align="center"><a href="#readme"><img src="https://gh.kaos.st/path.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/path"><img src="https://kaos.sh/r/path.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/l/path"><img src="https://kaos.sh/l/6d6a56ab8cf3884d8523.svg" alt="Code Climate Maintainability" /></a>
  <a href="https://kaos.sh/b/path"><img src="https://kaos.sh/b/ac5eb5c7-1a0d-4223-884c-f99d4efaf77a.svg" alt="Codebeat badge" /></a>
  <a href="https://kaos.sh/w/path/ci"><img src="https://kaos.sh/w/path/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/path/codeql"><img src="https://kaos.sh/w/path/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`path` is a dead simple tool for working with paths. This tool provides commands which you can use to replace such tools as `basename`, `dirname`, and `readlink` and many more. But unlike these tools, `path` allows you to pass input not only as arguments, but also using standard input (_for example with pipes_). It's easy to use and doesn't require to know all this kung-fu with `find` or `xargs`.

Simple examples:

```bash
find . -iname '*.txt' -print0 | xargs -0 -n1 -- basename
# or
find . -iname '*.txt' | xargs -L1 -I{} basename "{}"
# with path
find . -iname '*.txt' | path basename
```

```bash
# Note that there is two spaces between {} and \; and if you forget
# about this it will don't work. Also in this case we will run 'basename'
# for each item in find output.
find . -mindepth 1 -maxdepth 1 -type d -exec basename {}  \;
# with path
find . -mindepth 1 -maxdepth 1 -type d | path basename
```

Also, it works MUCH faster:

```bash
time find . -iname '*.go' -print0 | xargs -0 -n1 -- basename
# find . -iname '*.go' -print0  0.07s user 0.12s system 4% cpu 4.255 total
# xargs -0 -n1 -- basename  2.59s user 2.74s system 102% cpu 5.195 total

time find . -iname '*.go' | path basename
# find . -iname '*.go'  0.08s user 0.09s system 99% cpu 0.174 total
# path basename  0.02s user 0.02s system 21% cpu 0.207 total
```

### Installation

#### From source

To build the `path` from scratch, make sure you have a working Go 1.20+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/path@latest
```

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st) for EL 7/8/9

```bash
sudo yum install -y https://yum.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo yum install path
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/path/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) path
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo path --completion=bash 1> /etc/bash_completion.d/path
```

ZSH:
```bash
sudo path --completion=zsh 1> /usr/share/zsh/site-functions/path
```

Fish:
```bash
sudo path --completion=fish 1> /usr/share/fish/vendor_completions.d/path.fish
```

### Man documentation

You can generate man page using next command:

```bash
path --generate-man | sudo gzip > /usr/share/man/man1/path.1.gz
```

### Usage

```
Usage: path {options} {command}

Commands

  base                 Strip directory and suffix from filenames
  dir                  Strip last component from file name
  link                 Print resolved symbolic links or canonical file names
  clean                Print shortest path name equivalent to path by purely lexical processing
  compact              Converts path to compact representation
  abs                  Print absolute representation of path
  ext                  Print file extension
  match pattern        Filter given path using pattern
  join root            Join path elements
  add-prefix prefix    Add the substring at the beginning
  del-prefix prefix    Remove the substring at the beginning
  add-suffix suffix    Add the substring at the end
  del-suffix suffix    Remove the substring at the end
  exclude substr       Exclude part of the string
  is-abs               Check if given path is absolute
  is-local             Check if given path is local
  is-safe              Check if given path is safe
  is-match pattern     Check if given path is match to pattern

Options

  --zero, -z         End each output line with NUL, not newline
  --space, -s        End each output line with space, not newline
  --quiet, -q        Suppress all error messages
  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

Examples

  path base /path/to/file.txt
  → file.txt

  path dir /path/to/file.txt
  → /path/to

  path compact /very/long/path/to/some/file.txt
  → /v/l/p/t/s/file.txt

  ls -1 | path is-match '*.txt' && echo MATCH!
  Check if all files in current directory is match to pattern
```

### CI Status

| Branch | Status |
|--------|----------|
| `master` | [![CI](https://kaos.sh/w/path/ci.svg?branch=master)](https://kaos.sh/w/path/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/path/ci.svg?branch=develop)](https://kaos.sh/w/path/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
