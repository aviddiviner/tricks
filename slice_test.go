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
	assert.Equal(t, 0, len(Slice().Value().([]interface{})))
	assert.Equal(t, 1, len(Slice(nil).Value().([]interface{})))
	assert.Equal(t, 3, len(Slice(nil, nil, nil).Value().([]interface{})))
}

func TestSliceInterface(t *testing.T) {
	actual := Slice([]interface{}{1, "2", 3.0, nil}).Value()
	expected := []interface{}{1, "2", 3.0, nil}
	assert.EqualValues(t, expected, actual)

	values := actual.([]interface{})
	assert.Equal(t, expected[0], values[0])
	assert.Equal(t, expected[1], values[1])
	assert.Equal(t, expected[2], values[2])
	assert.Equal(t, expected[3], values[3])

	// These type assertions should all work.
	assert.Equal(t, 5, len(Slice(1, 2, 3.0, 4, "5").Value().([]interface{})))
	assert.Equal(t, 4, len(Slice(1, 2, 3, nil).Value().([]interface{})))
	assert.Equal(t, 3, len(Slice(Slice(1, 2, 3), 4, 5).Value().([]interface{})))
	assert.Equal(t, 2, len(Slice(nil, struct{ string }{""}).Value().([]interface{})))
}

func TestSlicePanics(t *testing.T) {
	t.Log("TODO")
	// assert.Panics(t, func() { Slice(...) }) // TODO
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

func TestSliceOfSlices(t *testing.T) {
	ss := Slice([]int{1, 2}, []int{3, 4})
	assert.Equal(t, [][]int{[]int{1, 2}, []int{3, 4}}, ss.Value().([][]int))
}

// TODO: Get this to pass:
// func TestSliceOfTrickSlices(t *testing.T) {
// 	ss := Slice(Slice(1, 2), Slice(3, 4))
// 	assert.Equal(t, [][]int{[]int{1, 2}, []int{3, 4}}, ss.Value().([][]int))
// }

func TestSliceOfChans(t *testing.T) {
	a := make(chan struct{})
	b := make(chan struct{})
	assert.Equal(t, []chan struct{}{a, b}, Slice(a, b).Value().([]chan struct{}))
}

func TestSliceOfStructLiterals(t *testing.T) {
	ss := Slice(struct{ string }{"abc"}, struct{ string }{"def"})
	expected := []struct{ string }{struct{ string }{"abc"}, struct{ string }{"def"}}
	assert.Equal(t, expected, ss.Value().([]struct{ string }))

	ss2 := Slice(struct {
		string
		int
	}{"abc", 123}, struct {
		string
		int
	}{"def", 456})

	expected2 := []struct {
		string
		int
	}{
		struct {
			string
			int
		}{"abc", 123},
		struct {
			string
			int
		}{"def", 456},
	}

	assert.Equal(t, expected2, ss2.Value().([]struct {
		string
		int
	}))
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
	assert.Equal(t, 1, Slice(nil).Len())
	assert.Equal(t, 0, Slice().Len())
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

func TestAny(t *testing.T) {
	creatures := Slice("ant", "bear", "cat")
	assert.True(t, creatures.Any(func(word string) bool { return len(word) < 4 }))
	assert.True(t, creatures.Any(func(word string) bool { return len(word) <= 4 }))
	assert.True(t, creatures.Any(func(word string) bool { return len(word) == 4 }))
	assert.False(t, creatures.Any(func(word string) bool { return len(word) > 4 }))
}

func TestAll(t *testing.T) {
	creatures := Slice("ant", "bear", "cat")
	assert.False(t, creatures.All(func(word string) bool { return len(word) < 4 }))
	assert.True(t, creatures.All(func(word string) bool { return len(word) <= 4 }))
	assert.False(t, creatures.All(func(word string) bool { return len(word) == 4 }))
	assert.False(t, creatures.All(func(word string) bool { return len(word) > 4 }))
}

func TestNone(t *testing.T) {
	creatures := Slice("ant", "bear", "cat")
	assert.False(t, creatures.None(func(word string) bool { return len(word) < 4 }))
	assert.False(t, creatures.None(func(word string) bool { return len(word) <= 4 }))
	assert.False(t, creatures.None(func(word string) bool { return len(word) == 4 }))
	assert.True(t, creatures.None(func(word string) bool { return len(word) > 4 }))
}

func TestOne(t *testing.T) {
	creatures := Slice("ant", "bear", "cat")
	assert.False(t, creatures.One(func(word string) bool { return len(word) < 4 }))
	assert.False(t, creatures.One(func(word string) bool { return len(word) <= 4 }))
	assert.True(t, creatures.One(func(word string) bool { return len(word) == 4 }))
	assert.False(t, creatures.One(func(word string) bool { return len(word) > 4 }))
}

func TestMany(t *testing.T) {
	creatures := Slice("ant", "bear", "cat")
	assert.True(t, creatures.Many(func(word string) bool { return len(word) < 4 }))
	assert.True(t, creatures.Many(func(word string) bool { return len(word) <= 4 }))
	assert.False(t, creatures.Many(func(word string) bool { return len(word) == 4 }))
	assert.False(t, creatures.Many(func(word string) bool { return len(word) > 4 }))
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

func TestJoin(t *testing.T) {
	foodChain := Slice("dog", "cat", "bear", "cow").Join(" eats ")
	assert.Equal(t, "dog eats cat eats bear eats cow", foodChain)

	assert.Panics(t, func() { Slice(1, 2, 3).Join("") })
	assert.Panics(t, func() { Slice('a', 'b', 'c').Join("") }) // TODO: Allow joining []rune.
}
