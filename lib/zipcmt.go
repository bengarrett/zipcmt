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
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	humanize "github.com/dustin/go-humanize"
	"github.com/gookit/color"
)

type Config struct {
	Dirs      []string
	SaveName  string
	Dupes     bool
	Export    bool
	Log       bool
	Overwrite bool
	Now       bool
	NoWalk    bool
	Raw       bool
	Print     bool
	Quiet     bool
	internal
}

type internal struct {
	test    bool
	log     string
	zips    int
	cmmts   int
	names   int
	saved   int
	exports export
	hashes  hash
	timer   time.Time
}

// SetLog sets the full path to a new log file with a name based on the current date and time.
func (i *internal) SetLog() {
	i.log = logName()
}

// SetTest toggles the unit test mode flag.
func (i *internal) SetTest() {
	i.test = true
}

// SetTimer initializes a timer for process time.
func (i *internal) SetTimer() {
	i.timer = time.Now()
}

// SetLog returns the full path to the log file.
func (i *internal) LogName() string {
	return i.log
}

// Timer returns the time since the SetTimer was triggered.
func (i *internal) Timer() time.Duration {
	return time.Since(i.timer)
}

type (
	export map[string]bool
	hash   map[[32]byte]bool
	save   struct {
		name string
		src  string
		cmmt string
		mod  time.Time
		ow   bool
	}
)

const filename = "-zipcomment.txt"

var (
	ErrIsFile  = errors.New("directory is a file")
	ErrMissing = errors.New("directory cannot be found")
	ErrPath    = errors.New("directory path cannot be found or points to a file")
	ErrPerm    = errors.New("directory access is blocked due to its permissions")
	ErrValid   = errors.New("the operating system reports this directory is invalid")
)

// Read the named zip file and return the zip comment.
// The Raw config will return the comment in its original legacy encoding.
// Otherwise the comment is returned as Unicode text.
func Read(name string, raw bool) (cmmt string, err error) {
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
		if !raw {
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

// WalkDirs walks the directories provided by the Arg slice for zip archives to extract any found comments.
func (c *Config) WalkDirs() {
	c.init()
	// sanitize the export directory
	if err := c.clean(); err != nil {
		c.Error(err)
	}
	// walk through the directories provided
	for _, root := range c.Dirs {
		if err := c.WalkDir(root); err != nil {
			c.Error(err)
		}
	}
}

// WalkDir walks the root directory for zip archives and to extract any found comments.
func (c *Config) WalkDir(root string) error {
	c.init()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				return nil
			}
			return err
		}
		// skip directories and non-zip files
		if d.IsDir() || !valid(d.Name()) {
			return nil
		}
		// skip sub-directories
		if c.NoWalk && filepath.Dir(path) != filepath.Dir(root) {
			return nil
		}
		c.zips++
		if !c.test && !c.Print && !c.Quiet {
			fmt.Print("\r", color.Secondary.Sprint("Scanned "), color.Primary.Sprintf("%d zip archives", c.zips))
		}
		// read zip file comment
		cmmt, err := Read(path, c.Raw)
		if err != nil {
			c.Error(err)
			return nil
		}
		if cmmt == "" {
			return nil
		}
		// hash the comment
		if !c.Dupes {
			hash := sha256.Sum256([]byte(strings.TrimSpace(cmmt)))
			if c.hashes[hash] {
				return nil
			}
			c.hashes[hash] = true
		}
		c.cmmts++
		// print the comment
		fmt.Print(c.separator(path))
		if c.Print {
			stdout(cmmt)
		}
		// save the comment to a text file
		dat := save{
			src:  path,
			cmmt: cmmt,
			mod:  c.lastMod(d),
			ow:   c.Overwrite,
		}
		if c.Export {
			dat.name = exportName(path)
			if c.save(dat) {
				c.WriteLog("SAVED: " + dat.name + humanize.Bytes(uint64(len(cmmt))))
				c.saved++
			}
		}
		if c.SaveName != "" {
			dat.name = c.exports.unique(path, c.SaveName)
			c.names += len(dat.name)
			if c.save(dat) {
				c.WriteLog(fmt.Sprintf("SAVED: %s (%s) << %s", dat.name, humanize.Bytes(uint64(len(cmmt))), path))
				c.saved++
			}
		}
		return err
	})
	return err
}

// clean the syntax of the target export directory path.
func (c *Config) clean() error {
	if name := c.SaveName; name != "" {
		name = filepath.Clean(name)
		p := strings.Split(name, string(filepath.Separator))
		if p[0] == "~" {
			hd, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			name = strings.Replace(name, "~", hd, 1)
		}
		s, err := os.Stat(name)
		if errors.Is(err, fs.ErrInvalid) {
			return fmt.Errorf("%s: export %w", name, ErrValid)
		}
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("%s: export %w", name, ErrMissing)
		}
		if errors.Is(err, fs.ErrPermission) {
			return fmt.Errorf("%s: export %w", name, ErrPerm)
		}
		if err != nil {
			return fmt.Errorf("%s: export %w", name, err)
		}
		if !s.IsDir() {
			return fmt.Errorf("%s: export %w", name, ErrIsFile)
		}
		c.SaveName = name
	}
	return nil
}

// init initialise the Config maps.
func (c *Config) init() {
	if c.exports == nil {
		c.exports = make(export)
	}
	if c.hashes == nil {
		c.hashes = make(hash)
	}
}

// lastMod preserves the zip files last modification date.
func (c *Config) lastMod(file fs.DirEntry) time.Time {
	zero := time.Date(0001, 1, 1, 00, 00, 00, 00, time.UTC)
	if c.Now {
		return zero
	}
	i, err := file.Info()
	if err != nil {
		c.Error(err)
		return zero
	}
	return i.ModTime()
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
	if c.Log {
		if c.SaveName != "" {
			s := fmt.Sprintf("Saved %d comments from %d finds", c.saved, c.cmmts)
			c.WriteLog(s)
		}
		s := fmt.Sprintf("Scan finished, time taken: %s", c.Timer())
		c.WriteLog(s)
	}
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
	s := "\n"
	if !c.test && !c.Print {
		s = "\r"
	}
	s += color.Secondary.Sprint("Scanned ") +
		color.Primary.Sprintf("%d zip %s", c.zips, a)
	if c.SaveName != "" && c.saved != c.cmmts {
		s += color.Secondary.Sprint(", saved ") +
			color.Primary.Sprintf("%d text files", c.saved)
	}
	s += color.Secondary.Sprint(" and found ") +
		color.Primary.Sprintf("%d %s%s", c.cmmts, unq, cm)
	if !c.test {
		s += color.Secondary.Sprint(", taking ") +
			color.Primary.Sprintf("%s", c.Timer()) + "\n"
	}
	return s
}

// SaveName a zip cmmt to the file path.
// Unless the overwrite argument is set, any previous cmmt text files are skipped.
func (c *Config) save(dat save) bool {
	// name, cmmt string, mod time.Time, ow bool
	if dat.cmmt == "" {
		return false
	}
	if !dat.ow {
		if s, err := os.Stat(dat.name); err == nil {
			size := humanize.Bytes(uint64(s.Size()))
			info := fmt.Sprintf("export skipped, file already exists: %s (%s)", dat.name, size)
			color.Info.Tips(info)
			c.WriteLog(fmt.Sprintf("SKIP (exists): %s (%s)", dat.name, size))
			return false
		}
	}
	f, err := os.Create(dat.name)
	if err != nil {
		c.Error(fmt.Errorf("%s: %w", dat.name, err))
	}
	defer f.Close()
	if !dat.mod.IsZero() {
		defer os.Chtimes(dat.name, time.Now(), dat.mod)
	}
	if dat.cmmt[len(dat.cmmt)-1:] != "\n" {
		dat.cmmt += "\n"
	}
	if i, err := f.Write([]byte(dat.cmmt)); err != nil {
		c.Error(fmt.Errorf("%s: %w", dat.name, err))
	} else if i == 0 {
		if err1 := os.Remove(dat.name); err1 != nil {
			c.Error(fmt.Errorf("%s: %w", dat.name, err1))
		}
	}
	return true
}

// exportName returns a text file file path for the Export config.
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
// SaveName config.
func (e export) unique(zipPath, dest string) string {
	base := filepath.Base(zipPath)
	name := strings.TrimSuffix(base, filepath.Ext(base)) + filename
	if runtime.GOOS == "windows" {
		name = strings.ToLower(name)
	}
	path := filepath.Join(dest, name)
	if f := e.find(path); f != path {
		path = f
	}
	e[path] = true
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
