// © Ben Garrett https://github.com/bengarrett/zipcmt

// Package zipcmt is a viewer and an extractor of zip archive comments.
package zipcmt

import (
	"archive/zip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/gookit/color"
)

type Config struct {
	Timer     time.Time
	Save      string
	Export    bool
	Dupes     bool
	Overwrite bool
	Raw       bool
	Print     bool
	Quiet     bool
	zips      int
	cmmts     int
}

type (
	export map[string]bool
	hash   map[[32]byte]bool
)

const filename = "-zipcomment.txt"

var (
	ErrIsFile  = errors.New("directory is a file")
	ErrMissing = errors.New("directory cannot be found")
	ErrPath    = errors.New("directory path cannot be found or points to a file")
	ErrPerm    = errors.New("directory access is blocked due to its permissions")
	ErrValid   = errors.New("the operating system reports this directory is invalid")
)

// Clean the syntax of the target export directory path.
func (c *Config) Clean() error {
	if c.Save != "" {
		c.Save = filepath.Clean(c.Save)
		p := strings.Split(c.Save, string(filepath.Separator))
		if p[0] == "~" {
			hd, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			c.Save = strings.Replace(c.Save, "~", hd, 1)
		}
		s, err := os.Stat(c.Save)
		if errors.Is(err, fs.ErrInvalid) {
			return fmt.Errorf("%s: export %w", c.Save, ErrValid)
		}
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("%s: export %w", c.Save, ErrMissing)
		}
		if errors.Is(err, fs.ErrPermission) {
			return fmt.Errorf("%s: export %w", c.Save, ErrPerm)
		}
		if err != nil {
			return fmt.Errorf("%s: export %w", c.Save, err)
		}
		if !s.IsDir() {
			return fmt.Errorf("%s: export %w", c.Save, ErrIsFile)
		}
	}
	return nil
}

// Read the named zip file and return the zip comment.
// The Raw config will return the comment in its original legacy encoding.
// Otherwise the comment is returned as Unicode text.
func (c Config) Read(name string) (cmmt string, err error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		return "", nil
	}
	defer r.Close()
	if cmmt := r.Comment; cmmt != "" {
		if strings.HasPrefix(cmmt, "TORRENTZIPPED-") {
			return "", nil
		}
		if strings.TrimSpace(cmmt) == "" {
			return "", nil
		}
		if !c.Raw {
			b, err := convert.D437(cmmt)
			if err != nil {
				return "", err
			}
			cmmt = string(b)
		}
		return cmmt, nil
	}
	return "", nil
}

// Scan the root directory for zip archives and parse any found comments.
func (c *Config) Scan(root string) error {
	exports, hashes := export{}, hash{}
	files, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, file := range files {
		path := filepath.Join(root, file.Name())
		if !valid(path) {
			continue
		}
		c.zips++
		cmmt, err := c.Read(path)
		if err != nil {
			color.Error.Tips(fmt.Sprint(err))
			continue
		}
		if cmmt == "" {
			continue
		}
		if !c.Dupes {
			hash := sha256.Sum256([]byte(strings.TrimSpace(cmmt)))
			if hashes[hash] {
				continue
			}
			hashes[hash] = true
		}
		c.cmmts++
		fmt.Print(c.separator(path))
		if c.Print {
			stdout(cmmt)
		}
		if c.Export {
			save(exportName(path), cmmt, c.Overwrite)
		}
		if c.Save != "" {
			path = exports.unique(path, c.Save)
			save(path, cmmt, c.Overwrite)
		}
	}
	return nil
}

// Walk the root directory plus all subdirectories for zip archives and parse any found comments.
func (c *Config) Walk(root string) error {
	exports, hashes := export{}, hash{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				return nil
			}
			return err
		}
		if !valid(d.Name()) {
			return nil
		}
		c.zips++
		cmmt, err := c.Read(path)
		if err != nil {
			color.Error.Tips(fmt.Sprint(err))
			return nil
		}
		if cmmt == "" {
			return nil
		}
		if !c.Dupes {
			hash := sha256.Sum256([]byte(strings.TrimSpace(cmmt)))
			if hashes[hash] {
				return nil
			}
			hashes[hash] = true
		}
		c.cmmts++
		fmt.Print(c.separator(path))
		if c.Print {
			stdout(cmmt)
		}
		if c.Export {
			save(exportName(path), cmmt, c.Overwrite)
		}
		if c.Save != "" {
			path = exports.unique(path, c.Save)
			save(path, cmmt, c.Overwrite)
		}
		return err
	})
	return err
}

// separator prints and stylises the named file.
func (c Config) separator(name string) string {
	if !c.Print || c.Quiet {
		return ""
	}
	const fileID = 45
	const pointer = " \u2500\u2500 "
	if h, err := os.UserHomeDir(); err == nil {
		if len(name) > len(h) && name[0:len(h)] == h {
			name = strings.Replace(name, h, "~", 1)
		}
	}
	l := len(pointer) + len(name)
	if l >= fileID {
		return fmt.Sprintf("%s%s\n", pointer, name)
	}
	return fmt.Sprintf("\n%s%s %s\u2510\n", pointer, name, strings.Repeat("\u2500", fileID-l))
}

// Status summarizes the zip files scan.
func (c Config) Status() string {
	if c.Quiet {
		return ""
	}
	a, cm, unq := "archive", "comment", ""
	if c.zips != 1 {
		a += "s"
	}
	if c.cmmts != 1 {
		cm += "s"
	}
	if !c.Dupes {
		unq = "unique "
	}
	s := ""
	if c.Print {
		s = "\n"
	}
	return s + color.Secondary.Sprint("Scanned ") +
		color.Primary.Sprintf("%d zip %s", c.zips, a) +
		color.Secondary.Sprint(" and found ") +
		color.Primary.Sprintf("%d %s%s", c.cmmts, unq, cm) +
		color.Secondary.Sprint(", taking ") +
		color.Primary.Sprintf("%s", time.Since(c.Timer))
}

// Save a zip cmmt to the file path.
// Unless the overwrite argument is set, any previous cmmt textfiles are skipped.
func save(name, cmmt string, ow bool) bool {
	if cmmt == "" {
		return false
	}
	if !ow {
		if s, err := os.Stat(name); err == nil {
			color.Info.Tips(fmt.Sprintf("export skipped, file already exists: %s (%dB)\n", name, s.Size()))
			return false
		}
	}
	f, err := os.Create(name)
	if err != nil {
		color.Error.Tips(fmt.Sprint(fmt.Errorf("%s: %w", name, err)))
	}
	defer f.Close()
	if cmmt[len(cmmt)-1:] != "\n" {
		cmmt += "\n"
	}
	if i, err := f.Write([]byte(cmmt)); err != nil {
		color.Error.Tips(fmt.Sprint(fmt.Errorf("%s: %w", name, err)))
	} else if i == 0 {
		if err1 := os.Remove(name); err1 != nil {
			color.Error.Tips(fmt.Sprint(fmt.Errorf("%s: %w", name, err1)))
		}
	}
	return true
}

// exportName returns a textfile filepath for the Export config.
func exportName(path string) string {
	if path == "" {
		return ""
	}
	return strings.TrimSuffix(path, filepath.Ext(path)) + filename
}

// stdout prints the cmmt with an ANSI reset command.
func stdout(cmmt string) {
	const resetCmd = "\033[0m"
	fmt.Printf("%s%s\n", cmmt, resetCmd)
}

// valid checks that the named file is a known zip archive.
func valid(name string) bool {
	const z = ".zip"
	return filepath.Ext(strings.ToLower(name)) == z
}

// unique checks the destination path against an export map.
// The map contains a unique collection of previously used destination
// paths, to avoid creating duplicate text filenames while using the
// Save config.
func (exports export) unique(zipPath, dest string) string {
	base := filepath.Base(zipPath)
	name := strings.TrimSuffix(base, filepath.Ext(base)) + filename
	path := filepath.Join(dest, name)
	if f := exports.find(path); f != path {
		path = f
	}
	exports[path] = true
	return path
}

// find searches for the name in the export map.
// If no matches exist, the name is unique and returned.
// Otherwise find attempts to append a `_1` suffix to the
// name. If the name already has this suffix, the digit
// is incrementally increased until a unique name is returned.
func (e export) find(name string) string {
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
