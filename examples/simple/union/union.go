package union

import "natsu/examples/simple/union/nested"

//go:generate go run natsu natsu/examples/simple/union.Union
type Union interface {
	A | nested.B
}
