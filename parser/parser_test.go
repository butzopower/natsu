package parser

import (
	"github.com/butzopower/natsu/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParsingAUnionReturnsTheUnion(t *testing.T) {
	result, err := Parse("github.com/butzopower/natsu/examples/simple/union", "U")

	require.NoError(t, err)

	require.Equal(t, "github.com/butzopower/natsu/examples/simple/union", result.Union.Package)
	require.Equal(t, "U", result.Union.Local)
}

func TestParsingAUnionReturnsTheTerms(t *testing.T) {
	result, err := Parse("github.com/butzopower/natsu/examples/simple/union", "U")

	require.NoError(t, err)

	require.ElementsMatch(t, result.Terms, []core.TermPath{
		{
			Package: "github.com/butzopower/natsu/examples/simple/union",
			Local:   "A",
		},
		{
			Package: "github.com/butzopower/natsu/examples/simple/union/nested",
			Local:   "B",
		},
	})
}

func TestParsingReturnsThePath(t *testing.T) {
	result, err := Parse("github.com/butzopower/natsu/examples/simple/union", "U")

	require.NoError(t, err)

	require.Equal(t, "github.com/butzopower/natsu/examples/simple/union", result.Path)
}
