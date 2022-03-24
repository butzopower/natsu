package generator_test

import (
	"github.com/butzopower/natsu/core"
	"github.com/butzopower/natsu/generator"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var result = core.Result{
	Path: "some/pkg",
	Union: core.TermPath{
		Full:    "some/pkg.CoolUnion",
		Package: "some/pkg",
		Local:   "CoolUnion",
	},
	Terms: []core.TermPath{
		{
			Local:   "TypeA",
			Package: "some/pkg",
			Full:    "some/pkg.TypeA",
		},
	},
}

func TestItGeneratesAValidFile(t *testing.T) {
	generate(t)
}

func TestItGeneratesTheTaggedUnionType(t *testing.T) {
	generated := generate(t)
	require.Regexp(t, "type TaggedCoolUnion struct", generated)
}

func TestItGeneratesTheConstructor(t *testing.T) {
	generated := generate(t)
	require.Regexp(t, "CoolUnionOf\\[.* CoolUnion\\]\\(", generated)
}

func TestItGeneratesTheContainers(t *testing.T) {
	generated := generate(t)
	require.Regexp(t, "type containerCoolUnion interface", generated)
	require.Regexp(t, "type containerTypeA struct", generated)
}

func generate(t *testing.T) string {
	var builder strings.Builder

	file := generator.Generate(result)

	err := file.Render(&builder)

	require.NoError(t, err)

	return builder.String()
}
