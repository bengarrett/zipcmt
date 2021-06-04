// © Ben Garrett https://github.com/bengarrett/zipcmt

// Zipcmt is the super-fast, batch, zip file comment viewer, and extractor.
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bengarrett/zipcmt/zipcmt"
	"github.com/gookit/color"
)

var (
	//go:embed embed/logo.txt
	brand string

	version = "0.0.0"
	commit  = "unset" // nolint: gochecknoglobals
	date    = "unset" // nolint: gochecknoglobals
)

const winOS = "windows"

func main() {
	var c zipcmt.Config
	var noprint, norecursive bool
	c.Timer = time.Now()
	flag.BoolVar(&noprint, "noprint", false, "do not print comments to the terminal")
	flag.BoolVar(&c.Raw, "raw", false, "use the original comment text encoding instead of Unicode")
	flag.BoolVar(&c.ExportFile, "export", false, fmt.Sprintf("save the comments to textfiles stored alongside the archive (%s)",
		color.Danger.Sprint("use at your own risk")))
	flag.StringVar(&c.Save, "save", "", "save the comments to uniquely named textfiles in this directory")
	flag.BoolVar(&norecursive, "norecursive", false, "do not recursively walk through any subdirectories while scanning for zip archives")
	flag.BoolVar(&c.Overwrite, "overwrite", false, "overwrite any previously exported comment textfiles")
	flag.BoolVar(&c.Quiet, "quiet", false, "suppress zipcmt feedback except for errors")
	flag.BoolVar(&c.Dupes, "dupes", false, "show duplicate comments from different zips")
	ver := flag.Bool("version", false, "version and information for this program")
	s := flag.String("s", "", "alias for save")
	e := flag.Bool("e", false, "alias for export")
	d := flag.Bool("d", false, "alias for dupes")
	o := flag.Bool("o", false, "alias for overwrite")
	q := flag.Bool("q", false, "alias for quiet")
	r := flag.Bool("r", false, "alias for norecursive")
	u := flag.Bool("p", false, "alias for noprint")
	v := flag.Bool("v", false, "alias for version")
	flag.Usage = func() {
		help(true)
	}
	flag.Parse()
	flags(ver, v)
	// parse aliases
	if *r {
		norecursive = true
	}
	if *u || noprint {
		c.Print = false
	} else {
		c.Print = true
	}
	if *e {
		c.ExportFile = true
	}
	if *s != "" {
		c.Save = *s
	}
	if *o {
		c.Overwrite = true
	}
	if *q {
		c.Quiet = true
	}
	if *d {
		c.Dupes = true
	}
	// sanitize the export directory
	if err := c.Clean(); err != nil {
		color.Error.Tips(fmt.Sprint(err))
	}
	// recursive directory scan
	if !norecursive {
		for _, root := range flag.Args() {
			if err := c.Walk(root); err != nil {
				color.Error.Tips(fmt.Sprint(err))
			}
			fmt.Print(c.Status())
		}
		return
	}
	// default flat directory scan
	for _, root := range flag.Args() {
		if err := c.Scan(root); err != nil {
			color.Error.Tips(fmt.Sprint(err))
		}
	}
	fmt.Print(c.Status())
}

func flags(ver, v *bool) {
	// convience for when a help or version flag is passed as an argument
	for _, arg := range flag.Args() {
		switch strings.ToLower(arg) {
		case "-h", "-help", "--help":
			help(false)
			os.Exit(0)
		case "-v", "-version", "--version":
			info()
			os.Exit(0)
		}
	}
	// print version information
	if *ver || *v {
		info()
		os.Exit(0)
	}
	// print help if no arguments are given
	if len(flag.Args()) == 0 {
		if runtime.GOOS == winOS {
			color.Warn.Println("zipcmt requires at least one directory or drive letter to scan")
		} else {
			color.Warn.Println("zipcmt requires at least one directory to scan")
		}
		fmt.Println()
		help(false)
		os.Exit(0)
	}
}

// Help, usage and examples.
func help(logo bool) {
	var f *flag.Flag
	const ps = string(os.PathSeparator)
	if logo {
		fmt.Fprintln(os.Stderr, brand)
	}
	fmt.Fprintln(os.Stderr, "Usage:")
	if runtime.GOOS == winOS {
		fmt.Fprintln(os.Stderr, "    zipcmt [options] <directories or drive letters>")
	} else {
		fmt.Fprintln(os.Stderr, "    zipcmt [options] <directories>")
	}
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprint(os.Stderr, color.Info.Sprint("    zipcmt .\t\t\t\t"))
	fmt.Fprintln(os.Stderr, color.Note.Sprint("# scan the current directory and subdirectories for unique zipfile comments"))
	if runtime.GOOS == winOS {
		if hd, err := os.UserHomeDir(); err == nil {
			fmt.Fprintln(os.Stderr, color.Info.Sprintf("    zipcmt -save=C:\\text %s%sDownloads\t\t", hd, ps))
			fmt.Fprintln(os.Stderr, color.Note.Sprint("\t\t\t\t\t# recursively scan the download directory and save found comments to the C:\\text directory"))
		}
		fmt.Fprint(os.Stderr, color.Info.Sprint("    zipcmt -s=C:\\text C:\t\t"))
		fmt.Fprintln(os.Stderr, color.Note.Sprint("# recursively scan the C: drive and save any found comments to the C:\\text directory"))
		fmt.Fprint(os.Stderr, color.Info.Sprint("    zipcmt -q C: D: | more\t"))
		fmt.Fprintln(os.Stderr, color.Note.Sprint("# scan the C: plus the D: drive and view unique comments in a page reader"))
	} else {
		fmt.Fprintln(os.Stderr, color.Info.Sprintf("    zipcmt -save=~%stext ~%sDownloads\t\t", ps, ps))
		fmt.Fprintln(os.Stderr, color.Note.Sprint("\t\t\t\t\t# recursively scan the download directory and save found comments to a directory"))
		fmt.Fprint(os.Stderr, color.Info.Sprintf("    zipcmt -r -s=~%stext ~%sDownloads\t", ps, ps))
		fmt.Fprintln(os.Stderr, color.Note.Sprint("# recursively scan the download directory and save all comments to a directory"))
		fmt.Fprint(os.Stderr, color.Info.Sprintf("    zipcmt -r -n -q %s | less\t\t", ps))
		fmt.Fprintln(os.Stderr, color.Note.Sprint("# scan the whole system and view unique comments in a page reader"))
	}
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 4, ' ', 0)
	f = flag.Lookup("norecursive")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", "r", f.Name, f.Usage)
	f = flag.Lookup("dupes")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("noprint")
	fmt.Fprintf(w, "    -p, -%v\t%v\n", f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("save")
	fmt.Fprintf(w, "    -%v, -%v=DIRECTORY\t%v\n", "s", f.Name, f.Usage)
	f = flag.Lookup("overwrite")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("export")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("raw")
	fmt.Fprintf(w, "        -%v\t%v\n", f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("quiet")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("version")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	fmt.Fprintln(w, "    -h, -help\tshow this list of options")
	fmt.Fprintln(w)
	w.Flush()
}

// Info prints out the program information and version.
func info() {
	const copyright = "\u00A9"
	fmt.Println(brand)
	fmt.Printf("zipcmt v%s\n%s 2021 Ben Garrett, logo by sensenstahl\n", version, copyright)
	fmt.Printf("https://github.com/bengarrett/zipcmt\n\n")
	fmt.Printf("build: %s (%s)\n", commit, date)
	exe, err := self()
	if err != nil {
		fmt.Printf("path: %s\n", err)
		return
	}
	fmt.Printf("path: %s\n", exe)
}

func self() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("self error: %w", err)
	}
	return exe, nil
}
