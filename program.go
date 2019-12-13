package nanogo

import (
	"io"
)

type Program struct {
	Assignments  []*Assignment
	Declarations []*Declaration
}

func (p *Program) Run(w io.Writer) {
	environments := Environments{Environment{}}
	for _, declaration := range p.Declarations {
		declaration.Run(w, environments)
	}
	for _, assignment := range p.Assignments {
		assignment.Run(w, environments)
	}
	closure := environments.Get("main").(*Closure)
	closure.Function.Body.Run(w, closure.Environments)
}
