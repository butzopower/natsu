package union

import "github.com/butzopower/natsu/examples/simple/union/nested"

//go:generate go run github.com/butzopower/natsu github.com/butzopower/natsu/examples/simple/union union SumType
type union interface {
	A | nested.B
}