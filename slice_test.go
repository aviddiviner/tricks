package tricks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstReslicesOriginal(t *testing.T) {
	var numbers = []int{1, 2, 3, 4, 5}

	first3 := Slice(numbers).First(3).Value().([]int)
	first3[0] = 7

	assert.Equal(t, first3, []int{7, 2, 3})
	assert.Equal(t, numbers, []int{7, 2, 3, 4, 5})
}

func TestLastReslicesOriginal(t *testing.T) {
	var numbers = []int{1, 2, 3, 4, 5}

	last3 := Slice(numbers).Last(3).Value().([]int)
	last3[0] = 7

	assert.Equal(t, last3, []int{7, 4, 5})
	assert.Equal(t, numbers, []int{1, 2, 7, 4, 5})
}

func TestCopyPreservesOriginal(t *testing.T) {
	var numbers = []int{1, 2, 3, 4}

	last2 := Slice(numbers).Copy().Last(2).Value().([]int)
	last2[0] = 7

	assert.Equal(t, last2, []int{7, 4})
	assert.Equal(t, numbers, []int{1, 2, 3, 4})
}

func TestGroupBy(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	grouped := Slice(animals).
		GroupBy(func(s string) int { return len(s) }).
		Value().(map[int][]string)

	t.Log(grouped) // TODO
}
