package misc

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type (
	Export map[string]bool
)

const Filename = "-zipcomment.txt"

// unique checks the destination path against an export map.
// The map contains a unique collection of previously used destination
// paths, to avoid creating duplicate text filenames while using the
// SaveName config.
func (e Export) Unique(zipPath, dest string) string {
	base := filepath.Base(zipPath)
	name := strings.TrimSuffix(base, filepath.Ext(base)) + Filename
	if runtime.GOOS == "windows" {
		name = strings.ToLower(name)
	}
	path := filepath.Join(dest, name)
	if f := e.Find(path); f != path {
		path = f
	}
	e[path] = true
	return path
}

// Find searches for the name in the export map.
// If no matches exist, the name is unique and returned.
// Otherwise find attempts to append a `_1` suffix to the
// name. If the name already has this suffix, the digit
// is incrementally increased until a unique name is returned.
func (e Export) Find(name string) string {
	if !e[name] {
		return name
	}
	i, new := 0, ""
	for {
		i++
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		a := strings.Split(base, "_")
		n, err := strconv.Atoi(a[len(a)-1])
		if err == nil {
			i = n
			base = strings.Join(a[0:len(a)-1], "_")
		}
		suf := fmt.Sprintf("_%d", i+1)
		new = fmt.Sprintf("%s%s%s", base, suf, ext)
		if !e[new] {
			return new
		}
		if i > 9999 {
			break
		}
	}
	return ""
}

// ExportName returns a text file file path for the Export config.
func ExportName(path string) string {
	if path == "" {
		return ""
	}
	return strings.TrimSuffix(path, filepath.Ext(path)) + Filename
}

// Self returns the path for the zipcmt executable.
func Self() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("self error: %w", err)
	}
	return exe, nil
}

// Valid checks that the named file is a known zip archive.
func Valid(name string) bool {
	const z = ".zip"
	return filepath.Ext(strings.ToLower(name)) == z
}
