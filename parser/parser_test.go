package parser

import (
	"github.com/stretchr/testify/require"
	"natsu/core"
	"testing"
)

func TestParsingAUnionReturnsTheUnion(t *testing.T) {
	result, err := Parse("natsu/examples/simple/union.Union")

	require.NoError(t, err)

	require.Equal(t, "natsu/examples/simple/union.Union", result.Union.Full)
	require.Equal(t, "natsu/examples/simple/union", result.Union.Package)
	require.Equal(t, "Union", result.Union.Local)
}

func TestParsingAUnionReturnsTheTerms(t *testing.T) {
	result, err := Parse("natsu/examples/simple/union.Union")

	require.NoError(t, err)

	require.ElementsMatch(t, result.Terms, []core.TermPath{
		{
			Full:    "natsu/examples/simple/union.A",
			Package: "natsu/examples/simple/union",
			Local:   "A",
		},
		{
			Full:    "natsu/examples/simple/union.B",
			Package: "natsu/examples/simple/union",
			Local:   "B",
		},
	})
}

func TestParsingReturnsThePath(t *testing.T) {
	result, err := Parse("natsu/examples/simple/union.Union")

	require.NoError(t, err)

	require.Equal(t, "natsu/examples/simple/union", result.Path)
}
