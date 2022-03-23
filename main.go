package main

import (
	"fmt"
	"natsu/generator"
	"natsu/parser"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		failErr(fmt.Errorf("expected exactly one argument: <source type>"))
	}
	sourceType := os.Args[1]

	result, err := parser.Parse(sourceType)

	if err != nil {
		failErr(err)
	}

	sourceFile := generator.Generate(result)

	goFile := os.Getenv("GOFILE")
	dir := filepath.Dir(goFile)

	targetFilename := filepath.Join(
		dir,
		strings.ToLower(result.Union.Local)+"_sum.go",
	)

	failErr(sourceFile.Save(targetFilename))
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
