// © Ben Garrett https://github.com/bengarrett/zipcmt

// Zipcmt is the super-fast, batch, zip file comment viewer, and extractor.
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/zipcmt/internal/cmnt"
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
	var configs zipcmt.Config
	var noprint bool
	configs.SetTimer()
	flag.BoolVar(&noprint, "noprint", false,
		"do not print comments to the terminal to improve the performance of the scan")
	flag.BoolVar(&configs.NoWalk, "norecursive", false,
		"do not recursively walk through any subdirectories while scanning for zip archives")
	flag.BoolVar(&configs.Export, "export", false,
		fmt.Sprintf("save comments to the directories that contain the zip files (%s)",
			color.Danger.Sprint("not advised")))
	flag.BoolVar(&configs.Dupes, "all", false,
		"show every comment, including all the duplicates")
	flag.BoolVar(&configs.Now, "now", false,
		"do not use the last modification date sourced from the zip files")
	flag.BoolVar(&configs.Log, "log", false,
		"create a logfile for debugging")
	flag.BoolVar(&configs.Overwrite, "overwrite", false,
		"overwrite any previously exported comment text files")
	flag.BoolVar(&configs.Quiet, "quiet", false,
		"suppress zipcmt feedback except for errors")
	flag.BoolVar(&configs.Raw, "raw", false,
		"use the original comment text encoding (CP437, ISO-8859"+ellipsis+") instead of Unicode")
	flag.StringVar(&configs.SaveName, "save", "",
		"save the comments to this directory as unique named text files")
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
		help(os.Stderr, true)
	}
	flag.Parse()
	flags(ver, aliasV, aliasQ)
	// parse aliases
	if *aliasR {
		configs.NoWalk = true
	}
	if *aliasU || noprint {
		configs.Print = false
	} else {
		configs.Print = true
	}
	if *aliasS != "" {
		configs.SaveName = *aliasS
	}
	if *aliasO {
		configs.Overwrite = true
	}
	if *aliasQ {
		configs.Quiet = true
	}
	if *aliasA {
		configs.Dupes = true
	}
	// directories to scan
	configs.Dirs = flag.Args()
	// file and directory scan
	configs.WalkDirs()
	// summaries
	fmt.Fprintln(os.Stdout, configs.Status())
	if s := configs.LogName(); s != "" {
		fmt.Fprintf(os.Stdout, "%s %s\n", "The log is found at", color.Primary.Sprint(s))
	}
}

func flags(ver, aliasV, quiet *bool) {
	// convenience for when a help or version flag is passed as an argument
	for _, arg := range flag.Args() {
		showLogo := !*quiet
		switch strings.ToLower(arg) {
		case "-h", "-help", "--help":
			help(os.Stderr, showLogo)
			os.Exit(0)
		case "-v", "-version", "--version":
			info(os.Stdout, quiet)
			os.Exit(0)
		}
	}
	// print version information
	if *ver || *aliasV {
		info(os.Stdout, quiet)
		os.Exit(0)
	}
	// print help if no arguments are given
	w := os.Stderr
	if len(flag.Args()) == 0 {
		s := "zipcmt requires at least one directory to scan"
		if runtime.GOOS == winOS {
			s = "zipcmt requires at least one directory or drive letter to scan"
		}
		fmt.Fprintln(w, color.Warn.Sprint(s)+"\n")
		help(w, false)
		os.Exit(0)
	}
}

func helpPosix(w io.Writer) {
	const ps = string(os.PathSeparator)
	fmt.Fprintln(w, "    zipcmt [options] <directories>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Examples:")
	fmt.Fprint(w, color.Info.Sprint("    zipcmt .\t\t\t\t"))
	fmt.Fprintln(w,
		color.Note.Sprint("# scan the current directory and subdirectories for unique comments"))
	fmt.Fprint(w, color.Info.Sprintf("    zipcmt -save=~%swork ~%sDownloads\t", ps, ps))
	fmt.Fprintln(w,
		color.Note.Sprint("# scan the user downloads directory, then save unique comments to a directory"))
	fmt.Fprint(w, color.Info.Sprintf("    zipcmt -a -s=~%swork ~%sDownloads\t", ps, ps))
	fmt.Fprintln(w,
		color.Note.Sprint("# scan the user downloads directory, then save all comments to a directory"))
	fmt.Fprint(w, color.Info.Sprintf("    zipcmt -quiet %s | less\t\t", ps))
	fmt.Fprintln(w,
		color.Note.Sprint("# scan the whole system to view the unique comments in a page reader"))
}

func helpWin(w io.Writer) {
	const ps = string(os.PathSeparator)
	fmt.Fprintln(w, "    zipcmt [options] <directories or drive letters>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Examples:")
	fmt.Fprint(w, color.Info.Sprint("    zipcmt .\t\t\t"))
	fmt.Fprintln(w,
		color.Note.Sprint("# scan the current directory and subdirectories for unique comments"))
	if hd, err := os.UserHomeDir(); err == nil {
		fmt.Fprintln(w, color.Info.Sprintf("    zipcmt -save=C:\\work %s%sDownloads\t\t", hd, ps))
		fmt.Fprintln(w, color.Note.Sprint("\t\t\t\t# scan the files and directories in Downloads"+
			" and save the unique comments to 'C:\\work'"))
	}
	fmt.Fprint(w, color.Info.Sprint("    zipcmt -save=C:\\work C:\t"))
	fmt.Fprintln(w,
		color.Note.Sprint("# scan the 'C' drive and save the unique comments to the 'C:\\work' directory"))
	fmt.Fprint(w, color.Info.Sprint("    zipcmt -quiet C: D: | more\t"))
	fmt.Fprintln(w,
		color.Note.Sprint("# scan the 'C' and 'D' drives to view the unique comments in a page reader"))
}

// Help, usage and examples.
func help(w io.Writer, logo bool) {
	var f *flag.Flag
	if logo {
		fmt.Fprintln(w, brand)
		fmt.Fprint(w, " Zip Comment is the super-fast, batch zip file-comment viewer and extractor.\n"+
			" Using a modern PC, zipcmt handles many thousands of archives per second.\n\n")
	}
	fmt.Fprintln(w, "Usage:")
	if runtime.GOOS == winOS {
		helpWin(w)
	} else {
		helpPosix(w)
	}
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	const padding = 4
	tw := tabwriter.NewWriter(w, 0, 0, padding, ' ', 0)
	names := []string{
		"save", "overwrite", "noprint", "norecursive", "all", "now", "raw", "export", "quiet", "version",
	}
	for name := range slices.Values(names) {
		f = flag.Lookup(name)
		if f == nil {
			fmt.Fprintf(os.Stderr, "flag lookup failure %q: %s", name, flag.ErrHelp)
			continue
		}
		helper(tw, f, name)
	}
	optimial(tw)
	tw.Flush()
}

func helper(tw *tabwriter.Writer, f *flag.Flag, name string) {
	if tw == nil || f == nil {
		return
	}
	switch name {
	case "save":
		fmt.Fprintf(tw, "    -%v, -%v=DIRECTORY\t%v\n", "s", f.Name, f.Usage)
	case "overwrite":
		fmt.Fprintf(tw, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	case "noprint":
		fmt.Fprintf(tw, "    -p, -%v\t%v\n", f.Name, f.Usage)
		fmt.Fprintln(tw, "                \t")
	case "norecursive":
		fmt.Fprintf(tw, "    -%v, -%v\t%v\n", "r", f.Name, f.Usage)
	case "all":
		fmt.Fprintf(tw, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	case "now":
		fmt.Fprintf(tw, "        -%v\t%v\n", f.Name, f.Usage)
	case "raw":
		fmt.Fprintf(tw, "        -%v\t%v\n", f.Name, f.Usage)
		fmt.Fprintln(tw, "                \t")
	case "export":
		fmt.Fprintf(tw, "        -%v\t%v\n", f.Name, f.Usage)
		fmt.Fprintln(tw, "                \t")
	case "quiet":
		fmt.Fprintf(tw, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
	case "version":
		fmt.Fprintf(tw, "    -%v, -%v\t%v\n", f.Name[:1], f.Name, f.Usage)
		fmt.Fprintln(tw, "    -h, -help\tshow this list of options")
		fmt.Fprintln(tw)
	}
}

func optimial(tw *tabwriter.Writer) {
	if runtime.GOOS != winOS || tw == nil {
		return
	}
	fmt.Fprintln(tw, "For optimal performance Windows users may wish to temporarily disable"+
		" the Virus & threat 'Real-time protection' under Windows Security.")
	fmt.Fprintln(tw, "Or create Windows Security Exclusions for the directories to be scanned.")
	fmt.Fprintln(tw, "https://support.microsoft.com/en-us/windows/"+
		"add-an-exclusion-to-windows-security-811816c0-4dfd-af4a-47e4-c301afe13b26")
}

// Info prints out the program information and version.
func info(w io.Writer, quiet *bool) {
	const copyright = "\u00A9"
	if !*quiet {
		fmt.Fprintln(w, brand)
	}
	fmt.Fprintf(w, "zipcmt v%s\n%s 2021-25 Ben Garrett, logo by sensenstahl\n",
		version, copyright)
	fmt.Fprintf(w, "https://github.com/bengarrett/zipcmt\n\n")
	fmt.Fprintf(w, "build: %s (%s)\n", commit, date)
	exe, err := cmnt.Self()
	if err != nil {
		fmt.Fprintf(w, "path: %s\n", err)
		return
	}
	fmt.Fprintf(w, "path: %s\n", exe)
}
