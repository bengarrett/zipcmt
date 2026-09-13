//nolint:exhaustruct_v5,testableexamples
package main_test

// © Ben Garrett https://github.com/bengarrett/zipcmt

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bengarrett/zipcmt/app"
	"github.com/gookit/color"
)

func init() {
	color.Enable = false
}

func ExampleConfig() {
	// print all comments found in the test directory
	example := []string{"testdata"}
	a := app.Config{
		Dirs:  example,
		Dupes: true,
		Print: true,
	}
	a.WalkDirs()
	if s := a.Status(); s != "" {
		fmt.Fprintln(os.Stdout, s)
	}

	// quietly scan and save only the unique comments as text files in the home directory
	const homeDir = "~"
	b := app.Config{
		Dirs:     example,
		SaveName: homeDir,
		Quiet:    true,
	}
	b.WalkDirs()
	if s := b.Status(); s != "" {
		fmt.Fprintln(os.Stdout, s)
	}

	// quietly scan and count the unique comments
	c := app.Config{
		Dirs:  example,
		Quiet: true,
	}
	c.WalkDirs()
	fmt.Fprintf(os.Stdout, "Scanned %d zip archives and found %d unique comments\n", c.Zips, c.Cmmts)
}

func ExampleRead() {
	s, err := app.Read(filepath.Join("testdata", "test-with-comment.zip"), false)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(os.Stdout, s)
	// Output:
	// This is an example test comment for zipcmmt.
	//
}

func ExampleConfig_Status() {
	c := app.Config{}
	c.SetTest()
	if err := c.WalkDir("testdata"); err != nil {
		log.Panicln(err)
	}
	fmt.Fprint(os.Stdout, c.Status())

	c = app.Config{
		Dupes: true,
	}
	c.SetTest()
	if err := c.WalkDir("testdata"); err != nil {
		log.Panicln(err)
	}
	fmt.Fprint(os.Stdout, c.Status())
	// Output:
	// Scanned 4 zip archives and found 1 unique comment
	// Scanned 4 zip archives and found 2 comments
}
