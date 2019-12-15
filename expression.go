package nanogo

import (
	"io"
)

type Expression interface {
	// Evaluate returns int64, float64, bool or *Closure.
	Evaluate(io.Writer, Environments) interface{}
}

type Int struct{ Value int64 }
type Float struct{ Value float64 }
type Bool struct{ Value bool }
type Variable struct{ Name string }
type Add struct{ Left, Right Expression }
type Sub struct{ Left, Right Expression }
type Mul struct{ Left, Right Expression }
type Div struct{ Left, Right Expression }
type Equal struct{ Left, Right Expression }
type Not struct{ Inner Expression }
type LessThan struct{ Left, Right Expression }

type Application struct {
	Function Expression
	Args     []Expression
}

type Function struct {
	Type Type
	Args []string
	Body []Statement
}

type Closure struct {
	Function     *Function
	Environments Environments
}

func (e *Int) Evaluate(w io.Writer, environments Environments) interface{} {
	return e.Value
}

func (e *Float) Evaluate(w io.Writer, environments Environments) interface{} {
	return e.Value
}

func (e *Bool) Evaluate(w io.Writer, environments Environments) interface{} {
	return e.Value
}

func (e *Variable) Evaluate(w io.Writer, environments Environments) interface{} {
	return environments.Get(e.Name)
}

func (e *Add) Evaluate(w io.Writer, environments Environments) interface{} {
	if left, ok := e.Left.Evaluate(w, environments).(int64); ok {
		return left + e.Right.Evaluate(w, environments).(int64)
	}

	return e.Left.Evaluate(w, environments).(float64) + e.Right.Evaluate(w, environments).(float64)
}

func (e *Sub) Evaluate(w io.Writer, environments Environments) interface{} {
	if left, ok := e.Left.Evaluate(w, environments).(int64); ok {
		return left - e.Right.Evaluate(w, environments).(int64)
	}

	return e.Left.Evaluate(w, environments).(float64) - e.Right.Evaluate(w, environments).(float64)
}

func (e *Mul) Evaluate(w io.Writer, environments Environments) interface{} {
	if left, ok := e.Left.Evaluate(w, environments).(int64); ok {
		return left * e.Right.Evaluate(w, environments).(int64)
	}

	return e.Left.Evaluate(w, environments).(float64) * e.Right.Evaluate(w, environments).(float64)
}

func (e *Div) Evaluate(w io.Writer, environments Environments) interface{} {
	if left, ok := e.Left.Evaluate(w, environments).(int64); ok {
		return left / e.Right.Evaluate(w, environments).(int64)
	}

	return e.Left.Evaluate(w, environments).(float64) / e.Right.Evaluate(w, environments).(float64)
}

func (e *Equal) Evaluate(w io.Writer, environments Environments) interface{} {
	if left, ok := e.Left.Evaluate(w, environments).(int64); ok {
		return left == e.Right.Evaluate(w, environments).(int64)
	}

	return e.Left.Evaluate(w, environments).(float64) == e.Right.Evaluate(w, environments).(float64)
}

func (e *Not) Evaluate(w io.Writer, environments Environments) interface{} {
	return !e.Inner.Evaluate(w, environments).(bool)
}

func (e *LessThan) Evaluate(w io.Writer, environments Environments) interface{} {
	if left, ok := e.Left.Evaluate(w, environments).(int64); ok {
		return left < e.Right.Evaluate(w, environments).(int64)
	}

	return e.Left.Evaluate(w, environments).(float64) < e.Right.Evaluate(w, environments).(float64)
}

func (e *Application) Evaluate(w io.Writer, environments Environments) interface{} {
	closure := e.Function.Evaluate(w, environments).(*Closure)
	environments = append(closure.Environments, Environment{})
	for i, arg := range e.Args {
		environments.Add(closure.Function.Args[i], arg.Evaluate(w, environments))
	}
	for _, statement := range closure.Function.Body {
		if v := statement.Run(w, environments); v != nil {
			return v
		}
	}
	return nil
}

func (e *Function) Evaluate(w io.Writer, environments Environments) interface{} {
	return &Closure{e, environments}
}
