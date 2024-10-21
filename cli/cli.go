package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic utility info
const (
	APP  = "path"
	VER  = "1.0.3"
	DESC = "Dead simple tool for working with paths"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_ZERO     = "z:zero"
	OPT_SPACE    = "s:space"
	OPT_QUIET    = "q:quiet"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	CMD_BASENAME    = "base"
	CMD_DIRNAME     = "dir"
	CMD_DIRNAME_NUM = "dirn"
	CMD_READLINK    = "link"
	CMD_CLEAN       = "clean"
	CMD_COMPACT     = "compact"
	CMD_EXT         = "ext"
	CMD_ABS         = "abs"
	CMD_MATCH       = "match"
	CMD_JOIN        = "join"

	CMD_ADD_PREFIX = "add-prefix"
	CMD_DEL_PREFIX = "del-prefix"
	CMD_ADD_SUFFIX = "add-suffix"
	CMD_DEL_SUFFIX = "del-suffix"
	CMD_EXCLUDE    = "exclude"

	CMD_IS_ABS   = "is-abs"
	CMD_IS_LOCAL = "is-local"
	CMD_IS_SAFE  = "is-safe"
	CMD_IS_MATCH = "is-match"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap contains information about all supported options
var optMap = options.Map{
	OPT_ZERO:     {Type: options.BOOL},
	OPT_SPACE:    {Type: options.BOOL},
	OPT_QUIET:    {Type: options.BOOL},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// quietMode is quiet mode flag
var quietMode bool

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	runtime.GOMAXPROCS(1)

	preConfigureUI()

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.Error("- "))
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			Print()
		os.Exit(0)
	case options.GetB(OPT_HELP) || len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	err, ok := process(args)

	if err != nil {
		printError(err.Error())
	}

	if !ok {
		os.Exit(1)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	quietMode = options.GetB(OPT_QUIET) || os.Getenv("PATH_QUIET") != ""
}

// process starts arguments processing
func process(args options.Arguments) (error, bool) {
	cmd := args.Get(0).String()
	cmdArgs := args[1:]

	switch strings.ToLower(cmd) {
	case CMD_BASENAME, "basename":
		return cmdBasename(cmdArgs)
	case CMD_DIRNAME, "dirname":
		return cmdDirname(cmdArgs)
	case CMD_DIRNAME_NUM:
		return cmdDirnameNum(cmdArgs)
	case CMD_READLINK, "readlink":
		return cmdReadlink(cmdArgs)
	case CMD_CLEAN:
		return cmdClean(cmdArgs)
	case CMD_COMPACT:
		return cmdCompact(cmdArgs)
	case CMD_ABS:
		return cmdAbs(cmdArgs)
	case CMD_EXT:
		return cmdExt(cmdArgs)
	case CMD_MATCH:
		return cmdMatch(cmdArgs)
	case CMD_JOIN:
		return cmdJoin(cmdArgs)

	case CMD_ADD_PREFIX:
		return cmdAddPrefix(cmdArgs)
	case CMD_DEL_PREFIX:
		return cmdDelPrefix(cmdArgs)
	case CMD_ADD_SUFFIX:
		return cmdAddSuffix(cmdArgs)
	case CMD_DEL_SUFFIX:
		return cmdDelSuffix(cmdArgs)
	case CMD_EXCLUDE:
		return cmdExclude(cmdArgs)

	case CMD_IS_ABS:
		return cmdIsAbs(cmdArgs)
	case CMD_IS_LOCAL:
		return cmdIsLocal(cmdArgs)
	case CMD_IS_SAFE:
		return cmdIsSafe(cmdArgs)
	case CMD_IS_MATCH:
		return cmdIsMatch(cmdArgs)
	}

	return fmt.Errorf("Unknown command %q", cmd), false
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	if quietMode {
		return
	}

	terminal.Error(f, a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, APP))
	case "fish":
		fmt.Print(fish.Generate(info, APP))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, APP))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo()

	info.AddCommand(CMD_BASENAME, "Strip directory and suffix from filenames", "?path")
	info.AddCommand(CMD_DIRNAME, "Strip last component from file name", "?path")
	info.AddCommand(CMD_DIRNAME_NUM, "Return N elements from path", "num", "?path")
	info.AddCommand(CMD_READLINK, "Print resolved symbolic links or canonical file names", "?path")
	info.AddCommand(CMD_CLEAN, "Print shortest path name equivalent to path by purely lexical processing", "?path")
	info.AddCommand(CMD_COMPACT, "Converts path to compact representation", "?path")
	info.AddCommand(CMD_ABS, "Print absolute representation of path", "?path")
	info.AddCommand(CMD_EXT, "Print file extension", "?path")
	info.AddCommand(CMD_MATCH, "Filter given path using pattern", "pattern", "?path")
	info.AddCommand(CMD_JOIN, "Join path elements", "root", "?path")

	info.AddCommand(CMD_ADD_PREFIX, "Add the substring at the beginning", "prefix", "?path")
	info.AddCommand(CMD_DEL_PREFIX, "Remove the substring at the beginning", "prefix", "?path")
	info.AddCommand(CMD_ADD_SUFFIX, "Add the substring at the end", "suffix", "?path")
	info.AddCommand(CMD_DEL_SUFFIX, "Remove the substring at the end", "suffix", "?path")
	info.AddCommand(CMD_EXCLUDE, "Exclude part of the string", "substr", "?path")

	info.AddCommand(CMD_IS_ABS, "Check if given path is absolute", "?path")
	info.AddCommand(CMD_IS_LOCAL, "Check if given path is local", "?path")
	info.AddCommand(CMD_IS_SAFE, "Check if given path is safe", "?path")
	info.AddCommand(CMD_IS_MATCH, "Check if given path is match to pattern", "pattern", "?path")

	info.AddOption(OPT_ZERO, "End each output line with NUL, not newline")
	info.AddOption(OPT_SPACE, "End each output line with space, not newline")
	info.AddOption(OPT_QUIET, "Suppress all error messages")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		CMD_BASENAME+" /path/to/file.txt",
		"→ file.txt",
	)

	info.AddExample(
		CMD_DIRNAME+" /path/to/file.txt",
		"→ /path/to",
	)

	info.AddExample(
		CMD_COMPACT+" /very/long/path/to/some/file.txt",
		"→ /v/l/p/t/s/file.txt",
	)

	info.AddRawExample(
		"ls -1 | path "+CMD_IS_MATCH+" '*.txt' && echo MATCH!",
		"Check if all files in current directory is match to pattern",
	)

	info.AddRawExample(
		"PATH_QUIET=1 path "+CMD_DIRNAME+" /path/to/file.txt",
		"Run "+CMD_DIRNAME+" command in quiet mode enabled by environment variable",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //
