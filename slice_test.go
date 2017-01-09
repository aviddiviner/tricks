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
	var numbers = []int{4, 3, 2, 1}

	sorted := Slice(numbers).Copy().Sort().Value().([]int)

	assert.Equal(t, []int{1, 2, 3, 4}, sorted)
	assert.Equal(t, []int{4, 3, 2, 1}, numbers)
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
