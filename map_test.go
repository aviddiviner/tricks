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

	subset := Map(alphabet).Only(Slice("A", "C", "F")).Value().(map[string]string)
	var expected = map[string]string{
		"A": "Apple",
		"C": "Cat",
		"F": "Frog",
	}
	assert.Equal(t, expected, subset)

	keys := []string{"A", "C", "F"}
	again := Map(alphabet).Only(keys).Value().(map[string]string)
	assert.Equal(t, expected, again)

	single := Map(alphabet).Only("A").Value().(map[string]string)
	assert.Equal(t, map[string]string{"A": "Apple"}, single)

	missing := Map(alphabet).Only("G").Value().(map[string]string)
	assert.Equal(t, map[string]string{}, missing)

	empty := Map(alphabet).Only(nil).Value().(map[string]string)
	assert.Equal(t, map[string]string{}, empty)
}
