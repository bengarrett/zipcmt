// © Ben Garrett https://github.com/bengarrett/zipcmt

package zipcmt

import (
	"fmt"
	"log"

	"github.com/gookit/color"
)

func init() {
	color.Enable = false
}

func ExampleConfig_Read() {
	c := Config{
		Raw: false,
	}
	s, err := c.Read("../test/test-with-comment.zip")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(s)
	// Output:
	//This is an example test comment for zipcmmt.
	//
}

func ExampleConfig_Status() {
	c := Config{
		test: true,
	}
	if err := c.Walk("../test"); err != nil {
		log.Panicln(err)
	}
	fmt.Print(c.Status())

	c = Config{
		Dupes: true,
		test:  true,
	}
	if err := c.Walk("../test"); err != nil {
		log.Panicln(err)
	}
	fmt.Print(c.Status())
	// Output:
	// Scanned 4 zip archives and found 1 unique comment
	// Scanned 4 zip archives and found 2 comments
}
