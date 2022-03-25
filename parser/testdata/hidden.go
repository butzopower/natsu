package testdata

import "github.com/butzopower/natsu/parser/testdata/nested"

type hidden interface {
	A | nested.B
}
