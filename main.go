// Package main is a batch viewer and extractor for large collections of zip archives.
// © Ben Garrett https://github.com/bengarrett/zipcmt
package main

import (
	"archive/zip"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/convert"
)

type config struct {
	exportdir  string
	exportfile bool
	nodupes    bool
	overwrite  bool
	raw        bool
	recursive  bool
	print      bool
	quiet      bool
	zips       int
	cmmts      int
}

var (
	version = "0.0.0"
	commit  = "unset" // nolint: gochecknoglobals
	date    = "unset" // nolint: gochecknoglobals
)

func main() {
	var c config
	flag.BoolVar(&c.print, "print", false, "print the comments to the terminal")
	flag.BoolVar(&c.raw, "raw", false, "use the original comment text encoding instead of Unicode")
	flag.BoolVar(&c.exportfile, "export", false, "save the comment to a textfile stored alongside the archive (use at your own risk)")
	flag.StringVar(&c.exportdir, "exportdir", "", "save the comment to a textfile stored in this directory")
	flag.BoolVar(&c.recursive, "recursive", false, "recursively walk through all subdirectories while scanning for zip archives")
	flag.BoolVar(&c.overwrite, "overwrite", false, "overwrite any previously exported comment textfiles")
	flag.BoolVar(&c.quiet, "quiet", false, "suppress zipcmt feedback except for errors")
	flag.BoolVar(&c.nodupes, "nodupes", false, "no duplicate comments, only show unique finds")
	ver := flag.Bool("version", false, "version and information for this program")
	r := flag.Bool("r", false, "alias for recursive")
	u := flag.Bool("p", false, "alias for print")
	e := flag.Bool("e", false, "alias for export")
	d := flag.String("d", "", "alias for exportdir")
	o := flag.Bool("o", false, "alias for overwrite")
	v := flag.Bool("v", false, "alias for version")
	q := flag.Bool("q", false, "alias for quiet")
	n := flag.Bool("n", false, "alias for nodupes")

	flag.Usage = func() {
		help()
	}
	flag.Parse()

	// help convience for when a help flag is passed as an argument
	for _, arg := range flag.Args() {
		arg = strings.ToLower(arg)
		if arg == "-h" || arg == "-help" || arg == "--help" {
			flag.Usage()
			return
		}
		if arg == "-v" || arg == "-version" || arg == "--version" {
			info()
			return
		}
	}
	// print help if no arguments are given
	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	// version information
	if *ver || *v {
		info()
		return
	}
	// parse aliases
	if *r {
		c.recursive = true
	}
	if *u {
		c.print = true
	}
	if *e {
		c.exportfile = true
	}
	if *d != "" {
		c.exportdir = *d
	}
	if *o {
		c.overwrite = true
	}
	if *q {
		c.quiet = true
	}
	if *n {
		c.nodupes = true
	}
	// export to directory sanity check
	if c.exportdir != "" {
		c.exportdir = filepath.Clean(c.exportdir)
		p := strings.Split(c.exportdir, string(filepath.Separator))
		if p[0] == "~" {
			hd, err := os.UserHomeDir()
			if err != nil {
				log.Fatalln(err)
			}
			c.exportdir = strings.Replace(c.exportdir, "~", hd, 1)
		}
		s, err := os.Stat(c.exportdir)
		if err != nil {
			log.Fatalln(err)
		}
		if !s.IsDir() {
			log.Fatalln(os.ErrInvalid)
		}
	}
	// recursive directory scan
	if c.recursive {
		for _, root := range flag.Args() {
			if err := c.scans(root); err != nil {
				log.Println(err)
			}
			c.status()
		}
		return
	}
	// default flat directory scan
	for _, root := range flag.Args() {
		if err := c.scan(root); err != nil {
			log.Println(err)
		}
	}
	c.status()
}

func help() {
	var f *flag.Flag
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "    zipcmt [options] [directories]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "    zipcmt --print --nodupes .\t\t# scan the working directory and only show unique comments")
	fmt.Fprintln(os.Stderr, "    zipcmt --export ~/Downloads\t\t# scan the download directory and save all comments")
	fmt.Fprintln(os.Stderr, "    zipcmt -r -d=~/text ~/Downloads\t# recursively scan the download directory and save all comments to a directory")
	fmt.Fprintln(os.Stderr, "    zipcmt -n -p -q -r / | less\t\t# scan the whole system and view unique comments in a page reader")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 4, ' ', 0)
	f = flag.Lookup("recursive")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("nodupes")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("print")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("export")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("exportdir")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", "d", f.Name, f.Usage)
	f = flag.Lookup("overwrite")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("raw")
	fmt.Fprintf(w, "        --%v\t%v\n", f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("quiet")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("version")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	fmt.Fprintln(w, "    -h, --help\tshow this list of options")
	fmt.Fprintln(w)
	w.Flush()
}

func self() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("self error: %w", err)
	}
	return exe, nil
}

// Info prints out the program information and version.
func info() {
	const copyright = "\u00A9"
	fmt.Printf("zipcmt v%s\n%s 2021 Ben Garrett\n", version, copyright)
	fmt.Printf("https://github.com/bengarrett/zipcmt\n\n")
	fmt.Printf("build: %s (%s)\n", commit, date)
	exe, err := self()
	if err != nil {
		fmt.Printf("path: %s\n", err)
		return
	}
	fmt.Printf("path:  %s\n", exe)
}

func (c *config) scan(root string) error {
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
		c.zips = c.zips + 1
		cmmt, err := c.read(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(cmmt) == 0 {
			continue
		}
		if c.nodupes {
			hash := sha256.Sum256([]byte(cmmt))
			if hashes[hash] {
				continue
			}
			hashes[hash] = true
		}
		c.cmmts = c.cmmts + 1
		c.separator(path)
		if c.print {
			if err := c.stdout(cmmt); err != nil {
				fmt.Println(err)
			}
		}
		if c.exportfile {
			go save(path, cmmt, c.overwrite)
		}
		if c.exportdir != "" {
			go save(filepath.Join(c.exportdir, file.Name()), cmmt, c.overwrite)
		}
	}
	return nil
}

func (c *config) scans(root string) error {
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
		c.zips = c.zips + 1
		cmmt, err := c.read(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if len(cmmt) == 0 {
			return nil
		}
		if c.nodupes {
			hash := sha256.Sum256([]byte(cmmt))
			if hashes[hash] {
				return nil
			}
			hashes[hash] = true
		}
		c.cmmts = c.cmmts + 1
		c.separator(path)
		if c.print {
			if err := c.stdout(cmmt); err != nil {
				fmt.Println(err)
				return nil
			}
		}
		if c.exportfile {
			fmt.Println(path)
			save(path, cmmt, c.overwrite)
		}
		return err
	})
	return err
}

func valid(name string) bool {
	const zip = ".zip"
	switch filepath.Ext(strings.ToLower(name)) {
	case zip:
		return true
	}
	return false
}

func (c config) read(name string) (cmmt string, err error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		return "", nil
	}
	defer r.Close()
	if cmmt := r.Comment; cmmt != "" {
		if strings.TrimSpace(cmmt) == "" {
			return "", nil
		}
		if !c.raw {
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

func save(path, cmmt string, ow bool) {
	if cmmt == "" {
		return
	}
	name := strings.TrimSuffix(path, filepath.Ext(path)) + "-zipcomment.txt"
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

func (c config) stdout(cmmt string) error {
	const resetCmd = "\033[0m"
	fmt.Printf("%s%s\n", cmmt, resetCmd)
	return nil
}

func (c config) separator(name string) {
	if !c.print || c.quiet {
		return
	}
	const file_id = 45
	const pointer = " \u2500\u2500 "
	if h, err := os.UserHomeDir(); err == nil {
		if len(name) > len(h) && name[0:len(h)] == h {
			name = strings.Replace(name, h, "~", 1)
		}
	}
	l := len(pointer) + len(name)
	if l >= file_id {
		fmt.Printf("%s%s\n", pointer, name)
		return
	}
	fmt.Printf("\n%s%s %s\u2510\n", pointer, name, strings.Repeat("\u2500", file_id-l))
}

func (c config) status() {
	if c.quiet {
		return
	}
	a, cm, unq := "archive", "comment", ""
	if c.zips != 1 {
		a += "s"
	}
	if c.cmmts != 1 {
		cm += "s"
	}
	if c.nodupes {
		unq = "unique "
	}
	fmt.Printf("Scanned %d zip %s and found %d %s%s\n", c.zips, a, c.cmmts, unq, cm)
}
