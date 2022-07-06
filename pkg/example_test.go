// © Ben Garrett https://github.com/bengarrett/zipcmt
package zipcmt_test

import (
	"fmt"
	"log"

	zipcmt "github.com/bengarrett/zipcmt/pkg"
	"github.com/gookit/color"
)

func init() { // nolint:gochecknoinits
	color.Enable = false
}

func ExampleConfig() {
	// print all comments found in the test directory
	example := []string{"../test"}
	a := zipcmt.Config{
		Dirs:  example,
		Dupes: true,
		Print: true,
	}
	a.WalkDirs()
	if s := a.Status(); s != "" {
		fmt.Println(s)
	}

	// quietly scan and save only the unique comments as text files in the home directory
	const homeDir = "~"
	b := zipcmt.Config{
		Dirs:     example,
		SaveName: homeDir,
		Quiet:    true,
	}
	if err := b.Clean(); err != nil {
		log.Fatalln(err)
	}
	b.WalkDirs()
	if s := b.Status(); s != "" {
		fmt.Println(s)
	}

	// quietly scan and count the unique comments
	c := zipcmt.Config{
		Dirs:  example,
		Quiet: true,
	}
	c.WalkDirs()
	fmt.Printf("Scanned %d zip archives and found %d unique comments\n", c.Zips, c.Cmmts)
}

func ExampleRead() {
	s, err := zipcmt.Read("../test/test-with-comment.zip", false)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(s)
	// Output:
	// This is an example test comment for zipcmmt.
	//
}

func ExampleConfig_Status() {
	c := zipcmt.Config{}
	c.SetTest()
	if err := c.WalkDir("../test"); err != nil {
		log.Panicln(err)
	}
	fmt.Print(c.Status())

	c = zipcmt.Config{
		Dupes: true,
	}
	c.SetTest()
	if err := c.WalkDir("../test"); err != nil {
		log.Panicln(err)
	}
	fmt.Print(c.Status())
	// Output:
	// Scanned 4 zip archives and found 1 unique comment
	// Scanned 4 zip archives and found 2 comments
}
