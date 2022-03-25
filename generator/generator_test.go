package generator_test

import (
	"github.com/butzopower/natsu/core"
	"github.com/butzopower/natsu/generator"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var sumTypeName = "CoolUnion"
var result = core.UnionDetails{
	Path: "some/pkg",
	Union: core.TermPath{
		Package: "some/pkg",
		Local:   "Union",
	},
	Terms: []core.TermPath{
		{
			Local:   "TypeA",
			Package: "some/pkg",
		},
	},
}

func TestItGeneratesAValidFile(t *testing.T) {
	generate(t)
}

func generate(t *testing.T) string {
	var builder strings.Builder

	file := generator.Generate(sumTypeName, result)

	err := file.Render(&builder)

	require.NoError(t, err)

	return builder.String()
}
