package generator_test

import (
	"github.com/stretchr/testify/require"
	"natsu/core"
	"natsu/generator"
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
	generated := generator.Generate(result)
	require.NotContains(t, generated, "PANIC")
}

func TestItGeneratesTheTaggedUnionType(t *testing.T) {
	generated := generator.Generate(result)
	require.Regexp(t, "type TaggedCoolUnion struct", generated)
}

func TestItGeneratesTheConstructor(t *testing.T) {
	generated := generator.Generate(result)
	require.Regexp(t, "CoolUnionOf\\[.* CoolUnion\\]\\(", generated)
}

func TestItGeneratesTheContainers(t *testing.T) {
	generated := generator.Generate(result)
	require.Regexp(t, "type containerCoolUnion interface", generated)
	require.Regexp(t, "type containerTypeA struct", generated)
}
