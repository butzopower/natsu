package union_test

import (
	"github.com/stretchr/testify/require"
	"natsu/examples/simple/union"
	"natsu/examples/simple/union/nested"
	"testing"
)

func TestInstantiation(t *testing.T) {
	var tagOfA union.TaggedUnion
	var tagOfB union.TaggedUnion

	tagOfA = union.UnionOf(union.A{})
	tagOfB = union.UnionOf(nested.B{})

	require.NotNil(t, tagOfA)
	require.NotNil(t, tagOfB)
}

func TestShouldNotCompileWithTypeOutsideUnion(t *testing.T) {
	// uncomment to show code has errors

	// var shouldNotWork union.TaggedUnion
	// shouldNotWork = union.UnionOf(union.C{})
	// require.NotNil(t, shouldNotWork)
}
