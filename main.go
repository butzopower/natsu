package main

import (
	"fmt"
	"github.com/butzopower/natsu/generator"
	"github.com/butzopower/natsu/parser"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		failErr(fmt.Errorf("expected exactly two arguments: <pkg> <original union type> <name to generate>"))
	}
	sourcePkg := os.Args[1]
	sourceType := os.Args[2]
	nameToGenerate := os.Args[3]

	unionDetails, err := parser.Parse(sourcePkg, sourceType)

	if err != nil {
		failErr(err)
	}

	sourceFile := generator.Generate(nameToGenerate, unionDetails)

	goFile := os.Getenv("GOFILE")
	dir := filepath.Dir(goFile)

	targetFilename := filepath.Join(
		dir,
		strings.ToLower(nameToGenerate)+".go",
	)

	failErr(sourceFile.Save(targetFilename))
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
