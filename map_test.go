package tricks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKeysAndValues(t *testing.T) {
	var alphabet = map[string]string{
		"A": "Apple",
		"B": "Ball",
		"C": "Cat",
		"D": "Doll",
		"E": "Egg",
		"F": "Frog",
	}

	letters := Map(alphabet).Keys().Sort().Value().([]string)
	words := Map(alphabet).Values().Sort().Value().([]string)

	assert.Equal(t, []string{"A", "B", "C", "D", "E", "F"}, letters)
	assert.Equal(t, []string{"Apple", "Ball", "Cat", "Doll", "Egg", "Frog"}, words)
}

func TestMapCopy(t *testing.T) {
	var alphabet = map[string]string{
		"A": "Apple",
		"B": "Ball",
		"C": "Cat",
	}
	orig := Map(alphabet)
	abc := Map(alphabet).Copy()
	delete(alphabet, "A")

	assert.Equal(t, []string{"B", "C"}, orig.Keys().Sort().Value().([]string))
	assert.Equal(t, []string{"A", "B", "C"}, abc.Keys().Sort().Value().([]string))

	underlying := abc.Value().(map[string]string)
	delete(underlying, "C")

	assert.Equal(t, []string{"B", "C"}, orig.Keys().Sort().Value().([]string))
	assert.Equal(t, []string{"A", "B"}, abc.Keys().Sort().Value().([]string))
}

func TestMapLen(t *testing.T) {
	assert.Equal(t, 3, Map(map[int]interface{}{1: nil, 3: nil, 5: nil}).Len())
	assert.Equal(t, 1, Map(map[interface{}]int{nil: 1}).Len())
	assert.Equal(t, 0, Map(map[int]bool{}).Len())
}

func TestMapIsEmpty(t *testing.T) {
	assert.True(t, Map(map[int]bool{}).IsEmpty())
	assert.False(t, Map(map[int]bool{1: true}).IsEmpty())
}

func TestMapNil(t *testing.T) {
	// These type assertions should all work.
	assert.Equal(t, 0, len(Map(map[interface{}]interface{}{}).Value().(map[interface{}]interface{})))
	assert.Equal(t, 1, len(Map(map[interface{}]interface{}{nil: nil}).Value().(map[interface{}]interface{})))
	assert.Equal(t, 2, len(Map(map[interface{}]interface{}{nil: nil, "abc": 123}).Value().(map[interface{}]interface{})))
	assert.Equal(t, 0, len(Map(nil).Value().(map[interface{}]interface{})))
}

func TestMapPanics(t *testing.T) {
	t.Log("TODO")
	// assert.Panics(t, func() { Map(...) }) // TODO
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

	missing := Map(alphabet).Only("X").Value().(map[string]string)
	assert.Equal(t, map[string]string{}, missing)

	empty := Map(alphabet).Only().Value().(map[string]string)
	assert.Equal(t, map[string]string{}, empty)

	assert.Panics(t, func() { Map(alphabet).Only(nil) })
	assert.Panics(t, func() { Map(alphabet).Only(1) })
}

func TestMapHasKeys(t *testing.T) {
	var alphabet = map[string]string{
		"A": "Apple",
		"B": "Ball",
		"C": "Cat",
		"D": "Doll",
		"E": "Egg",
		"F": "Frog",
	}

	assert.True(t, Map(alphabet).HasKeys("A", "B"))
	assert.True(t, Map(alphabet).HasKeys("A", "B", "C", "D", "E", "F"))
	assert.True(t, Map(alphabet).HasKeys("F"))
	assert.False(t, Map(alphabet).HasKeys("f"))
	assert.False(t, Map(alphabet).HasKeys("F", "G"))
}
