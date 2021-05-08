// Package main is a batch viewer and extractor for large collections of zip archives.
// © Ben Garrett https://github.com/bengarrett/zipcmt
package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

type config struct {
	unicode    bool
	raw        bool
	exportfile bool
	exportdir  string
	recursive  bool
	overwrite  bool
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

	flag.Usage = func() {
		var f *flag.Flag
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "    zipcmt [directories] [options]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "    zipcmt --unicode\t\t\t# scan the working directory and show any found comments")
		fmt.Fprintln(os.Stderr, "    zipcmt ~/Downloads --export\t\t# scan the download directory and save any comments")
		fmt.Fprintln(os.Stderr, "    zipcmt ~/Downloads -r -d=~/text\t# recursively scan the download directory and save comments to a directory")
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
	flag.Parse()

	if *ver {
		info()
		return
	}

	flag.Usage()
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
