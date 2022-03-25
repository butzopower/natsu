package happypath

import (
	"github.com/butzopower/natsu/example/happypath/nested"
)

//go:generate go run github.com/butzopower/natsu github.com/butzopower/natsu/example/happypath union SumType
type union interface {
	A | nested.B
}
