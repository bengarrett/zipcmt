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
