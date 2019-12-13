package nanogo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironments(t *testing.T) {
	environments := Environments{Environment{}}

	environments.Add("foo", 1)
	assert.Equal(t, 1, environments.Get("foo"))

	environments = append(environments, Environment{})
	assert.Equal(t, 1, environments.Get("foo"))

	environments.Add("foo", 2)
	assert.Equal(t, 2, environments.Get("foo"))

	environments.Set("foo", 3)
	assert.Equal(t, 3, environments.Get("foo"))

	environments = environments[:len(environments)-1]
	assert.Equal(t, 1, environments.Get("foo"))
}
