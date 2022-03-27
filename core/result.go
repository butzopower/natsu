package core

type UnionDetails struct {
	Union TermPath
	Path  string
	Terms []TermPath
}

type TermPath struct {
	Local   string
	Package string
	Pointer bool
}
