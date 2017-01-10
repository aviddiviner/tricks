package tricks_test

import (
	"fmt"
	"strings"

	"github.com/aviddiviner/tricks"
)

var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

func ExampleSlice() {
	bearCow := tricks.Slice(animals).
		Map(strings.ToUpper).
		Last(5).
		First(2).
		Value().([]string)

	fmt.Println(bearCow)

	// Output: [BEAR COW]
}

func ExampleSlice_groupby() {
	byLength := tricks.Slice(animals).
		Copy().Sort().
		GroupBy(func(s string) int { return len(s) }).
		Value().(map[int][]string)

	fmt.Println(byLength[3])

	// Output: [cat cow dog pig]
}
