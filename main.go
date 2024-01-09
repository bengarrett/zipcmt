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

	"github.com/bengarrett/zipcmt/internal/misc"
	zipcmt "github.com/bengarrett/zipcmt/pkg"
	"github.com/gookit/color"
)

var (
	//go:embed embed/logo.txt
	brand string

	version = "0.0.0"
	commit  = "unset"
	date    = "unset"
)

const winOS = "windows"

func main() {
	const ellipsis = "\u2026"
	var c zipcmt.Config
	var noprint bool
	c.SetTimer()
	flag.BoolVar(&noprint, "noprint", false,
		"do not print comments to the terminal to improve the performance of the scan")
	flag.BoolVar(&c.NoWalk, "norecursive", false,
		"do not recursively walk through any subdirectories while scanning for zip archives")
	flag.BoolVar(&c.Export, "export", false,
		fmt.Sprintf("save the comments as text files stored alongside the zip files (%s)",
			color.Danger.Sprint("use at your own risk")))
	flag.BoolVar(&c.Dupes, "all", false,
		"show all comments, including duplicates in multiple zips")
	flag.BoolVar(&c.Now, "now", false,
		"do not use the last modification date sourced from the zip files")
	flag.BoolVar(&c.Log, "log", false,
		"create a logfile for debugging")
	flag.BoolVar(&c.Overwrite, "overwrite", false,
		"overwrite any previously exported comment text files")
	flag.BoolVar(&c.Quiet, "quiet", false,
		"suppress zipcmt feedback except for errors")
	flag.BoolVar(&c.Raw, "raw", false,
		"use the original comment text encoding (CP437, ISO-8859"+ellipsis+") instead of Unicode")
	flag.StringVar(&c.SaveName, "save", "",
		"save the comments to uniquely named textfiles in this directory")
	ver := flag.Bool("version", false,
		"version and information for this program")
	aliasA := flag.Bool("a", false, "alias for all")
	aliasO := flag.Bool("o", false, "alias for overwrite")
	aliasQ := flag.Bool("q", false, "alias for quiet")
	aliasR := flag.Bool("r", false, "alias for norecursive")
	aliasS := flag.String("s", "", "alias for save")
	aliasU := flag.Bool("p", false, "alias for noprint")
	aliasV := flag.Bool("v", false, "alias for version")
	flag.Usage = func() {
		help(true)
	}
	flag.Parse()
	flags(ver, aliasV, aliasQ)
	// parse aliases
	if *aliasR {
		c.NoWalk = true
	}
	if *aliasU || noprint {
		c.Print = false
	} else {
		c.Print = true
	}
	if *aliasS != "" {
		c.SaveName = *aliasS
	}
	if *aliasO {
		c.Overwrite = true
	}
	if *aliasQ {
		c.Quiet = true
	}
	if *aliasA {
		c.Dupes = true
	}
	// directories to scan
	c.Dirs = flag.Args()
	// file and directory scan
	c.WalkDirs()
	// summaries
	fmt.Fprintln(os.Stdout, c.Status())
	if s := c.LogName(); s != "" {
		fmt.Fprintf(os.Stdout, "%s %s\n", "The log is found at", color.Primary.Sprint(s))
	}
}

func flags(ver, aliasV, quiet *bool) {
	// convenience for when a help or version flag is passed as an argument
	for _, arg := range flag.Args() {
		showLogo := !*quiet
		switch strings.ToLower(arg) {
		case "-h", "-help", "--help":
			help(showLogo)
			os.Exit(0)
		case "-v", "-version", "--version":
			info(quiet)
			os.Exit(0)
		}
	}
	// print version information
	if *ver || *aliasV {
		info(quiet)
		os.Exit(0)
	}
	// print help if no arguments are given
	if len(flag.Args()) == 0 {
		if runtime.GOOS == winOS {
			fmt.Fprintln(os.Stderr, color.Warn.Sprint("zipcmt requires at least one directory or drive letter to scan"))
		} else {
			fmt.Fprintln(os.Stderr, color.Warn.Sprint("zipcmt requires at least one directory to scan"))
		}
		fmt.Fprintln(os.Stderr)
		help(false)
		os.Exit(0)
	}
}

func helpPosix() {
	const ps = string(os.PathSeparator)
	fmt.Fprintln(os.Stderr, "    zipcmt [options] <directories>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprint(os.Stderr, color.Info.Sprint("    zipcmt .\t\t\t\t"))
	fmt.Fprintln(os.Stderr,
		color.Note.Sprint("# scan the current directory and subdirectories for unique comments"))
	fmt.Fprint(os.Stderr, color.Info.Sprintf("    zipcmt -save=~%stext ~%sDownloads\t", ps, ps))
	fmt.Fprintln(os.Stderr,
		color.Note.Sprint("# scan the user download directories and save unique comments to a directory"))
	fmt.Fprint(os.Stderr, color.Info.Sprintf("    zipcmt -a -s=~%stext ~%sDownloads\t", ps, ps))
	fmt.Fprintln(os.Stderr,
		color.Note.Sprint("# scan the user download directories and save all comments to a directory"))
	fmt.Fprint(os.Stderr, color.Info.Sprintf("    zipcmt -quiet %s | less\t\t", ps))
	fmt.Fprintln(os.Stderr,
		color.Note.Sprint("# scan the whole system to view the unique comments in a page reader"))
}

func helpWin() {
	const ps = string(os.PathSeparator)
	fmt.Fprintln(os.Stderr, "    zipcmt [options] <directories or drive letters>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprint(os.Stderr, color.Info.Sprint("    zipcmt .\t\t\t"))
	fmt.Fprintln(os.Stderr,
		color.Note.Sprint("# scan the current directory and subdirectories for unique comments"))
	if hd, err := os.UserHomeDir(); err == nil {
		fmt.Fprintln(os.Stderr, color.Info.Sprintf("    zipcmt -save=C:\\text %s%sDownloads\t\t", hd, ps))
		fmt.Fprintln(os.Stderr, color.Note.Sprint("\t\t\t\t# scan the files and directories in Downloads"+
			" and save the unique comments to 'C:\\text'"))
	}
	fmt.Fprint(os.Stderr, color.Info.Sprint("    zipcmt -save=C:\\text C:\t"))
	fmt.Fprintln(os.Stderr,
		color.Note.Sprint("# scan the 'C' drive and save the unique comments to the 'C:\\text' directory"))
	fmt.Fprint(os.Stderr, color.Info.Sprint("    zipcmt -quiet C: D: | more\t"))
	fmt.Fprintln(os.Stderr,
		color.Note.Sprint("# scan the 'C' and 'D' drives to view the unique comments in a page reader"))
}

// Help, usage and examples.
func help(logo bool) {
	var f *flag.Flag
	if logo {
		fmt.Fprintln(os.Stderr, brand)
		fmt.Fprint(os.Stderr, " Zip Comment is the super-fast batch, zip file comment viewer, and extractor.\n"+
			" Using a modern PC with the zip files stored on a solid-state drive,\n"+
			" zipcmt handles many thousands of archives per second.\n\n")
	}
	fmt.Fprintln(os.Stderr, "Usage:")
	if runtime.GOOS == winOS {
		helpWin()
	} else {
		helpPosix()
	}
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	const padding = 4
	w := tabwriter.NewWriter(os.Stderr, 0, 0, padding, ' ', 0)
	f = flag.Lookup("save")
	fmt.Fprintf(w, "    -%v, -%v=DIRECTORY\t%v\n", "s", f.Name, f.Usage)
	f = flag.Lookup("overwrite")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("noprint")
	fmt.Fprintf(w, "    -p, -%v\t%v\n", f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("norecursive")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", "r", f.Name, f.Usage)
	f = flag.Lookup("all")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("now")
	fmt.Fprintf(w, "        -%v\t%v\n", f.Name, f.Usage)
	f = flag.Lookup("raw")
	fmt.Fprintf(w, "        -%v\t%v\n", f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("export")
	fmt.Fprintf(w, "        -%v\t%v\n", f.Name, f.Usage)
	fmt.Fprintln(w, "                \t")
	f = flag.Lookup("quiet")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	f = flag.Lookup("version")
	fmt.Fprintf(w, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	fmt.Fprintln(w, "    -h, -help\tshow this list of options")
	fmt.Fprintln(w)
	optimial(w)
	w.Flush()
}

func optimial(w *tabwriter.Writer) {
	if runtime.GOOS == winOS {
		fmt.Fprintln(w, "For optimal performance Windows users may wish to temporarily disable"+
			" the Virus & threat 'Real-time protection' under Windows Security.")
		fmt.Fprintln(w, "Or create Windows Security Exclusions for the directories to be scanned.")
		fmt.Fprintln(w, "https://support.microsoft.com/en-us/windows/"+
			"add-an-exclusion-to-windows-security-811816c0-4dfd-af4a-47e4-c301afe13b26")
	}
}

// Info prints out the program information and version.
func info(quiet *bool) {
	const copyright = "\u00A9"
	if !*quiet {
		fmt.Fprintln(os.Stdout, brand)
	}
	fmt.Fprintf(os.Stdout, "zipcmt v%s\n%s 2021-24 Ben Garrett, logo by sensenstahl\n", version, copyright)
	fmt.Fprintf(os.Stdout, "https://github.com/bengarrett/zipcmt\n\n")
	fmt.Fprintf(os.Stdout, "build: %s (%s)\n", commit, date)
	exe, err := misc.Self()
	if err != nil {
		fmt.Fprintf(os.Stdout, "path: %s\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, "path: %s\n", exe)
}
