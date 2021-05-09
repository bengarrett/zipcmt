// Package main is a batch viewer and extractor for large collections of zip archives.
// © Ben Garrett https://github.com/bengarrett/zipcmt
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/zipcmt/zipcmt"
)

var (
	version = "0.0.0"
	commit  = "unset" // nolint: gochecknoglobals
	date    = "unset" // nolint: gochecknoglobals
)

func main() {
	var c zipcmt.Config
	var recursive bool
	flag.BoolVar(&c.Print, "print", false, "print the comments to the terminal")
	flag.BoolVar(&c.Raw, "raw", false, "use the original comment text encoding instead of Unicode")
	flag.BoolVar(&c.ExportFile, "export", false, "save the comment to a textfile stored alongside the archive (use at your own risk)")
	flag.StringVar(&c.ExportDir, "exportdir", "", "save the comment to a textfile stored in this directory")
	flag.BoolVar(&recursive, "recursive", false, "recursively walk through all subdirectories while scanning for zip archives")
	flag.BoolVar(&c.Overwrite, "overwrite", false, "overwrite any previously exported comment textfiles")
	flag.BoolVar(&c.Quiet, "quiet", false, "suppress zipcmt feedback except for errors")
	flag.BoolVar(&c.NoDupes, "nodupes", false, "no duplicate comments, only show unique finds")
	ver := flag.Bool("version", false, "version and information for this program")
	d := flag.String("d", "", "alias for exportdir")
	e := flag.Bool("e", false, "alias for export")
	n := flag.Bool("n", false, "alias for nodupes")
	o := flag.Bool("o", false, "alias for overwrite")
	q := flag.Bool("q", false, "alias for quiet")
	r := flag.Bool("r", false, "alias for recursive")
	u := flag.Bool("p", false, "alias for print")
	v := flag.Bool("v", false, "alias for version")
	flag.Usage = func() {
		help()
	}
	flag.Parse()
	flags(*ver, *v)
	// parse aliases
	if *r {
		recursive = true
	}
	if *u {
		c.Print = true
	}
	if *e {
		c.ExportFile = true
	}
	if *d != "" {
		c.ExportDir = *d
	}
	if *o {
		c.Overwrite = true
	}
	if *q {
		c.Quiet = true
	}
	if *n {
		c.NoDupes = true
	}
	// sanitize the export directory
	if err := c.Clean(); err != nil {
		log.Fatalln(err)
	}
	// recursive directory scan
	if recursive {
		for _, root := range flag.Args() {
			if err := c.Scans(root); err != nil {
				log.Println(err)
			}
			fmt.Print(c.Status())
		}
		return
	}
	// default flat directory scan
	for _, root := range flag.Args() {
		if err := c.Scan(root); err != nil {
			log.Println(err)
		}
	}
	fmt.Print(c.Status())
}

func flags(ver, v bool) {
	// convience for when a help or version flag is passed as an argument
	for _, arg := range flag.Args() {
		switch strings.ToLower(arg) {
		case "-h", "-help", "--help":
			flag.Usage()
			os.Exit(0)
		case "-v", "-version", "--version":
			info()
			os.Exit(0)
		}
	}
	// print help if no arguments are given
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}
	// print version information
	if ver || v {
		info()
		os.Exit(0)
	}
}

// Help, usage and examples.
func help() {
	var f *flag.Flag
	const ps = string(os.PathSeparator)
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "    zipcmt [options] [directories]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprint(os.Stderr, "    zipcmt --print --nodupes .\t\t")
	fmt.Fprintln(os.Stderr, "# scan the working directory and only show unique comments")
	fmt.Fprintf(os.Stderr, "    zipcmt --export ~%sDownloads\t\t", ps)
	fmt.Fprintln(os.Stderr, "# scan the download directory and save all comments")
	fmt.Fprintf(os.Stderr, "    zipcmt -r -d=~%stext ~%sDownloads\t", ps, ps)
	fmt.Fprintln(os.Stderr, "# recursively scan the download directory and save all comments to a directory")
	fmt.Fprintf(os.Stderr, "    zipcmt -n -p -q -r %s | less\t\t", ps)
	fmt.Fprintln(os.Stderr, "# scan the whole system and view unique comments in a page reader")
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

func self() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("self error: %w", err)
	}
	return exe, nil
}
