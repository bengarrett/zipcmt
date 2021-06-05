// +build windows
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
	// Output: ..\test
}

func ExampleConfig_Walk() {
	c := Config{
		Print: true,
		Dupes: true,
	}
	if err := c.Walk("../test"); err != nil {
		log.Panicln(err)
	}
	// Output:
	// ── ..\test\subdir\test-with-comment.zip ─┐
	//    This is an example test comment for zipcmmt.[0m
	//
	//  ── ..\test\test-with-comment.zip ────────┐
	//    This is an example test comment for zipcmmt.[0m
}
