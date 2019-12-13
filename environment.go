package nanogo

type Environment map[string]interface{}
type Environments []Environment

func (e Environments) Get(name string) interface{} {
	for i := len(e) - 1; i >= 0; i-- {
		if v, exists := e[i][name]; exists {
			return v
		}
	}
	return nil
}

func (e Environments) Set(k string, v interface{}) {
	for i := len(e) - 1; i >= 0; i-- {
		if _, exists := e[i][k]; exists {
			e[i][k] = v
			return
		}
	}
}

func (e Environments) Add(k string, v interface{}) {
	e[len(e)-1][k] = v
}
