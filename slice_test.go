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

func TestSliceNil(t *testing.T) {
	// These type assertions should all work.
	assert.Equal(t, 0, len(Slice([]interface{}{}).Value().([]interface{})))
	assert.Equal(t, 1, len(Slice([]interface{}{nil}).Value().([]interface{})))
	assert.Equal(t, 3, len(Slice([]interface{}{nil, nil, nil}).Value().([]interface{})))
	assert.Equal(t, 0, len(Slice(nil).Value().([]interface{})))
	assert.Equal(t, 0, len(Slice(nil, nil, nil).Value().([]interface{})))
}

func TestSlicePanics(t *testing.T) {
	assert.Panics(t, func() { Slice(1, 2, 3, 4, "5") })
	assert.Panics(t, func() { Slice(1, 2, 3, 4, nil) })
	assert.Panics(t, func() { Slice(Slice(1, 2, 3), 4, 5) })
}

func TestSliceSingle(t *testing.T) {
	assert.Equal(t, []int{1}, Slice(1).Value().([]int))
	assert.Equal(t, []string{"abc"}, Slice("abc").Value().([]string))
}

func TestSliceTrickSlice(t *testing.T) {
	first := Slice(1, 2, 3)
	second := Slice(first)
	assert.Equal(t, []int{1, 2, 3}, first.Value().([]int))
	assert.Equal(t, []int{1, 2, 3}, second.Value().([]int))
	assert.Equal(t, first, second)
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

func TestSliceLen(t *testing.T) {
	assert.Equal(t, 4, Slice([]int{4, 3, 2, 1}).Len())
	assert.Equal(t, 4, Slice(4, 3, 2, 1).Len())
	assert.Equal(t, 1, Slice([]interface{}{nil}).Len())
	assert.Equal(t, 0, Slice(nil).Len())
}

func TestSortStringsInPlace(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	sorted := Slice(animals).Sort().Value().([]string)

	assert.Equal(t, []string{"bear", "bull", "cat", "cow", "dog", "iguana", "pig"}, sorted)
	assert.Equal(t, sorted, animals)
}

func TestSortIntsInPlace(t *testing.T) {
	var numbers = []int8{3, 5, 21, 1, 34, 55, 13, 2, 8, 89, 1}

	sorted := Slice(numbers).Sort().Value().([]int8)

	assert.Equal(t, []int8{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}, sorted)
	assert.Equal(t, sorted, numbers)
}

func TestSortFloatsInPlace(t *testing.T) {
	var numbers = []float64{3.5, 21.1, 34.55, 13.2, 8.89, 1}

	sorted := Slice(numbers).Sort().Value().([]float64)

	assert.Equal(t, []float64{1, 3.5, 8.89, 13.2, 21.1, 34.55}, sorted)
	assert.Equal(t, sorted, numbers)
}

type testSortByLen []string

func (t testSortByLen) Len() int           { return len(t) }
func (t testSortByLen) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t testSortByLen) Less(i, j int) bool { return len(t[i]) < len(t[j]) }

type testUnsortable []struct{}
type testUnsortableInts []int

func TestSortInterfaceInPlace(t *testing.T) {
	animals := testSortByLen{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}
	sorted := Slice(animals).Sort().Value().(testSortByLen)
	assert.Equal(t, sorted, animals)
	expected := testSortByLen{"dog", "cat", "cow", "pig", "bear", "bull", "iguana"}
	assert.Equal(t, expected, animals)
}

func TestSortUnsortableTypes(t *testing.T) {
	assert.Panics(t, func() { Slice(testUnsortable{}).Sort() })
	assert.Panics(t, func() { Slice(testUnsortableInts{}).Sort() })
	assert.NotPanics(t, func() { Slice([]int(testUnsortableInts{})).Sort() })

	unsortable := testUnsortableInts{3, 5, 21, 1, 34, 55, 13, 2, 8, 89, 1}
	assert.NotPanics(t, func() { Slice([]int(unsortable)).Sort() })
	assert.Equal(t, []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}, []int(unsortable))
}

func TestAllAny(t *testing.T) {
	threes := []int{3, 6, 9, 12, 15, 18, 21, 24, 27, 30}

	allMod3 := Slice(threes).All(func(n int) bool { return n%3 == 0 })
	allMod2 := Slice(threes).All(func(n int) bool { return n%2 == 0 })
	anyMod2 := Slice(threes).Any(func(n int) bool { return n%2 == 0 })
	anyMod11 := Slice(threes).Any(func(n int) bool { return n%11 == 0 })

	assert.True(t, allMod3)
	assert.False(t, allMod2)
	assert.True(t, anyMod2)
	assert.False(t, anyMod11)
}

func TestReverse(t *testing.T) {
	animals := []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}
	reversed := Slice(animals).Reverse().Value().([]string)
	assert.Equal(t, reversed, animals)
	expected := []string{"iguana", "pig", "bull", "cow", "bear", "cat", "dog"}
	assert.Equal(t, expected, animals)

	assert.Equal(t, []int{}, Slice([]int{}).Reverse().Value().([]int))
	assert.Equal(t, []int{1}, Slice([]int{1}).Reverse().Value().([]int))
	assert.Equal(t, []int{1, 2}, Slice([]int{2, 1}).Reverse().Value().([]int))
	assert.Equal(t, []int{1, 2, 3}, Slice([]int{3, 2, 1}).Reverse().Value().([]int))
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
