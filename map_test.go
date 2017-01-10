package tricks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKeys(t *testing.T) {
	var alphabet = map[string]string{
		"A": "Apple",
		"B": "Ball",
		"C": "Cat",
		"D": "Doll",
		"E": "Egg",
		"F": "Frog",
	}

	letters := Map(alphabet).Keys().Sort().Value().([]string)

	assert.Equal(t, []string{"A", "B", "C", "D", "E", "F"}, letters)
}
