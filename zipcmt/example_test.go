package zipcmt

import (
	"fmt"
	"log"
)

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
