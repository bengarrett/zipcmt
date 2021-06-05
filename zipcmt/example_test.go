// © Ben Garrett https://github.com/bengarrett/zipcmt

package zipcmt

import (
	"fmt"
	"log"
)

func ExampleConfig_Clean() {
	c := Config{
		Save: "..//test///.",
	}
	if err := c.Clean(); err != nil {
		log.Fatalln(err)
	}
	fmt.Print(c.Save)
	// Output: ../test
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

func ExampleConfig_Scan() {
	c := Config{
		Print: true,
		Quiet: true,
	}
	if err := c.Scan("../test"); err != nil {
		log.Println(err)
	}
	// Output:
	//This is an example test comment for zipcmmt.[0m
	//
}

func ExampleConfig_Walk() {
	c := Config{
		Print: true,
	}
	if err := c.Walk("../test"); err != nil {
		log.Panicln(err)
	}
	// Output:
	// ── ../test/subdir/test-with-comment.zip ─┐
	//    This is an example test comment for zipcmmt.[0m
	//
	//  ── ../test/test-with-comment.zip ────────┐
	//    This is an example test comment for zipcmmt.[0m
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
