package tricks_test

import (
	"fmt"

	"github.com/aviddiviner/tricks"
)

var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

func ExampleSlice() {
	bearCow := tricks.Slice(animals).Last(5).First(2).Value().([]string)

	fmt.Println(bearCow, animals)

	// Output: [bear cow] [dog cat bear cow bull pig iguana]
}

func ExampleSlice_groupby() {
	animalsByLength := tricks.Slice(animals).
		GroupBy(func(s string) int { return len(s) }).
		Value().(map[int][]string)

	fmt.Println(animalsByLength[3])

	// Output: [dog cat cow pig]
}
