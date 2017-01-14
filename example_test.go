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

func ExampleSlice_strings() {
	password := tricks.Slice([]rune("abracadabra")).Reverse().Value().([]rune)
	fmt.Println(string(password))
	// Output: arbadacarba
}

func ExampleSlice_variadic() {
	runes := tricks.Slice(1, 2, 18, 1, 3, 1, 4, 1, 2, 18, 1).
		Map(func(i int) rune { return rune(i + 96) }).
		Value().([]rune)

	fmt.Println(string(runes))
	// Output: abracadabra
}
