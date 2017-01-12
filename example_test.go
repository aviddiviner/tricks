package tricks_test

import (
	"fmt"
	"strings"

	"github.com/aviddiviner/tricks"
)

func ExampleSlice() {
	animals := []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	bearCow := tricks.Slice(animals).
		Map(strings.ToUpper).
		Last(5).
		First(2).
		Value().([]string)

	fmt.Println(bearCow)
	// Output: [BEAR COW]
}

func ExampleSlice_strings() {
	password := tricks.Slice([]rune("abracadabra")).Reverse().Value().([]rune)
	fmt.Println(string(password))
	// Output: arbadacarba
}

// TODO: Add variadic example.

func ExampleTrickSlice_GroupBy() {
	animals := []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}

	byLength := tricks.Slice(animals).
		Copy().
		Sort().
		Reverse().
		GroupBy(func(s string) int { return len(s) }).
		Value().(map[int][]string)

	fmt.Println(byLength[3])
	// Output: [pig dog cow cat]
}
