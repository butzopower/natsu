package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParsingAUnionReturnsTheTypes(t *testing.T) {
	terms, err := Parse("natsu/examples/simple/union.Union")

	require.NoError(t, err)

	var termNames []string

	for _, term := range terms {
		termNames = append(termNames, term.String())
	}

	require.ElementsMatch(t, termNames, []string{
		"natsu/examples/simple/union.A",
		"natsu/examples/simple/union.B",
	})
}
