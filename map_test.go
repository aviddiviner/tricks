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

func TestMapOnly(t *testing.T) {
	var alphabet = map[string]string{
		"A": "Apple",
		"B": "Ball",
		"C": "Cat",
		"D": "Doll",
		"E": "Egg",
		"F": "Frog",
	}

	subset := Map(alphabet).Only("A", "C", "F").Value().(map[string]string)
	var expected = map[string]string{
		"A": "Apple",
		"C": "Cat",
		"F": "Frog",
	}
	assert.Equal(t, expected, subset)

	missing := Map(alphabet).Only("G").Value().(map[string]string)
	assert.Equal(t, map[string]string{}, missing)

	empty := Map(alphabet).Only().Value().(map[string]string)
	assert.Equal(t, map[string]string{}, empty)
}
