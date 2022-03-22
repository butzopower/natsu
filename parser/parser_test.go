package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParsingAUnionReturnsTheTerms(t *testing.T) {
	result, err := Parse("natsu/examples/simple/union.Union")

	require.NoError(t, err)

	var termNames []string

	for _, term := range result.Terms {
		termNames = append(termNames, term.String())
	}

	require.ElementsMatch(t, termNames, []string{
		"natsu/examples/simple/union.A",
		"natsu/examples/simple/union.B",
	})
}

func TestParsingReturnsThePath(t *testing.T) {
	result, err := Parse("natsu/examples/simple/union.Union")

	require.NoError(t, err)

	require.Equal(t, "natsu/examples/simple/union", result.Path)
}
