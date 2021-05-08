// Package main is a batch viewer and extractor for large collections of zip archives.
// © Ben Garrett https://github.com/bengarrett/zipcmt
package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	rt "github.com/bengarrett/retrotxtgo/lib/convert"
)

type config struct {
	exportdir  string
	exportfile bool
	nodupes    bool
	overwrite  bool
	raw        bool
	recursive  bool
	unicode    bool
}

var (
	version = "0.0.0"
	commit  = "unset" // nolint: gochecknoglobals
	date    = "unset" // nolint: gochecknoglobals
)

func main() {
	var c config
	flag.BoolVar(&c.unicode, "unicode", false, "convert the zip comment to Unicode and print to the terminal")
	flag.BoolVar(&c.raw, "raw", false, "print the zip comment to the terminal")
	flag.BoolVar(&c.exportfile, "export", false, "save the zip comment to a textfile stored alongside the archive")
	flag.StringVar(&c.exportdir, "exportdir", "", "save the zip comment to a unique textfile stored this directory")
	flag.BoolVar(&c.recursive, "recursive", false, "recursively walk through all subdirectories while scanning for zip archives")
	flag.BoolVar(&c.overwrite, "overwrite", false, "overwrite any previously exported zip comment textfiles")
	ver := flag.Bool("version", false, "version and information for this program")
	r := flag.Bool("r", false, "alias for recursive")
	u := flag.Bool("u", false, "alias for unicode")
	e := flag.Bool("e", false, "alias for export")
	d := flag.String("d", "", "alias for exportdir")
	o := flag.Bool("o", false, "alias for overwrite")
	v := flag.Bool("v", false, "alias for version")

	flag.Usage = func() {
		help()
	}
	flag.Parse()

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
		c.unicode = true
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
	// print help if no arguments are given
	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}
	// if err := c.scan(flag.Arg(0)); err != nil {
	// 	log.Println(err)
	// }
	if c.recursive {
		if err := c.scans(flag.Arg(0)); err != nil {
			log.Println(err)
		}
		return
	}
	if err := c.scan(flag.Arg(0)); err != nil {
		log.Println(err)
	}
}

func help() {
	var f *flag.Flag
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "    zipcmt [options] [directories]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "    zipcmt --unicode .\t\t\t# scan the working directory and show any found comments")
	fmt.Fprintln(os.Stderr, "    zipcmt --export ~/Downloads\t\t# scan the download directory and save any comments")
	fmt.Fprintln(os.Stderr, "    zipcmt -r -d=~/text ~/Downloads \t# recursively scan the download directory and save comments to a directory")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 4, ' ', 0)
	f = flag.Lookup("recursive")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("raw")
	fmt.Fprintf(w, "        --%v\t%v\n", f.Name, f.Usage)
	f = flag.Lookup("unicode")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("export")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("exportdir")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", "d", f.Name, f.Usage)
	f = flag.Lookup("overwrite")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	fmt.Fprintln(w, "    -h, --help\tshow this list of options")
	f = flag.Lookup("version")
	fmt.Fprintf(w, "    -%v, --%v\t%v\n", f.Name[:1], f.Name, f.Usage)
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

func (c config) scan(root string) error {
	files, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, file := range files {
		name := filepath.Join(root, file.Name())
		if !valid(name) {
			continue
		}
		cmmt, err := read(name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(cmmt) == 0 {
			continue
		}
		separator(name)
		if c.raw {
			fmt.Println(cmmt)
		}
		if c.unicode {
			if err := unicode(cmmt); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func (c config) scans(root string) error {
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
		cmmt, err := read(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if len(cmmt) == 0 {
			return nil
		}
		separator(path)
		if c.raw {
			fmt.Println(cmmt)
		}
		if c.unicode {
			if err := unicode(cmmt); err != nil {
				fmt.Println(err)
				return nil
			}
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

func read(name string) (cmmt string, err error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		return "", err
	}
	defer r.Close()
	if cmmt := r.Comment; cmmt != "" {
		if strings.TrimSpace(cmmt) == "" {
			return "", nil
		}
		return cmmt, nil
	}
	return "", nil
}

func unicode(cmmt string) error {
	const resetCmd = "\033[0m"
	b, err := rt.D437(cmmt)
	if err != nil {
		return err
	}
	fmt.Printf("%s%s\n", b, resetCmd)
	return nil
}

func separator(name string) {
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
