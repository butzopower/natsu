package union_test

import (
	"github.com/butzopower/natsu/examples/simple/union"
	"github.com/butzopower/natsu/examples/simple/union/nested"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInstantiation(t *testing.T) {
	var tagOfA union.Union
	var tagOfB union.Union

	tagOfA = union.UnionOf(union.A{})
	tagOfB = union.UnionOf(nested.B{})

	require.NotNil(t, tagOfA)
	require.NotNil(t, tagOfB)
}

func TestExecutor(t *testing.T) {
	calledWith := "none"

	executor := union.UnionExecutor().
		WithA(func(a union.A) { calledWith = "a" }).
		WithB(func(b nested.B) { calledWith = "b" })

	executor.Exec(union.UnionOf(union.A{}))

	require.Equal(t, "a", calledWith)

	executor.Exec(union.UnionOf(nested.B{}))

	require.Equal(t, "b", calledWith)
}

func TestShouldNotCompileWithTypeOutsideUnion(t *testing.T) {
	// uncomment to show code has errors

	// var shouldNotWork union.Union
	// shouldNotWork = union.UnionOf(union.C{})
	// require.NotNil(t, shouldNotWork)
}
