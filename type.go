package nanogo

type Type interface{}

type IntType struct{}
type FloatType struct{}
type BoolType struct{}

type FunctionType struct {
	Args   []Type
	Return []Type
}
