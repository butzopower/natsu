package main

import (
	"fmt"
	"natsu/parser"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		failErr(fmt.Errorf("expected exactly one argument: <source type>"))
	}
	sourceType := os.Args[1]

	terms, err := parser.Parse(sourceType)

	if err != nil {
		failErr(err)
	}

	for _, term := range terms {
		fmt.Println(term.String(), term.Type())
	}
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
