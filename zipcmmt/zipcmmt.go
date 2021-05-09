// Package zipcmmmt is a viewer and an extractor of zip archive comments
package zipcmmt

import (
	"archive/zip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/convert"
)

type Config struct {
	ExportDir  string
	ExportFile bool
	NoDupes    bool
	Overwrite  bool
	Raw        bool
	Print      bool
	Quiet      bool
	zips       int
	cmmts      int
}

const filename = "-zipcomment.txt"

var ErrExpFile = errors.New("export directory is a file")

// Clean the syntax of the target export directory path.
func (c *Config) Clean() error {
	if c.ExportDir != "" {
		c.ExportDir = filepath.Clean(c.ExportDir)
		p := strings.Split(c.ExportDir, string(filepath.Separator))
		if p[0] == "~" {
			hd, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			c.ExportDir = strings.Replace(c.ExportDir, "~", hd, 1)
		}
		s, err := os.Stat(c.ExportDir)
		if err != nil {
			return err
		}
		if !s.IsDir() {
			return ErrExpFile
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
	hashes := map[[32]byte]bool{}
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
			fmt.Println(err)
			continue
		}
		if cmmt == "" {
			continue
		}
		if c.NoDupes {
			hash := sha256.Sum256([]byte(strings.TrimSpace(cmmt)))
			if hashes[hash] {
				continue
			}
			hashes[hash] = true
		}
		c.cmmts++
		fmt.Print(c.Separator(path))
		if c.Print {
			stdout(cmmt)
		}
		if c.ExportFile {
			go save(path, cmmt, c.Overwrite)
		}
		if c.ExportDir != "" {
			go save(filepath.Join(c.ExportDir, file.Name()), cmmt, c.Overwrite)
		}
	}
	return nil
}

// Scans the root directory plus all subdirectories for zip archives and parse any found comments.
func (c *Config) Scans(root string) error {
	hashes := map[[32]byte]bool{}
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
			fmt.Println(err)
			return nil
		}
		if cmmt == "" {
			return nil
		}
		if c.NoDupes {
			hash := sha256.Sum256([]byte(strings.TrimSpace(cmmt)))
			if hashes[hash] {
				return nil
			}
			hashes[hash] = true
		}
		c.cmmts++
		fmt.Print(c.Separator(path))
		if c.Print {
			stdout(cmmt)
		}
		if c.ExportFile {
			fmt.Println(path)
			save(path, cmmt, c.Overwrite)
		}
		return err
	})
	return err
}

// Separator prints and stylises the named file.
func (c Config) Separator(name string) string {
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
	if c.NoDupes {
		unq = "unique "
	}
	s := ""
	if c.Print {
		s = "\n"
	}
	return s + fmt.Sprintf("Scanned %d zip %s and found %d %s%s\n", c.zips, a, c.cmmts, unq, cm)
}

// Save a zip cmmt to the file path.
// Unless the overwrite argument is set, any previous cmmt textfiles are skipped.
func save(path, cmmt string, ow bool) {
	if cmmt == "" {
		return
	}
	name := strings.TrimSuffix(path, filepath.Ext(path)) + filename
	if !ow {
		if s, err := os.Stat(name); err == nil {
			fmt.Printf("export skipped, file already exists: %s (%dB)\n", name, s.Size())
			return
		}
	}
	f, err := os.Create(name)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if i, err := f.Write([]byte(cmmt)); err != nil {
		log.Println(err)
	} else if i == 0 {
		if err1 := os.Remove(name); err != nil {
			log.Println(err1)
		}
	}
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
