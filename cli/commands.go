package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/path"
	"github.com/essentialkaos/ek/v12/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// cmdBasename is handler for "base" command
func cmdBasename(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		result = append(result, path.Base(item))
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdDirname is handler for "dir" command
func cmdDirname(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		result = append(result, path.Dir(item))
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdDirnameNum is handler for "dirn" command
func cmdDirnameNum(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	num, err := strconv.Atoi(input[0])

	if err != nil {
		return err, false
	}

	for _, item := range input[1:] {
		result = append(result, path.DirN(item, num))
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdReadlink is handler for "link" command
func cmdReadlink(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		dest, _ := filepath.EvalSymlinks(item)

		if dest != "" {
			result = append(result, dest)
		}
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdClean is handler for "clean" command
func cmdClean(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		result = append(result, path.Clean(item))
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdCompact is handler for "compact" command
func cmdCompact(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		result = append(result, path.Compact(item))
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdExt is handler for "ext" command
func cmdExt(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		ext := path.Ext(item)

		if ext != "" {
			result = append(result, ext)
		}
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdAbs is handler for "abs" command
func cmdAbs(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	for _, item := range input {
		dest, _ := filepath.Abs(item)
		result = append(result, dest)
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdMatch is handler for "match" command
func cmdMatch(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	if len(input) < 2 {
		printError("Not enough arguments")
		return nil, false
	}

	pattern := input[0]

	for _, item := range input[1:] {
		isMatch, _ := filepath.Match(pattern, item)

		if !isMatch {
			continue
		}

		result = append(result, item)
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdJoin is handler for "join" command
func cmdJoin(args options.Arguments) (error, bool) {
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	if len(input) < 2 {
		printError("Not enough arguments")
		return nil, false
	}

	root := input[0]
	path, err := path.JoinSecure(root, input[1:]...)

	if err != nil {
		return err, false
	}

	fmt.Println(path)

	return nil, true
}

// cmdAddPrefix is handler for "add-prefix" command
func cmdAddPrefix(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	if len(input) < 2 {
		printError("Not enough arguments")
		return nil, false
	}

	prefix := input[0]

	for _, item := range input[1:] {
		result = append(result, prefix+item)
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdDelPrefix is handler for "del-prefix" command
func cmdDelPrefix(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	if len(input) < 2 {
		printError("Not enough arguments")
		return nil, false
	}

	prefix := input[0]

	for _, item := range input[1:] {
		data, _ := strings.CutPrefix(item, prefix)
		result = append(result, data)
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdAddSuffix is handler for "add-suffix" command
func cmdAddSuffix(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	if len(input) < 2 {
		printError("Not enough arguments")
		return nil, false
	}

	suffix := input[0]

	for _, item := range input[1:] {
		result = append(result, item+suffix)
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdDelSuffix is handler for "del-suffix" command
func cmdDelSuffix(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	if len(input) < 2 {
		printError("Not enough arguments")
		return nil, false
	}

	prefix := input[0]

	for _, item := range input[1:] {
		data, _ := strings.CutSuffix(item, prefix)
		result = append(result, data)
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

	return nil, true
}

// cmdExclude is handler for "exclude" command
func cmdExclude(args options.Arguments) (error, bool) {
	var result []string

	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	if len(input) < 2 {
		printError("Not enough arguments")
		return nil, false
	}

	substr := input[0]

	for _, item := range input[1:] {
		result = append(result, strutil.Exclude(item, substr))
	}

	if len(result) == 0 {
		return err, false
	}

	fmt.Println(strings.Join(result, getSeparator()))

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
	input, err := getInputData(args)

	if err != nil {
		return err, false
	}

	if len(input) < 2 {
		printError("Not enough arguments")
		return nil, false
	}

	pattern := input[0]

	for _, item := range input[1:] {
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
	var result []string

	data := strings.Split(args.Flatten(), " ")

	if !fsutil.IsCharacterDevice("/dev/stdin") {
		stdinData, err := io.ReadAll(os.Stdin)

		if err != nil {
			return nil, fmt.Errorf("Can't data from standard input: %v", err)
		}

		data = append(data, strings.Split(strings.ReplaceAll(string(stdinData), "\n", " "), " ")...)
	}

	for _, item := range data {
		if strings.Trim(item, " \r\n") == "" {
			continue
		}

		result = append(result, item)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Input is empty")
	}

	return result, nil
}

// getSeparator returns data separator
func getSeparator() string {
	switch {
	case options.GetB(OPT_SPACE):
		return " "
	case options.GetB(OPT_ZERO):
		return "\x00"
	default:
		return "\n"
	}
}
