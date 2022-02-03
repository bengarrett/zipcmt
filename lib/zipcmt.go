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
	"strings"
	"time"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/zipcmt/internal/misc"
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
	Zips      int
	Cmmts     int
	internal
}

type internal struct {
	test    bool
	log     string
	names   int         // nolint: structcheck
	saved   int         // nolint: structcheck
	exports misc.Export // nolint: structcheck
	hashes  hash        // nolint: structcheck
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
	hash map[[32]byte]bool
	save struct {
		name string
		src  string
		cmmt string
		mod  time.Time
		ow   bool
	}
)

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
func Read(name string, raw bool) (string, error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		return "", nil // nolint: nilerr
	}
	defer r.Close()
	cmmt := r.Comment
	if cmmt == "" {
		return "", nil
	}
	if strings.HasPrefix(cmmt, "TORRENTZIPPED-") {
		return "", nil
	}
	if strings.TrimSpace(cmmt) == "" {
		return "", nil
	}
	if !raw {
		b, err := convert.D437(cmmt)
		if err != nil {
			return "", fmt.Errorf("codepage 437 decoder: %w", err)
		}
		cmmt = string(b)
	}
	return cmmt, nil
}

// WalkDirs walks the directories provided by the Arg slice for zip archives to extract any found comments.
func (c *Config) WalkDirs() {
	c.init()
	// sanitize the export directory
	if err := c.Clean(); err != nil {
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
func (c *Config) WalkDir(root string) error { // nolint: cyclop,funlen,gocognit
	c.init()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				return nil
			}
			return err
		}
		// skip directories and non-zip files
		if d.IsDir() || !misc.Valid(d.Name()) {
			return nil
		}
		// skip sub-directories
		if c.NoWalk && filepath.Dir(path) != filepath.Dir(root) {
			return nil
		}
		c.Zips++
		if !c.test && !c.Print && !c.Quiet {
			fmt.Print("\r", color.Secondary.Sprint("Scanned "), color.Primary.Sprintf("%d zip archives", c.Zips))
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
		c.Cmmts++
		// print the comment
		fmt.Print(c.Separator(path))
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
			dat.name = misc.ExportName(path)
			if c.save(dat) {
				c.WriteLog("SAVED: " + dat.name + humanize.Bytes(uint64(len(cmmt))))
				c.saved++
			}
		}
		if c.SaveName != "" {
			dat.name = c.exports.Unique(path, c.SaveName)
			c.names += len(dat.name)
			if c.save(dat) {
				c.WriteLog(fmt.Sprintf("SAVED: %s (%s) << %s", dat.name, humanize.Bytes(uint64(len(cmmt))), path))
				c.saved++
			}
		}
		return err
	})
	if err != nil {
		return fmt.Errorf("walk directory: %s, %w", root, err)
	}
	return nil
}

// Clean the syntax of the target export directory path.
func (c *Config) Clean() error {
	name := c.SaveName
	if name == "" {
		return nil
	}
	name = filepath.Clean(name)
	p := strings.Split(name, string(filepath.Separator))
	if p[0] == "~" {
		hd, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("%s: export %w", name, err)
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
	return nil
}

// init initialise the Config maps.
func (c *Config) init() {
	if c.exports == nil {
		c.exports = make(misc.Export)
	}
	if c.hashes == nil {
		c.hashes = make(hash)
	}
}

// lastMod preserves the zip files last modification date.
func (c *Config) lastMod(file fs.DirEntry) time.Time {
	zero := time.Date(0o001, 1, 1, 0o0, 0o0, 0o0, 0o0, time.UTC)
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

// Separator prints and stylises the named file.
func (c *Config) Separator(name string) string {
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
			s := fmt.Sprintf("Saved %d comments from %d finds", c.saved, c.Cmmts)
			c.WriteLog(s)
		}
		s := fmt.Sprintf("Scan finished, time taken: %s", c.Timer())
		c.WriteLog(s)
	}
	if c.Quiet {
		return ""
	}
	a, cm, unq := "archive", "comment", ""
	if c.Zips != 1 {
		a += "s"
	}
	if c.Cmmts != 1 {
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
		color.Primary.Sprintf("%d zip %s", c.Zips, a)
	if c.SaveName != "" && c.saved != c.Cmmts {
		s += color.Secondary.Sprint(", saved ") +
			color.Primary.Sprintf("%d text files", c.saved)
	}
	s += color.Secondary.Sprint(" and found ") +
		color.Primary.Sprintf("%d %s%s", c.Cmmts, unq, cm)
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
		defer func() {
			if err := os.Chtimes(dat.name, time.Now(), dat.mod); err != nil {
				c.Error(fmt.Errorf("%s: %w", dat.name, err))
			}
		}()
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

// stdout prints the cmmt with an ANSI reset command.
func stdout(cmmt string) {
	const resetCmd = "\033[0m"
	fmt.Printf("%s%s\n", cmmt, resetCmd)
}
