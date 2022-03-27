package happypath_test

import (
	"fmt"
	"github.com/butzopower/natsu/example/happypath"
	"github.com/butzopower/natsu/example/happypath/nested"
	"testing"
)

func BenchmarkExecutorCreateEachTime(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkExecutorA(happypath.A{})
		benchmarkExecutorB(&nested.B{})
	}
}

func BenchmarkExecutorCreateOnce(b *testing.B) {
	executor := createExecutor()

	for n := 0; n < b.N; n++ {
		executor.Exec(happypath.SumTypeOf(happypath.A{}))
		executor.Exec(happypath.SumTypeOf(&nested.B{}))
	}
}

func BenchmarkExecutorGoImpl(b *testing.B) {
	for n := 0; n < b.N; n++ {
		executorGoImpl(happypath.A{})
		executorGoImpl(&nested.B{})
	}
}

func createExecutor() *happypath.ExecutorSumType {
	return happypath.SumTypeExecutor().
		WithA(func(a happypath.A) { fmt.Sprintf("%+v", a) }).
		WithB(func(b *nested.B) { fmt.Sprintf("%+v", b) })
}

func benchmarkExecutorA(a happypath.A) {
	createExecutor().Exec(happypath.SumTypeOf(a))
}

func benchmarkExecutorB(b *nested.B) {
	createExecutor().Exec(happypath.SumTypeOf(b))
}

type union interface {
	happypath.A | *nested.B
}

func executorGoImpl[T union](u T) {
	switch t := any(u).(type) {
	case happypath.A:
		fmt.Sprintf("%+v", t)
	case *nested.B:
		fmt.Sprintf("%+v", t)
	default:
		panic("should not get here")
	}
}
