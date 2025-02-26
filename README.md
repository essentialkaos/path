<p align="center"><a href="#readme"><img src=".github/images/card.svg"/>
</a></p>

<p align="center">
  <a href="https://kaos.sh/l/path"><img src="https://kaos.sh/l/6d6a56ab8cf3884d8523.svg" alt="Code Climate Maintainability" /></a>
  <a href="https://kaos.sh/y/path"><img src="https://kaos.sh/y/cc438c3893ea48e58d8e3d935400e9e2.svg" alt="Codacy badge" /></a>
  <a href="https://kaos.sh/w/path/ci"><img src="https://kaos.sh/w/path/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/path/codeql"><img src="https://kaos.sh/w/path/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
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

Also, it works **MUCH** faster (~120x):

```
$ git clone https://github.com/kubernetes/kubernetes.git --depth=1

$ cd kubernetes

$ hyperfine 'find . -iname *.go -print0 | xargs -0 -n1 -- basename' 'find . -iname *.go | path basename'

Benchmark 1: find . -iname *.go -print0 | xargs -0 -n1 -- basename
  Time (mean ± σ):     12.621 s ±  0.077 s    [User: 5.871 s, System: 7.043 s]
  Range (min … max):   12.512 s … 12.745 s    10 runs

Benchmark 2: find . -iname *.go | path basename
  Time (mean ± σ):     106.5 ms ±   1.5 ms    [User: 59.8 ms, System: 60.4 ms]
  Range (min … max):   104.1 ms … 111.1 ms    28 runs

Summary
  find . -iname *.go | path basename ran
  118.45 ± 1.80 times faster than find . -iname *.go -print0 | xargs -0 -n1 -- basename
```

### Installation

#### From source

To build the `path` from scratch, make sure you have a working [Go 1.23+](https://github.com/essentialkaos/.github/blob/master/GO-VERSION-SUPPORT.md) workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/path@latest
```

#### From [ESSENTIAL KAOS Public Repository](https://kaos.sh/kaos-repo) for EL 8/9

```bash
sudo dnf install -y https://pkgs.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo dnf install path
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

<img src=".github/images/usage.svg"/>

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
