package union_test

import (
	"github.com/butzopower/natsu/examples/simple/union"
	"github.com/butzopower/natsu/examples/simple/union/nested"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInstantiation(t *testing.T) {
	var tagOfA union.SumType
	var tagOfB union.SumType

	tagOfA = union.SumTypeOf(union.A{})
	tagOfB = union.SumTypeOf(nested.B{})

	require.NotNil(t, tagOfA)
	require.NotNil(t, tagOfB)
}

func TestExecutor(t *testing.T) {
	calledWith := "none"

	executor := union.SumTypeExecutor().
		WithA(func(a union.A) { calledWith = "a" }).
		WithB(func(b nested.B) { calledWith = "b" })

	executor.Exec(union.SumTypeOf(union.A{}))

	require.Equal(t, "a", calledWith)

	executor.Exec(union.SumTypeOf(nested.B{}))

	require.Equal(t, "b", calledWith)
}

func TestShouldNotCompileWithTypeOutsideSumType(t *testing.T) {
	// uncomment to show code has errors

	// var shouldNotWork union.SumType
	// shouldNotWork = union.SumTypeOf(union.C{})
	// require.NotNil(t, shouldNotWork)
}
