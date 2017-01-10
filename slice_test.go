package tricks

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceValid(t *testing.T) {
	normal := Slice([]int{1, 2, 3, 4, 5}).Value().([]int)
	variadic := Slice(1, 2, 3, 4, 5).Value().([]int)

	assert.Equal(t, []int{1, 2, 3, 4, 5}, normal)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, variadic)
	assert.Equal(t, 5, cap(normal))
	assert.Equal(t, 5, cap(variadic))
}

func TestSlicePanics(t *testing.T) {
	assert.Panics(t, func() { Slice(1, 2, 3, 4, "5") })
	assert.Panics(t, func() { Slice(1, 2, 3, 4, nil) })
}

func TestSliceSingle(t *testing.T) {
	assert.Equal(t, []int{1}, Slice(1).Value().([]int))
	assert.Equal(t, []string{"abc"}, Slice("abc").Value().([]string))
}

func TestFirstReslicesOriginal(t *testing.T) {
	var numbers = make([]int, 10)[:5]
	copy(numbers, []int{1, 2, 3, 4, 5})

	first3 := Slice(numbers).First(3).Value().([]int)
	first3[0] = 7

	assert.Equal(t, []int{7, 2, 3}, first3)
	assert.Equal(t, 3, cap(first3))
	assert.Equal(t, []int{7, 2, 3, 4, 5}, numbers)
	assert.Equal(t, 10, cap(numbers))
}

func TestLastReslicesOriginal(t *testing.T) {
	var numbers = make([]int, 10)[:5]
	copy(numbers, []int{1, 2, 3, 4, 5})

	last3 := Slice(numbers).Last(3).Value().([]int)
	last3[0] = 7

	assert.Equal(t, []int{7, 4, 5}, last3)
	assert.Equal(t, 3, cap(last3))
	assert.Equal(t, []int{1, 2, 7, 4, 5}, numbers)
	assert.Equal(t, 10, cap(numbers))
}

func TestSliceCopyPreservesOriginal(t *testing.T) {
	var numbers = []int{4, 3, 2, 1}

	sorted := Slice(numbers).Copy().Sort().Value().([]int)

	assert.Equal(t, []int{4, 3, 2, 1}, numbers)
	assert.NotEqual(t, numbers, sorted)
	assert.Equal(t, []int{1, 2, 3, 4}, sorted)
}

func TestSortStringsInPlace(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	sorted := Slice(animals).Sort().Value().([]string)

	assert.Equal(t, []string{"bear", "bull", "cat", "cow", "dog", "iguana", "pig"}, sorted)
	assert.Equal(t, sorted, animals)
}

func TestSortIntsInPlace(t *testing.T) {
	var numbers = []int{3, 5, 21, 1, 34, 55, 13, 2, 8, 89, 1}

	sorted := Slice(numbers).Sort().Value().([]int)

	assert.Equal(t, []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}, sorted)
	assert.Equal(t, sorted, numbers)
}

func TestSortFloatsInPlace(t *testing.T) {
	var numbers = []float64{3.5, 21.1, 34.55, 13.2, 8.89, 1}

	sorted := Slice(numbers).Sort().Value().([]float64)

	assert.Equal(t, []float64{1, 3.5, 8.89, 13.2, 21.1, 34.55}, sorted)
	assert.Equal(t, sorted, numbers)
}

type testSortByLen []struct{ string }
type testUnsortable []struct{ string }

func (t testSortByLen) Len() int           { return len(t) }
func (t testSortByLen) Swap(i, j int)      { t[i].string, t[j].string = t[j].string, t[i].string }
func (t testSortByLen) Less(i, j int) bool { return len(t[i].string) < len(t[j].string) }

func TestSortInterfaceInPlace(t *testing.T) {
	var animals = testSortByLen{
		struct{ string }{"dog"},
		struct{ string }{"cat"},
		struct{ string }{"bear"},
		struct{ string }{"cow"},
		struct{ string }{"bull"},
		struct{ string }{"pig"},
		struct{ string }{"iguana"},
	}

	sorted := Slice(animals).Sort().Value().(testSortByLen)

	var expected = testSortByLen{
		struct{ string }{"dog"},
		struct{ string }{"cat"},
		struct{ string }{"cow"},
		struct{ string }{"pig"},
		struct{ string }{"bear"},
		struct{ string }{"bull"},
		struct{ string }{"iguana"},
	}
	assert.Equal(t, expected, sorted)
	assert.Equal(t, sorted, animals)
}

func TestGroupBy(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	grouped := Slice(animals).
		GroupBy(func(s string) int { return len(s) }).
		Value().(map[int][]string)

	expected := map[int][]string{
		3: []string{"dog", "cat", "cow", "pig"},
		4: []string{"bear", "bull"},
		6: []string{"iguana"},
	}
	assert.Equal(t, expected, grouped)
}

func TestSliceMap(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow"}

	result := Slice(animals).Map(strings.ToUpper).Value().([]string)

	expected := []string{"DOG", "CAT", "BEAR", "COW"}
	assert.Equal(t, expected, result)
	assert.NotEqual(t, animals, result)
}

func TestMaxAndMin(t *testing.T) {
	var numbers = []int{3, 5, 21, 1, 34, 55, 13, 2, 8, 89, 1}
	var number = []int{42}

	max := Slice(numbers).Max().(int)
	assert.Equal(t, 89, max)

	max = Slice(number).Max().(int)
	assert.Equal(t, 42, max)

	_, ok := Slice([]int{}).Max().(int)
	assert.False(t, ok)

	min := Slice(numbers).Min().(int)
	assert.Equal(t, 1, min)

	min = Slice(number).Min().(int)
	assert.Equal(t, 42, min)

	_, ok = Slice([]int{}).Min().(int)
	assert.False(t, ok)
}
