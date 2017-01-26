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

func TestSliceVariadic(t *testing.T) {
	assert.Equal(t, 3, len(Slice(1, 2, 3).Value().([]int)))

	// These are quirky, but there's no good way around it. The variadic args
	// are always going to be received as []interface{}, so we'll have no idea
	// of knowing if we were passed an actual int, or a interface{}(int).
	assert.Equal(t, 3, len(Slice(interface{}(1), 2, 3).Value().([]int)))
	assert.Equal(t, 3, len(Slice(1, 2, interface{}(3)).Value().([]int)))

	// Also, the nil values are just treated as zeroes of that particular type
	// (in this case, zero valued ints, which are just 0).
	assert.Equal(t, 3, len(Slice(1, 2, nil).Value().([]int)))
	assert.Equal(t, 3, len(Slice(nil, 2, 3).Value().([]int)))

	// So, either we treat nils as zeroes of the shared type, or we're forced
	// to type all variadic constructions containing nil values as []interface{}.
	//
	// In particular, if we make this test case pass:
	//   assert.NotPanics(t, func() { Slice(nil, 2, 3).Value().([]interface{}) })
	//
	// Then this test case will fail:
	//   a, b := 1, 2
	//   assert.NotPanics(t, func() { Slice(&a, &b, nil).Value().([]*int) })
	//
	// And I think we should rather have this latter one pass.
	//
	// TODO: Document this fact: Slice(1, 2, nil) == []int{1, 2, 0}
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

	// Variadic. These have no common type; we resort to []interface{}.
	assert.Equal(t, 5, len(Slice(1, 2, 3.0, 4, "5").Value().([]interface{})))
	assert.Equal(t, 3, len(Slice(Slice(1, 2, 3), 4, 5).Value().([]interface{})))

	// Variadic. These share a common type, with some zero values.
	assert.Equal(t, 4, len(Slice(1, 2, 3, nil).Value().([]int)))
	assert.Equal(t, 2, len(Slice(nil, struct{ string }{""}).Value().([]struct{ string })))

	// Singular. These type assertions should work and it should come out the same type it went in with.

	// [[[1 2] [3]] [[4 5]]]
	assert.Equal(t, 2, len(Slice([][][]int{
		[][]int{[]int{1, 2}, []int{3}},
		[][]int{[]int{4, 5}},
	}).Value().([][][]int)))

	// [[[1 2] [3]] [[4 5]]]
	assert.Equal(t, 2, len(Slice([][][]interface{}{
		[][]interface{}{[]interface{}{1, 2}, []interface{}{3}},
		[][]interface{}{[]interface{}{4, 5}},
	}).Value().([][][]interface{})))
}

func TestSliceInterfaceSplat(t *testing.T) {
	numbers := []interface{}{1, 2, 3, 4, 5}
	ints := Slice(numbers...).Value()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, ints.([]int))
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

func TestSliceOfTrickSlices(t *testing.T) {
	ss := Slice(Slice(1, 2), Slice(3, 4))
	ints := ss.Map(func(ts TrickSlice) []int { return ts.Value().([]int) }).Value().([][]int)
	assert.EqualValues(t, [][]int{[]int{1, 2}, []int{3, 4}}, ints)
}

// TODO: Think about how to do a .DeepVal (or .Flatten) so we can make these kinds of assertions.
//
// ss := Slice(Slice(1, 2), Slice(3, 4))
// assert.Equal(t, [][]int{[]int{1, 2}, []int{3, 4}}, ss.DeepVal().([][]int))
//
// sss := Slice(Slice(Slice(1, 2), Slice(3, 4)), Slice(Slice(5), Slice(6)))
// exp := [][][]int{
// 	[][]int{[]int{1, 2}, []int{3, 4}},
// 	[][]int{[]int{5}, []int{6}},
// }
// assert.Equal(t, exp, sss.DeepVal().([][][]int))

func TestSliceOfChans(t *testing.T) {
	a := make(chan struct{})
	b := make(chan struct{})
	assert.Equal(t, []chan struct{}{a, b}, Slice(a, b).Value().([]chan struct{}))
}

func TestSliceOfPointers(t *testing.T) {
	a, b := 3, 4

	ptrSlice := Slice(&a, &b)
	expected := []*int{&a, &b}
	assert.Equal(t, expected, ptrSlice.Value().([]*int))

	ptrNilSlice := Slice(&a, &b, nil)
	expectedNil := []*int{&a, &b, nil}
	assert.Equal(t, expectedNil, ptrNilSlice.Value().([]*int))
}

func TestSliceOfStructLiterals(t *testing.T) {
	ss := Slice(struct{ string }{"abc"}, struct{ string }{"def"})
	expected := []struct{ string }{struct{ string }{"abc"}, struct{ string }{"def"}}
	assert.Equal(t, expected, ss.Value().([]struct{ string }))

	ssInt := Slice(struct {
		string
		int
	}{"abc", 123}, struct {
		string
		int
	}{"def", 456})

	expectedInt := []struct {
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

	assert.Equal(t, expectedInt, ssInt.Value().([]struct {
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

func TestSortByFunc(t *testing.T) {
	animals := []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}
	byLen := func(a, b string) bool { return len(a) < len(b) }
	sorted := Slice(animals).SortBy(byLen)
	expected1 := []string{"dog", "cat", "cow", "pig", "bear", "bull", "iguana"}
	assert.Equal(t, expected1, animals)
	assert.Equal(t, sorted.Value().([]string), animals)

	lexically := func(a, b string) bool { return a < b }
	sorted.SortBy(lexically)
	expected2 := []string{"bear", "bull", "cat", "cow", "dog", "iguana", "pig"}
	assert.Equal(t, expected2, animals)
}

func TestMinByMaxBy(t *testing.T) {
	animals := []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	byLen := func(a, b string) bool { return len(a) < len(b) }
	lexically := func(a, b string) bool { return a < b }

	assert.Equal(t, "dog", Slice(animals).MinBy(byLen).(string))
	assert.Equal(t, "iguana", Slice(animals).MaxBy(byLen).(string))

	assert.Equal(t, "bear", Slice(animals).MinBy(lexically).(string))
	assert.Equal(t, "pig", Slice(animals).MaxBy(lexically).(string))
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

func TestSliceFilter(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	result := Slice(animals).
		Filter(func(s string) bool { return len(s) > 3 }).
		Value().([]string)

	expected := []string{"bear", "bull", "iguana"}
	assert.Equal(t, expected, result)
	assert.NotEqual(t, animals, result)

	result = Slice(animals).
		Filter(func(s string) bool { return len(s) < 3 }).
		Value().([]string)

	expected = []string{}
	assert.Equal(t, expected, result)
	assert.NotEqual(t, animals, result)
}

func TestSliceMap(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow"}

	result := Slice(animals).Map(strings.ToUpper).Value().([]string)

	expected := []string{"DOG", "CAT", "BEAR", "COW"}
	assert.Equal(t, expected, result)
	assert.NotEqual(t, animals, result)

	ints := Slice().Map(func(i interface{}) int { return 0 }).Value().([]int)
	assert.Equal(t, []int{}, ints)

	strings := Slice([]int{}).Map(func(i int) string { return "" }).Value().([]string)
	assert.Equal(t, []string{}, strings)

	assert.Panics(t, func() {
		t.Log(Slice().Map(func(s string) int { return 0 }).Value().([]int))
	})
}

func TestSliceReduce(t *testing.T) {
	var animals = []string{"dog", "cat", "bear", "cow"}

	joiner := func(a, b string) string { return a + "-" + b }
	silly := Slice(animals).Reduce(joiner, "monkey").(string)
	assert.Equal(t, "monkey-dog-cat-bear-cow", silly)

	counter := func(a int, b string) int { return a + len(b) }
	total := Slice(animals).Reduce(counter, 0).(int)
	assert.Equal(t, 13, total)

	total = Slice([]string{}).Reduce(counter, 0).(int)
	assert.Equal(t, 0, total)

	assert.Panics(t, func() {
		t.Log(Slice(animals).Reduce(counter, "0"))
	})
	assert.Panics(t, func() {
		badFn := func(a int, b string) string { return b }
		t.Log(Slice(animals).Reduce(badFn, 0))
	})
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

func TestSliceFlatten(t *testing.T) {
	five := []int{1, 2, 3, 4, 5}
	assert.Equal(t, five, Slice(1, 2, 3, 4, 5).Flatten().Value().([]int))
	assert.Equal(t, five, Slice(Slice(1, 2, 3, 4, 5)).Flatten().Value().([]int))
	assert.Equal(t, five, Slice(Slice(1, 2), Slice(3, 4, 5)).Flatten().Value().([]int))
	assert.Equal(t, five, Slice(Slice(Slice(1), Slice(2)), Slice(3, 4, 5)).Flatten().Value().([]int))
	assert.Equal(t, five, Slice(1, Slice(2), Slice(3, Slice(4, Slice(5)))).Flatten().Value().([]int))

	mixedTypes := Slice(1, Slice(2), Slice(3, Slice(4, Slice("5"))))
	assert.Equal(t, 5, len(mixedTypes.Flatten().Value().([]interface{})))

	nilValues := Slice(1, Slice(2), Slice(3, Slice(4, Slice(nil))))
	assert.Equal(t, 5, len(nilValues.Flatten().Value().([]int)))

	assert.Equal(t,
		[]int{1, 2, 3, 0},
		Slice([]interface{}{1, 2, 3, nil}).Flatten().Value().([]int))
	assert.Equal(t,
		[]interface{}{1, 2, "a", nil},
		Slice([]interface{}{1, 2, "a", nil}).Flatten().Value().([]interface{}))
	assert.Equal(t,
		[]interface{}{nil, nil, nil},
		Slice([]interface{}{nil, nil, nil}).Flatten().Value().([]interface{}))

	six := []int{1, 2, 3, 4, 5, 6}
	// [[[1 2] [3]] [[4 5 6]]]
	assert.Equal(t, six, Slice([][][]int{[][]int{[]int{1, 2}, []int{3}}, [][]int{[]int{4, 5, 6}}}).Flatten().Value().([]int))
	// [[[1 2] [3]] [[4 5 6]]]
	assert.Equal(t, six, Slice([][][]interface{}{[][]interface{}{[]interface{}{1, 2}, []interface{}{3}}, [][]interface{}{[]interface{}{4, 5, 6}}}).Flatten().Value().([]int))
	// [[1 2] 3 [4 5 [6]]]
	assert.Equal(t, six, Slice([]interface{}{[]interface{}{1, 2}, 3, []interface{}{4, 5, []interface{}{6}}}).Flatten().Value().([]int))
	// [[1 2] [4 5 [6]]]
	assert.Equal(t, six, Slice([][]interface{}{[]interface{}{1, 2, 3}, []interface{}{4, 5, []interface{}{6}}}).Flatten().Value().([]int))
	// [[1 2] 3 [4 5 [6]]]
	assert.Equal(t, six, Slice([]interface{}{[]int{1, 2}, 3, []interface{}{4, 5, []int{6}}}).Flatten().Value().([]int))
}

func TestSliceIsEmpty(t *testing.T) {
	assert.True(t, Slice().IsEmpty())
	assert.True(t, Slice([]int{}).IsEmpty())

	assert.False(t, Slice(1).IsEmpty())
	assert.False(t, Slice([]int{1}).IsEmpty())
}

func TestSlicePushPop(t *testing.T) {
	ints := Slice(1, 1)
	ints.Push(2)
	assert.Equal(t, 3, ints.Len())
	ints.Push(3)
	assert.Equal(t, []int{1, 1, 2, 3}, ints.Value().([]int))

	abc := Slice([]rune("abc"))
	assert.Equal(t, 3, abc.Len())
	assert.Equal(t, 'c', abc.Pop().(rune))
	assert.Equal(t, 'b', abc.Pop().(rune))
	assert.Equal(t, 'a', abc.Pop().(rune))
	assert.True(t, abc.IsEmpty())
	assert.True(t, abc.Pop() == nil)
}

func TestSliceShiftUnshift(t *testing.T) {
	ints := Slice(1, 1)
	ints.Unshift(2)
	assert.Equal(t, 3, ints.Len())
	ints.Unshift(3)
	assert.Equal(t, []int{3, 2, 1, 1}, ints.Value().([]int))

	empty := Slice()
	empty.Unshift("a")
	assert.Equal(t, 1, empty.Len())

	abc := Slice([]rune("abc"))
	assert.Equal(t, 3, abc.Len())
	assert.Equal(t, 'a', abc.Shift().(rune))
	assert.Equal(t, 'b', abc.Shift().(rune))
	assert.Equal(t, 'c', abc.Shift().(rune))
	assert.True(t, abc.IsEmpty())
	assert.True(t, abc.Shift() == nil)
}

func TestSliceInsert(t *testing.T) {
	runes := Slice('a', 'b')
	assert.Panics(t, func() { runes.Insert('z', -1) })
	assert.Panics(t, func() { runes.Insert('z', 3) })
	runes.Insert('c', 2)
	assert.Equal(t, []rune{'a', 'b', 'c'}, runes.Value().([]rune))
	runes.Insert('d', 0)
	assert.Equal(t, []rune{'d', 'a', 'b', 'c'}, runes.Value().([]rune))
	runes.Insert('e', 2)
	assert.Equal(t, []rune{'d', 'a', 'e', 'b', 'c'}, runes.Value().([]rune))

	empty := Slice()
	empty.Insert('a', 0)
	assert.Equal(t, 1, empty.Len())
}

func TestSliceDelete(t *testing.T) {
	runes := Slice('a', 'b')
	assert.Panics(t, func() { runes.Delete(-1) })
	assert.Panics(t, func() { runes.Delete(3) })
	runes.Delete(0)
	assert.Equal(t, []rune{'b'}, runes.Value().([]rune))
	runes.Insert('d', 1)
	runes.Insert('e', 1)
	runes.Insert('a', 1)
	runes.Delete(2)
	assert.Equal(t, []rune{'b', 'a', 'd'}, runes.Value().([]rune))
	runes.Delete(0)
	runes.Delete(0)
	runes.Delete(0)
	assert.Panics(t, func() { runes.Delete(0) })
}
