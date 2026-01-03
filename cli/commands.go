package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/path"
	"github.com/essentialkaos/ek/v13/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// cmdBasename is handler for "base" command
func cmdBasename(data string, args options.Arguments) (string, error, bool) {
	return path.Base(data), nil, true
}

// cmdDirname is handler for "dir" command
func cmdDirname(data string, args options.Arguments) (string, error, bool) {
	return path.Dir(data), nil, true
}

// cmdDirnameNum is handler for "dirn" command
func cmdDirnameNum(data string, args options.Arguments) (string, error, bool) {
	num, err := strconv.Atoi(strings.ReplaceAll(args.Get(0).String(), "^", "-"))

	if err != nil {
		return "", fmt.Errorf("Can't parse number of directories: %v", err), false
	}

	return path.DirN(data, num), nil, true
}

// cmdReadlink is handler for "link" command
func cmdReadlink(data string, args options.Arguments) (string, error, bool) {
	dest, _ := filepath.EvalSymlinks(data)
	return strutil.B(dest != "", dest, data), nil, true
}

// cmdClean is handler for "clean" command
func cmdClean(data string, args options.Arguments) (string, error, bool) {
	return path.Clean(data), nil, true
}

// cmdCompact is handler for "compact" command
func cmdCompact(data string, args options.Arguments) (string, error, bool) {
	return path.Compact(data), nil, true
}

// cmdExt is handler for "ext" command
func cmdExt(data string, args options.Arguments) (string, error, bool) {
	return path.Ext(data), nil, true
}

// cmdAbs is handler for "abs" command
func cmdAbs(data string, args options.Arguments) (string, error, bool) {
	dest, _ := filepath.Abs(data)
	return strutil.B(dest != "", dest, data), nil, true
}

// cmdMatch is handler for "match" command
func cmdMatch(data string, args options.Arguments) (string, error, bool) {
	isMatch, _ := filepath.Match(args.Get(0).String(), data)
	return strutil.B(isMatch, data, ""), nil, true
}

// cmdJoin is handler for "join" command
func cmdJoin(data string, args options.Arguments) (string, error, bool) {
	path, err := path.JoinSecure(args.Get(0).String(), data)

	if err != nil {
		return path, err, false
	}

	return path, nil, true
}

// cmdAddPrefix is handler for "add-prefix" command
func cmdAddPrefix(data string, args options.Arguments) (string, error, bool) {
	return args.Get(0).String() + data, nil, true
}

// cmdDelPrefix is handler for "del-prefix" command
func cmdDelPrefix(data string, args options.Arguments) (string, error, bool) {
	data, _ = strings.CutPrefix(data, args.Get(0).String())
	return data, nil, true
}

// cmdAddSuffix is handler for "add-suffix" command
func cmdAddSuffix(data string, args options.Arguments) (string, error, bool) {
	return data + args.Get(0).String(), nil, true
}

// cmdDelSuffix is handler for "del-suffix" command
func cmdDelSuffix(data string, args options.Arguments) (string, error, bool) {
	data, _ = strings.CutSuffix(data, args.Get(0).String())
	return data, nil, true
}

// cmdExclude is handler for "exclude" command
func cmdExclude(data string, args options.Arguments) (string, error, bool) {
	return strutil.Exclude(data, args.Get(0).String()), nil, true
}

// cmdReplace is handler for "replace" command
func cmdReplace(data string, args options.Arguments) (string, error, bool) {
	return strings.ReplaceAll(
		data, args.Get(0).String(), args.Get(1).String(),
	), nil, true
}

// cmdLower is handler for "lower" command
func cmdLower(data string, args options.Arguments) (string, error, bool) {
	return strings.ToLower(data), nil, true
}

// cmdUpper is handler for "upper" command
func cmdUpper(data string, args options.Arguments) (string, error, bool) {
	return strings.ToUpper(data), nil, true
}

// cmdStripExt is handler for "strip-ext" command
func cmdStripExt(data string, args options.Arguments) (string, error, bool) {
	ext := path.Ext(data)

	if ext == "" {
		return data, nil, true
	}

	return strings.TrimSuffix(data, ext), nil, true
}

// cmdIsAbs is handler for "is-abs" command
func cmdIsAbs(data string, args options.Arguments) (string, error, bool) {
	return "", nil, filepath.IsAbs(data)
}

// cmdIsLocal is handler for "is-local" command
func cmdIsLocal(data string, args options.Arguments) (string, error, bool) {
	return "", nil, filepath.IsLocal(data)
}

// cmdIsSafe is handler for "is-safe" command
func cmdIsSafe(data string, args options.Arguments) (string, error, bool) {
	return "", nil, path.IsSafe(data)
}

// cmdIsMatch is handler for "is-match" command
func cmdIsMatch(data string, args options.Arguments) (string, error, bool) {
	isMatch, _ := filepath.Match(args.Get(0).String(), data)
	return "", nil, isMatch
}
