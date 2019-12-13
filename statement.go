package nanogo

import (
	"fmt"
	"io"
)

type Statement interface {
	Run(io.Writer, Environments) interface{}
}

type Block []Statement

type Declaration struct {
	Name string
	Type Type
}

type Assignment struct {
	Left  string
	Right Expression
}

type If struct {
	Condition Expression
	Block     Block
}

type For struct {
	Condition Expression
	Block     Block
}

type Return struct {
	Expression Expression
}

type Print struct{ Arg Expression }

func (s Block) Run(w io.Writer, environments Environments) interface{} {
	environments = append(environments, Environment{})
	for _, statement := range s {
		if v := statement.Run(w, environments); v != nil {
			return v
		}
	}

	return nil
}

func (s *Declaration) Run(w io.Writer, environments Environments) interface{} {
	switch s.Type.(type) {
	case *IntType:
		environments.Add(s.Name, int64(0))
	case *FloatType:
		environments.Add(s.Name, float64(0))
	case *BoolType:
		environments.Add(s.Name, bool(false))
	case *FunctionType:
		environments.Add(s.Name, new(Closure))
	}

	return nil
}

func (s *Assignment) Run(w io.Writer, environments Environments) interface{} {
	environments.Set(s.Left, s.Right.Evaluate(w, environments))
	return nil
}

func (s *If) Run(w io.Writer, environments Environments) interface{} {
	if s.Condition.Evaluate(w, environments).(bool) {
		return s.Block.Run(w, environments)
	}
	return nil
}

func (s *For) Run(w io.Writer, environments Environments) interface{} {
	for s.Condition.Evaluate(w, environments).(bool) {
		if v := s.Block.Run(w, environments); v != nil {
			return v
		}
	}
	return nil
}

func (s *Return) Run(w io.Writer, environments Environments) interface{} {
	return s.Expression.Evaluate(w, environments)
}

func (s *Print) Run(w io.Writer, environments Environments) interface{} {
	fmt.Fprint(w, s.Arg.Evaluate(w, environments))
	return nil
}

func (s *Application) Run(w io.Writer, environments Environments) interface{} {
	s.Evaluate(w, environments)
	return nil
}
