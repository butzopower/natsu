package parser

import (
	"github.com/butzopower/natsu/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParsingAUnionReturnsTheUnion(t *testing.T) {
	result, err := Parse("github.com/butzopower/natsu/parser/testdata", "Exported")

	require.NoError(t, err)

	require.Equal(t, "github.com/butzopower/natsu/parser/testdata", result.Union.Package)
	require.Equal(t, "Exported", result.Union.Local)
}

func TestParsingAnUnexposedUnionReturnsTheUnion(t *testing.T) {
	result, err := Parse("github.com/butzopower/natsu/parser/testdata", "hidden")

	require.NoError(t, err)

	require.Equal(t, "github.com/butzopower/natsu/parser/testdata", result.Union.Package)
	require.Equal(t, "hidden", result.Union.Local)
}

func TestParsingAUnionReturnsTheTerms(t *testing.T) {
	result, err := Parse("github.com/butzopower/natsu/parser/testdata", "Exported")

	require.NoError(t, err)

	require.ElementsMatch(t, result.Terms, []core.TermPath{
		{
			Package: "github.com/butzopower/natsu/parser/testdata",
			Local:   "A",
			Pointer: false,
		},
		{
			Package: "github.com/butzopower/natsu/parser/testdata/nested",
			Local:   "B",
			Pointer: true,
		},
	})
}

func TestParsingReturnsThePath(t *testing.T) {
	result, err := Parse("github.com/butzopower/natsu/parser/testdata", "Exported")

	require.NoError(t, err)

	require.Equal(t, "github.com/butzopower/natsu/parser/testdata", result.Path)
}
