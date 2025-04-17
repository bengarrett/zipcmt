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

	"github.com/bengarrett/retrotxtgo/byter"
	"github.com/bengarrett/sauce"
	"github.com/bengarrett/zipcmt/internal/misc"
	humanize "github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"golang.org/x/text/encoding/charmap"
)

// Config zipcmt to walk one or more directories.
type Config struct {
	Dirs      []string // Dirs are the directory paths to walk.
	SaveName  string   // SaveName is an optional directory path to save any found comments as uniquely named text files.
	Dupes     bool     // Dupes shows all comments, including duplicates found in multiple zips.
	Export    bool     // Export the comments as text files stored alongside the source zip files.
	Log       bool     // Log creates a logfile for debugging.
	Overwrite bool     // Overwrite any previously exported comment text files.
	// Now ignores the zip files last modification date,
	// which is otherwise applied to the comment text file.
	Now    bool
	NoWalk bool // NoWalk ignores all subdirectories while scanning for zip archives.
	Raw    bool // Raw uses the original comment text encoding (CP437, ISO-8859...) instead of Unicode.
	Print  bool // Print found comments to stdout.
	Quiet  bool // Quiet suppresses the scan activity feedback to stdout.
	Zips   int  // Zips is the number of zip files scanned.
	Cmmts  int  // Cmmts are the number of zip comments found.
	internal
}

type internal struct {
	test    bool
	log     string
	names   uint
	saved   int
	exports misc.Export
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
	ErrFlag     = errors.New("this option is used after a directory, it must be placed before any directories are listed")
	ErrDirExist = errors.New("directory does not exist")
	ErrIsFile   = errors.New("directory is a file")
	ErrMissing  = errors.New("directory cannot be found")
	ErrPath     = errors.New("directory path cannot be found or points to a file")
	ErrPerm     = errors.New("directory access is blocked due to its permissions")
	ErrRead     = errors.New("skip named zip file due to read error")
	ErrValid    = errors.New("the operating system reports this directory is invalid")
)

// Read the named zip file and return the zip comment.
// The Raw config will return the comment in its original legacy encoding.
// Otherwise the comment is returned as Unicode text.
func Read(name string, raw bool) (string, error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		return "", ErrRead
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

	if raw {
		return cmmt, nil
	}
	p := []byte(cmmt)
	if ok := sauce.Contains(p); ok {
		cmmt = string(sauce.Trim(p))
	}
	b, err := byter.Decode(charmap.CodePage437, cmmt)
	if err != nil {
		return "", fmt.Errorf("codepage 437 decoder: %w", err)
	}
	return string(b), nil
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
		_ = c.WalkDir(root)
	}
}

// WalkDir walks the root directory for zip archives and to extract any found comments.
// The returned error is only used for testing purposes.
func (c *Config) WalkDir(root string) error { //nolint: cyclop,funlen,gocognit
	c.init()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				// skip permission errors for subdirectories
				// but return an error if the root is inaccessible
				if root != path {
					return nil
				}
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
			fmt.Fprint(os.Stdout, "\r", color.Secondary.Sprint("Scanned "),
				color.Primary.Sprintf("%d zip archives", c.Zips))
		}
		// read zip file comment
		cmmt, err := Read(path, c.Raw)
		if err != nil {
			if !errors.Is(err, ErrRead) {
				c.Error(err)
			}
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
		fmt.Fprint(os.Stdout, c.Separator(path))
		if c.Print {
			stdout(cmmt)
		}
		// save the comment to a text file
		dat := save{
			name: "",
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
			c.names += uint(len(dat.name))
			if c.save(dat) {
				c.WriteLog(fmt.Sprintf("SAVED: %s (%s) << %s",
					dat.name, humanize.Bytes(uint64(len(cmmt))), path))
				c.saved++
			}
		}
		return err
	})
	if errs := walkErrs(root, err); errs != nil {
		color.Error.Tips(fmt.Sprint(errs))
	}
	if err != nil {
		return fmt.Errorf("walk dir %w: %s", err, root)
	}
	return nil
}

func walkErrs(root string, err error) error {
	var pathError *os.PathError
	if errors.As(err, &pathError) {
		if root != "" && root[:1] == "-" {
			return fmt.Errorf("%w: %s", ErrFlag, root)
		}
	}
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %s", ErrDirExist, root)
	}
	if errors.Is(err, fs.ErrPermission) {
		f, err := os.Stat(root)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrPerm, root)
		}
		return fmt.Errorf("%w: %s, %s", ErrPerm, f.Mode(), root)
	}
	if err != nil {
		return fmt.Errorf("walk directory: %s, %w %T", root, err, err.Error())
	}
	return nil
}

// Clean the syntax and check the usability of the SaveName directory path.
func (c *Config) Clean() error {
	name := c.SaveName
	if name == "" {
		return nil
	}
	name = filepath.Clean(name)
	before, _, _ := strings.Cut(name, string(filepath.Separator))
	if before == "~" {
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
func (c *Config) Status() string {
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

// SaveName a zip cmmt to the file path.
// Unless the overwrite argument is set, any previous cmmt text files are skipped.
func (c *Config) save(dat save) bool {
	// name, cmmt string, mod time.Time, ow bool
	if dat.cmmt == "" {
		return false
	}
	if !dat.ow {
		if s, err := os.Stat(dat.name); err == nil {
			size := humanize.Bytes(uint64(s.Size())) //nolint:gosec
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
	b := []byte(dat.cmmt)
	i, err := f.Write(b)
	if err != nil {
		c.Error(fmt.Errorf("%s: %w", dat.name, err))
		return true
	}
	if i == 0 {
		if err1 := os.Remove(dat.name); err1 != nil {
			c.Error(fmt.Errorf("%s: %w", dat.name, err1))
		}
	}
	return true
}

// stdout prints the cmmt with an ANSI reset command.
func stdout(cmmt string) {
	const resetCmd = "\033[0m"
	fmt.Fprintf(os.Stdout, "%s%s\n", cmmt, resetCmd)
}
