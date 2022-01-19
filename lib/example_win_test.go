//go:build windows
// +build windows

// © Ben Garrett https://github.com/bengarrett/zipcmt

package zipcmt

import (
	"fmt"
	"log"
)

func ExampleConfig_clean() {
	c := Config{
		SaveName: "..//test///.",
	}
	if err := c.clean(); err != nil {
		log.Fatalln(err)
	}
	fmt.Print(c.SaveName)
	// Output: ..\test
}

func ExampleConfig_WalkDir() {
	c := Config{
		Print: true,
		Dupes: true,
	}
	if err := c.WalkDir("../test"); err != nil {
		log.Panicln(err)
	}
	// Output:
	// ── ..\test\subdir\test-with-comment.zip ─┐
	//    This is an example test comment for zipcmmt.[0m
	//
	//  ── ..\test\test-with-comment.zip ────────┐
	//    This is an example test comment for zipcmmt.[0m
}
