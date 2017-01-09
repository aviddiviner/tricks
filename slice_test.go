package tricks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstReslicesOriginal(t *testing.T) {
	var numbers = []int{1, 2, 3, 4, 5}

	first3 := Slice(numbers).First(3).Value().([]int)
	first3[0] = 7

	assert.Equal(t, []int{7, 2, 3}, first3)
	assert.Equal(t, []int{7, 2, 3, 4, 5}, numbers)
}

func TestLastReslicesOriginal(t *testing.T) {
	var numbers = []int{1, 2, 3, 4, 5}

	last3 := Slice(numbers).Last(3).Value().([]int)
	last3[0] = 7

	assert.Equal(t, []int{7, 4, 5}, last3)
	assert.Equal(t, []int{1, 2, 7, 4, 5}, numbers)
}

func TestCopyPreservesOriginal(t *testing.T) {
	var numbers = []int{1, 2, 3, 4}

	last2 := Slice(numbers).Copy().Last(2).Value().([]int)
	last2[0] = 7

	assert.Equal(t, []int{7, 4}, last2)
	assert.Equal(t, []int{1, 2, 3, 4}, numbers)
}

func TestSortStringsInPlace(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	sorted := Slice(animals).Sort().Value().([]string)

	assert.Equal(t, []string{"bear", "bull", "cat", "cow", "dog", "iguana", "pig"}, sorted)
	assert.Equal(t, []string{"bear", "bull", "cat", "cow", "dog", "iguana", "pig"}, animals)
}

func TestSortIntsInPlace(t *testing.T) {
	var numbers = []int{3, 5, 21, 1, 34, 55, 13, 2, 8, 89, 1}

	sorted := Slice(numbers).Sort().Value().([]int)

	assert.Equal(t, []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}, sorted)
	assert.Equal(t, []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}, numbers)
}

func TestSortFloatsInPlace(t *testing.T) {
	var numbers = []float64{3.5, 21.1, 34.55, 13.2, 8.89, 1}

	sorted := Slice(numbers).Sort().Value().([]float64)

	assert.Equal(t, []float64{1, 3.5, 8.89, 13.2, 21.1, 34.55}, sorted)
	assert.Equal(t, []float64{1, 3.5, 8.89, 13.2, 21.1, 34.55}, numbers)
}

func TestGroupBy(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	grouped := Slice(animals).
		GroupBy(func(s string) int { return len(s) }).
		Value().(map[int][]string)

	t.Log(grouped) // TODO
}
