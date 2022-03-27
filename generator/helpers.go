package generator

import (
    "github.com/butzopower/natsu/core"
    . "github.com/dave/jennifer/jen"
)

func qualifiedTerm(term core.TermPath) *Statement {
    if term.Pointer {
        return Op("*").Qual(term.Package, term.Local)
    }

    return Qual(term.Package, term.Local)
}