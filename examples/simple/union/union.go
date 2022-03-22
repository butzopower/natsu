package union

type A struct{}

type B struct{}

type Union interface {
	A | B
}
