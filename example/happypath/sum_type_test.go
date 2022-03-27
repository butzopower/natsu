package happypath_test

import (
    "github.com/butzopower/natsu/example/happypath"
    "github.com/butzopower/natsu/example/happypath/nested"
    "github.com/stretchr/testify/require"
    "testing"
)

func TestInstantiation(t *testing.T) {
    var tagOfA happypath.SumType
    var tagOfB happypath.SumType

    tagOfA = happypath.SumTypeOf(happypath.A{})
    tagOfB = happypath.SumTypeOf(nested.B{})

    require.NotNil(t, tagOfA)
    require.NotNil(t, tagOfB)
}

func TestExecutor(t *testing.T) {
    calledWith := "none"

    executor := happypath.SumTypeExecutor().
        WithA(func(a happypath.A) { calledWith = "a" }).
        WithB(func(b nested.B) { calledWith = "b" })

    executor.Exec(happypath.SumTypeOf(happypath.A{}))

    require.Equal(t, "a", calledWith)

    executor.Exec(happypath.SumTypeOf(nested.B{}))

    require.Equal(t, "b", calledWith)
}

func TestMapper(t *testing.T) {
    mapper := happypath.SumTypeMapper[string]().
        WithA(func(a happypath.A) string { return "a" }).
        WithB(func(b nested.B) string { return "b" })

    mapper.Map(happypath.SumTypeOf(happypath.A{}))

    require.Equal(
        t,
        "a",
        mapper.Map(happypath.SumTypeOf(happypath.A{})),
    )

    require.Equal(
        t,
        "b",
        mapper.Map(happypath.SumTypeOf(nested.B{})),
    )
}

func TestShouldNotCompileWithTypeOutsideSumType(t *testing.T) {
    // uncomment to show code has errors
    // var shouldNotWork happypath.SumType
    // shouldNotWork = happypath.SumTypeOf(happypath.C{})
    // require.NotNil(t, shouldNotWork)
}
