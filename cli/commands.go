package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/path"
	"github.com/essentialkaos/ek/v12/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// cmdBasename is handler for "base" command
func cmdBasename(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for i, item := range input {
		fmt.Printf("%s", path.Base(item))

		if i+1 < len(input) {
			printSeparator()
		}
	}

	fmt.Println()

	return nil, true
}

// cmdDirname is handler for "dir" command
func cmdDirname(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for i, item := range input {
		fmt.Printf("%s", path.Dir(item))

		if i+1 < len(input) {
			printSeparator()
		}
	}

	fmt.Println()

	return nil, true
}

// cmdReadlink is handler for "link" command
func cmdReadlink(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for i, item := range input {
		dest, _ := filepath.EvalSymlinks(item)
		fmt.Printf("%s", dest)

		if i+1 < len(input) {
			printSeparator()
		}
	}

	fmt.Println()

	return nil, true
}

// cmdClean is handler for "clean" command
func cmdClean(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for i, item := range input {
		fmt.Printf("%s", path.Clean(item))

		if i+1 < len(input) {
			printSeparator()
		}
	}

	fmt.Println()

	return nil, true
}

// cmdCompact is handler for "compact" command
func cmdCompact(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for i, item := range input {
		fmt.Printf("%s", path.Compact(item))

		if i+1 < len(input) {
			printSeparator()
		}
	}

	fmt.Println()

	return nil, true
}

// cmdExt is handler for "ext" command
func cmdExt(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for i, item := range input {
		fmt.Printf("%s", path.Ext(item))

		if i+1 < len(input) {
			printSeparator()
		}
	}

	fmt.Println()

	return nil, true
}

// cmdAbs is handler for "abs" command
func cmdAbs(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for i, item := range input {
		dest, _ := filepath.Abs(item)
		fmt.Printf("%s", dest)

		if i+1 < len(input) {
			printSeparator()
		}
	}

	fmt.Println()

	return nil, true
}

// cmdVolume is handler for "volume" command
func cmdVolume(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for i, item := range input {
		fmt.Printf("%s", filepath.VolumeName(item))

		if i+1 < len(input) {
			printSeparator()
		}
	}

	fmt.Println()

	return nil, true
}

// cmdMatch is handler for "match" command
func cmdMatch(args options.Arguments) (error, bool) {
	pattern := args.Get(0).String()
	input, err := getInputData(args[1:])

	if err != nil {
		return err, false
	}

	isDataPrinted := false

	for i, item := range input {
		isMatch, _ := filepath.Match(pattern, item)

		if !isMatch {
			continue
		}

		isDataPrinted = true
		fmt.Printf("%s", item)

		if i+1 < len(input) {
			printSeparator()
		}
	}

	if !isDataPrinted {
		return nil, false
	}

	fmt.Println()

	return nil, true
}

// cmdIsAbs is handler for "is-abs" command
func cmdIsAbs(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		if !filepath.IsAbs(item) {
			return nil, false
		}
	}

	return nil, true
}

// cmdIsLocal is handler for "is-local" command
func cmdIsLocal(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		if !filepath.IsLocal(item) {
			return nil, false
		}
	}

	return nil, true
}

// cmdIsSafe is handler for "is-safe" command
func cmdIsSafe(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		if !path.IsSafe(item) {
			return nil, false
		}
	}

	return nil, true
}

// cmdIsMatch is handler for "is-match" command
func cmdIsMatch(args options.Arguments) (error, bool) {
	pattern := args.Get(0).String()
	input, err := getInputData(args[1:])

	if err != nil {
		return err, false
	}

	for _, item := range input {
		isMatch, _ := filepath.Match(pattern, item)

		if !isMatch {
			return nil, false
		}
	}

	return nil, true
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getInputData returns import from stdin, arguments or both
func getInputData(args options.Arguments) ([]string, error) {
	var rawData string

	if !fsutil.IsCharacterDevice("/dev/stdin") {
		stdinData, err := io.ReadAll(os.Stdin)

		if err != nil {
			return nil, fmt.Errorf("Can't read stdin data: %v", err)
		}

		rawData = strings.ReplaceAll(string(stdinData), "\n", " ")
	}

	if len(args) > 0 {
		if rawData == "" {
			rawData = args.Flatten()
		} else {
			rawData += " " + args.Flatten()
		}
	}

	if rawData == "" {
		return nil, fmt.Errorf("Input is empty")
	}

	return strutil.Fields(rawData), nil
}

// printSeparator prints data separator
func printSeparator() {
	switch {
	case options.GetB(OPT_SPACE):
		fmt.Printf(" ")
	case options.GetB(OPT_ZERO):
		fmt.Printf("\x00")
	default:
		fmt.Println()
	}
}
