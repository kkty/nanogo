package test

import (
	"bytes"
	"testing"

	"github.com/kkty/nanogo"
	"github.com/stretchr/testify/assert"
)

func TestParseAndRun(t *testing.T) {
	for _, c := range []struct {
		program  string
		expected string
	}{
		{
			`
			func main() () {
				print(3)
			}
			`,
			"3",
		},
		{
			`
			var cnt int64
			func main() () {
				var i int64
				for i < 10 {
					i = i + 1
					increment()
				}
				print(cnt)
			}

			func increment() () {
				cnt = cnt + 1
			}
			`,
			"10",
		},
		{
			`
			func main() () {
				var add func(int64) (int64)
				add = makeAdder(1)
				print(add(3))
				add = makeAdder(2)
				print(add(4))
			}
			func makeAdder(i int64) (func (int64) (int64)) {
				return func(j int64) (int64) {
					return i + j
				}
			}
			`,
			"46",
		},
	} {
		program := nanogo.Parse(c.program)
		buf := &bytes.Buffer{}
		program.Run(buf)
		assert.Equal(t, c.expected, buf.String())
	}
}
