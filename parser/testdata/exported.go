package testdata

import "github.com/butzopower/natsu/parser/testdata/nested"

type Exported interface {
	A | *nested.B
}
