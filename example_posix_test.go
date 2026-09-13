//go:build !windows

//nolint:exhaustruct_v5
package main_test

// © Ben Garrett https://github.com/bengarrett/zipcmt

import (
	"fmt"
	"log"
	"os"

	"github.com/bengarrett/zipcmt/app"
)

func ExampleConfig_Clean() {
	c := app.Config{
		SaveName: "testdata///.",
	}
	if err := c.Clean(); err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(os.Stdout, c.SaveName)
	// Output: testdata
}

func ExampleConfig_WalkDir() {
	c := app.Config{
		Print: true,
		Dupes: true,
	}
	if err := c.WalkDir("testdata"); err != nil {
		log.Panicln(err)
	}
	// Output:
	// ── testdata/subdir/test-with-comment.zip
	//    This is an example test comment for zipcmmt.[0m
	//
	//  ── testdata/test-with-comment.zip ───────┐
	//    This is an example test comment for zipcmmt.[0m
}
