//go:build windows

// © Ben Garrett https://github.com/bengarrett/zipcmt

package main_test

import (
	"fmt"
	"log"
	"os"

	zipcmt "github.com/bengarrett/zipcmt/pkg"
)

func ExampleConfig_clean() {
	c := zipcmt.Config{
		SaveName: "..//test///.",
	}
	if err := c.Clean(); err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(os.Stdout, c.SaveName)
	// Output: ..\test
}

func ExampleConfig_WalkDir() {
	c := zipcmt.Config{
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
