package main

import (
	"fmt"
	"github.com/butzopower/natsu/generator"
	"github.com/butzopower/natsu/parser"
	"os"
	"path/filepath"
	"regexp"
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
	filename := fmt.Sprintf("%s.go", toSnakeCase(nameToGenerate))

	targetFilename := filepath.Join(dir, filename)

	failErr(sourceFile.Save(targetFilename))
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
