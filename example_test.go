package tricks_test

import (
	"fmt"
	"strings"

	"github.com/aviddiviner/tricks"
)

func ExampleSlice() {
	animals := []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}
	bearCow := tricks.Slice(animals).
		Map(strings.Title).
		Last(5).
		First(2).
		Value().([]string) // [Bear Cow]

	fmt.Println(bearCow)
	// Output: [Bear Cow]
}

func ExampleSlice_runes() {
	numbers := tricks.Slice(1, 2, 18, 1, 3, 1, 4, 1, 2, 18, 1)
	password := numbers.
		Map(func(i int) rune { return rune(i + 104) }).
		Reverse().
		Last(5).
		Value().([]rune)

	fmt.Println(string(password))
	// Output: kizji
}

func ExampleSlice_groupBy() {
	animals := []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}
	byLength := func(s string) int { return len(s) }

	schweinehund := tricks.Slice(animals).
		GroupBy(byLength). // map[4:[bear bull] 6:[iguana] 3:[dog cat cow pig]] ¹
		Only(3, 4).        // map[4:[bear bull] 3:[dog cat cow pig]] ¹
		Values().          // [[dog cat cow pig] [bear bull]] ¹
		Flatten().         // [dog cat cow pig bear bull]
		Sort().            // [bear bull cat cow dog pig]
		Last(2).           // [dog pig]
		Reverse().         // [pig dog]
		Join("-")          // "pig-dog"

	// ¹ No guarantee on ordering in a map.

	fmt.Println(schweinehund)
	// Output: pig-dog
}

func ExampleSlice_variadic() {
	numbers := tricks.Slice(1, 2, 18, 1, 3, 1, 4, 1, 2, 18, 1)
	magic := numbers.
		Reduce(func(s string, i int) string { return s + string(i+64) + "~" }, nil)

	fmt.Println(magic)
	// Output: A~B~R~A~C~A~D~A~B~R~A~
}
