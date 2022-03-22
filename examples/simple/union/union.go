package union

type A struct{}

type B struct{}

//go:generate go run natsu natsu/examples/simple/union.Union
type Union interface {
	A | B
}
