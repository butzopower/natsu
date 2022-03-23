package core

type Result struct {
	Union TermPath
	Path  string
	Terms []TermPath
}

type TermPath struct {
	Local   string
	Package string
	Full    string
}
