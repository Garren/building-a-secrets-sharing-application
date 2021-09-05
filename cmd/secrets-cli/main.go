package main

import (
	"fmt"
	"os"
)

func main() {
	c, err := setupParseArgs(os.Stdout, os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	errors := validateConfig(c)
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Println(e)
		}
		os.Exit(1)
	}
	result, err := performAction(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, result)
}
