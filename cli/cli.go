package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fsutil"
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
	VER  = "2.0.0"
	DESC = "Tool for working with paths"
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

	OPT_UPDATE       = "U:update"
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
	CMD_STRIP_EXT  = "strip-ext"
	CMD_EXCLUDE    = "exclude"
	CMD_REPLACE    = "replace"
	CMD_LOWER      = "lower"
	CMD_UPPER      = "upper"

	CMD_IS_ABS   = "is-abs"
	CMD_IS_LOCAL = "is-local"
	CMD_IS_SAFE  = "is-safe"
	CMD_IS_MATCH = "is-match"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// handlerFunc is a function for processing command data
type handlerFunc func(data string, args options.Arguments) (string, error, bool)

// handler contains base info for command handler
type handler struct {
	Func handlerFunc       // Handler function
	Args options.Arguments // Command arguments
}

// pipe is a slice of handler to process data
type pipe []*handler

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

// colorTagApp is app name color tag
var colorTagApp string

// colorTagVer is app version color tag
var colorTagVer string

// separator is data separator
var separator string

// hasStdinData is marker that shows that there some data in stdin
var hasStdinData bool

// ////////////////////////////////////////////////////////////////////////////////// //

// minCmdArgs contains minimum number of arguments
var minCmdArgs = map[string]int{
	CMD_DIRNAME_NUM: 1,
	CMD_MATCH:       1,
	CMD_JOIN:        1,
	CMD_ADD_PREFIX:  1,
	CMD_DEL_PREFIX:  1,
	CMD_ADD_SUFFIX:  1,
	CMD_DEL_SUFFIX:  1,
	CMD_EXCLUDE:     1,
	CMD_REPLACE:     2,
	CMD_IS_MATCH:    1,
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	runtime.GOMAXPROCS(1)

	preConfigureUI()
	preConfigureOptions()

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.Error(" - "))
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
	case withSelfUpdate && options.GetB(OPT_UPDATE):
		os.Exit(updateBinary())
	case options.GetB(OPT_HELP) || len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	err, ok := runCommands(args)

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

	switch {
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{#99}", "{#99}"
	default:
		colorTagApp, colorTagVer = "{*}{m}", "{m}"
	}
}

// preConfigureOptions preconfigures command-line options based on build tags
func preConfigureOptions() {
	optMap.SetIf(withSelfUpdate, OPT_UPDATE, &options.V{Type: options.MIXED})
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	quietMode = options.GetB(OPT_QUIET) || os.Getenv("PATH_QUIET") != ""

	switch {
	case options.GetB(OPT_SPACE):
		separator = " "
	case options.GetB(OPT_ZERO):
		separator = "\x00"
	default:
		separator = "\n"
	}

	if !fsutil.IsCharacterDevice("/dev/stdin") {
		hasStdinData = true
	}
}

// runCommands starts arguments processing
func runCommands(args options.Arguments) (error, bool) {
	var cmds pipe
	var err error
	var data []string
	var hdlr *handler

	cmd := args.Get(0).String()

	if strings.ContainsRune(cmd, ',') {
		cmds, err = parseCommandPipe(cmd)
		data = args[1:].Strings()
	} else {
		hdlr, data, err = createCommandHandler(cmd, args[1:])
		cmds = pipe{hdlr}
	}

	if err != nil {
		return err, false
	}

	if !hasStdinData && len(data) == 0 {
		return fmt.Errorf("There is no data for command"), false
	}

	if len(data) > 0 {
		err, ok := processArgsData(cmds, data)

		if err != nil || !ok {
			return err, ok
		}
	}

	if hasStdinData {
		err, ok := processStdinData(cmds)

		if err != nil || !ok {
			return err, ok
		}
	}

	return nil, true
}

// processArgsData runs commands over data passed as CLI arguments
func processArgsData(cmds pipe, data []string) (error, bool) {
	for _, str := range data {
		err, ok := executePipeHandlers(cmds, str)

		if err != nil || !ok {
			return err, ok
		}
	}

	return nil, true
}

// processStdinData runs commands over data passed via standard input
func processStdinData(cmds pipe) (error, bool) {
	r := bufio.NewReader(os.Stdin)
	delim := separator[0]

	for {
		str, err := r.ReadString(delim)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("Can't read stdin data: %v", err), false
			}
		}

		str = strings.TrimRight(str, separator)
		err, ok := executePipeHandlers(cmds, str)

		if err != nil || !ok {
			return err, ok
		}
	}

	return nil, false
}

// parseCommandPipe parses command pipe
func parseCommandPipe(data string) (pipe, error) {
	var result pipe

	for _, cmd := range strings.Split(data, ",") {
		var args options.Arguments
		var rawArgs string

		if strings.ContainsRune(cmd, '+') {
			cmd, rawArgs, _ = strings.Cut(cmd, "+")
			args = options.NewArguments(strings.Split(rawArgs, "+")...)
		}

		h, _, err := createCommandHandler(cmd, args)

		if err != nil {
			return nil, err
		}

		result = append(result, h)
	}

	return result, nil
}

// createCommandHandler returns handler for command
func createCommandHandler(cmd string, args options.Arguments) (*handler, []string, error) {
	cmd = strings.ToLower(cmd)
	minArgs := minCmdArgs[cmd]

	if minArgs > 0 && len(args) < minArgs {
		return nil, nil, fmt.Errorf("Not enough arguments for command %q", cmd)
	}

	switch strings.ToLower(cmd) {
	case CMD_BASENAME, "basename":
		return &handler{cmdBasename, nil}, args.Strings(), nil

	case CMD_DIRNAME, "dirname":
		return &handler{cmdDirname, nil}, args.Strings(), nil

	case CMD_DIRNAME_NUM:
		return &handler{cmdDirnameNum, args[:minArgs]}, args[minArgs:].Strings(), nil

	case CMD_READLINK, "readlink":
		return &handler{cmdReadlink, nil}, args.Strings(), nil

	case CMD_CLEAN:
		return &handler{cmdClean, nil}, args.Strings(), nil

	case CMD_COMPACT:
		return &handler{cmdCompact, nil}, args.Strings(), nil

	case CMD_ABS:
		return &handler{cmdAbs, nil}, args.Strings(), nil

	case CMD_EXT:
		return &handler{cmdExt, nil}, args.Strings(), nil

	case CMD_MATCH:
		return &handler{cmdMatch, args[:minArgs]}, args[minArgs:].Strings(), nil

	case CMD_JOIN:
		return &handler{cmdJoin, args[:minArgs]}, args[minArgs:].Strings(), nil

	case CMD_ADD_PREFIX:
		return &handler{cmdAddPrefix, args[:minArgs]}, args[minArgs:].Strings(), nil

	case CMD_DEL_PREFIX:
		return &handler{cmdDelPrefix, args[:minArgs]}, args[minArgs:].Strings(), nil

	case CMD_ADD_SUFFIX:
		return &handler{cmdAddSuffix, args[:minArgs]}, args[minArgs:].Strings(), nil

	case CMD_DEL_SUFFIX:
		return &handler{cmdDelSuffix, args[:minArgs]}, args[minArgs:].Strings(), nil

	case CMD_EXCLUDE:
		return &handler{cmdExclude, args[:minArgs]}, args[minArgs:].Strings(), nil

	case CMD_REPLACE:
		return &handler{cmdReplace, args[:minArgs]}, args[minArgs:].Strings(), nil

	case CMD_LOWER, "lower-case":
		return &handler{cmdLower, nil}, args.Strings(), nil

	case CMD_UPPER, "upper-case":
		return &handler{cmdUpper, nil}, args.Strings(), nil

	case CMD_STRIP_EXT:
		return &handler{cmdStripExt, nil}, args.Strings(), nil

	case CMD_IS_ABS:
		return &handler{cmdIsAbs, nil}, args.Strings(), nil

	case CMD_IS_LOCAL:
		return &handler{cmdIsLocal, nil}, args.Strings(), nil

	case CMD_IS_SAFE:
		return &handler{cmdIsSafe, nil}, args.Strings(), nil

	case CMD_IS_MATCH:
		return &handler{cmdIsMatch, args[:minArgs]}, args[minArgs:].Strings(), nil
	}

	return nil, nil, fmt.Errorf("Unknown command %q", cmd)
}

// executePipeHandlers executes all handlers in pipe with given data
func executePipeHandlers(p pipe, data string) (error, bool) {
	var err error
	var ok bool

	for _, cmd := range p {
		data, err, ok = cmd.Func(data, cmd.Args)

		if err != nil || !ok {
			return err, ok
		}
	}

	if data != "" {
		fmt.Printf("%s%s", data, separator)
	}

	return nil, true
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
	info := usage.NewInfo("", "args…")

	info.AppNameColorTag = colorTagApp

	info.AddCommand(CMD_BASENAME, "Strip directory and suffix from filenames", "?path…")
	info.AddCommand(CMD_DIRNAME, "Strip last component from file name", "?path…")
	info.AddCommand(CMD_DIRNAME_NUM, "Return N elements from path", "num", "?path…")
	info.AddCommand(CMD_READLINK, "Print resolved symbolic links or canonical file names", "?path…")
	info.AddCommand(CMD_CLEAN, "Print shortest path name equivalent to path by purely lexical processing", "?path…")
	info.AddCommand(CMD_COMPACT, "Converts path to compact representation", "?path…")
	info.AddCommand(CMD_ABS, "Print absolute representation of path", "?path…")
	info.AddCommand(CMD_EXT, "Print file extension", "?path…")
	info.AddCommand(CMD_MATCH, "Filter given path using pattern", "pattern", "?path…")
	info.AddCommand(CMD_JOIN, "Join path elements", "root", "?path…")

	info.AddCommand(CMD_ADD_PREFIX, "Add the substring at the beginning", "prefix", "?path…")
	info.AddCommand(CMD_DEL_PREFIX, "Remove the substring at the beginning", "prefix", "?path…")
	info.AddCommand(CMD_ADD_SUFFIX, "Add the substring at the end", "suffix", "?path…")
	info.AddCommand(CMD_DEL_SUFFIX, "Remove the substring at the end", "suffix", "?path…")
	info.AddCommand(CMD_EXCLUDE, "Exclude part of the path", "substr", "?path…")
	info.AddCommand(CMD_REPLACE, "Replace part of the path", "old", "new", "?path…")
	info.AddCommand(CMD_LOWER, "Convert path to lower case", "?path…")
	info.AddCommand(CMD_UPPER, "Convert path to upper case", "?path…")
	info.AddCommand(CMD_STRIP_EXT, "Remove file extension", "?path…")

	info.AddCommand(CMD_IS_ABS, "Check if given path is absolute", "?path…")
	info.AddCommand(CMD_IS_LOCAL, "Check if given path is local", "?path…")
	info.AddCommand(CMD_IS_SAFE, "Check if given path is safe", "?path…")
	info.AddCommand(CMD_IS_MATCH, "Check if given path is match to pattern", "pattern", "?path…")

	info.AddOption(OPT_ZERO, "End each output line with NUL, not newline")
	info.AddOption(OPT_SPACE, "End each output line with space, not newline")
	info.AddOption(OPT_QUIET, "Suppress all error messages")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")

	if withSelfUpdate {
		info.AddOption(OPT_UPDATE, "Update application to the latest version")
	}

	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddEnv("PATH_QUIET", "Flag to suppress all error messages {s-}(Boolean){!}")

	info.AddExample(
		"base /path/to/file.txt",
		"→ file.txt",
	)

	info.AddExample(
		"dir /path/to/file.txt",
		"→ /path/to",
	)

	info.AddExample(
		"compact /very/long/path/to/some/file.txt",
		"→ /v/l/p/t/s/file.txt",
	)

	info.AddExample(
		"path abs,strip-ext *",
		"Run many commands at once using piping",
	)

	info.AddRawExample(
		`find . -type f | path 'base,match+*.md,strip-ext'`,
		"Run many commands at once using piping with stdin data",
	)

	info.AddRawExample(
		"ls -1 | path is-match '*.txt' && echo MATCH!",
		"Check if all files in current directory is match to pattern",
	)

	info.AddRawExample(
		"PATH_QUIET=1 path dir /path/to/file.txt",
		"Run dir command in quiet mode enabled by environment variable",
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

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //
