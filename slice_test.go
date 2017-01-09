package tricks_test

import (
	"fmt"

	"github.com/aviddiviner/tricks"
)

var animals = []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

func ExampleSlice() {
	slice := tricks.Slice(animals).Last(5).First(2).Value().([]string)
	fmt.Println(slice)
	fmt.Println(animals)

	// Output:
	// [bear cow]
	// [dog cat bear cow bull pig iguana]
}
