# tricks [![GoDoc](https://godoc.org/github.com/aviddiviner/tricks?status.svg)](https://godoc.org/github.com/aviddiviner/tricks)

An entirely unidiomatic approach to working with maps and slices in Go.

## What is this?

Originally started as an exercise in teaching myself Go reflection, I got a little carried away and created something... interesting. A fast way of working with maps and slices by simply chaining methods together, à la Ruby.

### Show me examples!

Sure. The best place to start is probably [the docs](https://godoc.org/github.com/aviddiviner/tricks), but here's some pretty code to admire:

```go
animals := []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}
bearCow := tricks.Slice(animals).
    Map(strings.Title).
    Last(5).
    First(2).
    Value().([]string) // [Bear Cow]
```

```go
numbers := tricks.Slice(1, 2, 18, 1, 3, 1, 4, 1, 2, 18, 1) // now []int
magic := numbers.
    Reduce("~", func(s string, i int) string { return s + string(i+64) + "~" })

magic.(string) // "~A~B~R~A~C~A~D~A~B~R~A~"
```

```go
byLength := func(s string) int { return len(s) }
pigDog := tricks.Slice(animals).
    GroupBy(byLength).  // map[4:[bear bull] 6:[iguana] 3:[dog cat cow pig]] ¹
    Only(3, 4).         // map[4:[bear bull] 3:[dog cat cow pig]] ¹
    Values().           // [[dog cat cow pig] [bear bull]] ¹
    Flatten().          // [dog cat cow pig bear bull]
    Sort().             // [bear bull cat cow dog pig]
    Last(2).            // [dog pig]
    Reverse().          // [pig dog]
    Join("-")           // "pig-dog"

// ¹ No guarantee on ordering in a map.
```

### TL;DR Docs

_(Click these to expand.)_

<details>
<summary>slice.{All, Any, Many, One, None}</summary>

These take a `func(T) bool` and tell you whether the elements in the slice: all return true, any return true, more than one returns true, exactly one returns true, or none return true.

</details>
<details>
<summary>slice.{Map, Reduce}</summary>

The classics. Apply a `func(T) X` to every element of the slice and create a new slice `[]X` of the results. Reduce all the elements down to a single value by some `func(a, b T) T`.

</details>
<details>
<summary>slice.Filter</summary>

Choose only the elements for which some `func(T) bool` returns true.

</details>
<details>
<summary>slice.{Push, Pop, Shift, Unshift}</summary>

Append or remove an element from the start or end of the slice.

</details>
<details>
<summary>slice.{Insert, Delete}</summary>

Add or remove an element at any position in the slice.

</details>
<details>
<summary>slice.{First, Last}</summary>

Reslice to only take the first or last `n` elements.

</details>
<details>
<summary>slice.{Sort, Min, Max}</summary>

Sort the elements of the slice. Find the smallest or biggest values. As long as the slice is a normal type (`[]string`, `[]int`, etc.) or it implements `sort.Interface`, these all work.

</details>
<details>
<summary>slice.{SortBy, MinBy, MaxBy}</summary>

Sort, or find the smallest / biggest values by some `func(a, b T) bool` that returns whether element `a < b`.

</details>
<details>
<summary>slice.GroupBy</summary>

Apply a `func(V) K` to every element of the slice and group them into a map (`map[K][]V`) of the results.

</details>
<details>
<summary>slice.{Reverse, Flatten, Join}</summary>

Reverse the order of elements in the slice. Flatten a nested slice of slices into a one-dimensional slice. Join a slice of strings into a single string.

</details>
<details>
<summary>slice.{Copy, Value, Len, IsEmpty}</summary>

Copy the contents to a new underlying slice. Get the underlying slice value. Get the number of elements in the slice. Check if the slice is empty.

</details>

<details>
<summary>map.{Keys, Values}</summary>

Get a slice of only the key or values of the map.

</details>
<details>
<summary>map.Only</summary>

Get a map containing only the entries matching some list of keys.

</details>
<details>
<summary>map.HasKeys</summary>

Return true if all of the given keys are present in the map.

</details>
<details>
<summary>map.{Copy, Value, Len, IsEmpty}</summary>

Copy the contents to a new underlying map. Get the underlying map value. Get the number of elements in the map. Check if the map is empty.

</details>

## Why did you do this?

**(The back-story.)**

So, there I was one day, merrily coding in Go. Feeling productive and happy with my life, I examined a piece of code I had just written. I had some logs which I'd read off disk, and I wanted to group them by date, and only take the last few days (with maybe an offset to paginate them). So I looked at my code:

```go
func groupLogsByDate(logs []Timelog, amount, offset int) map[string][]Timelog {
    // Group the logs by date, into a map
    grouped := make(map[string][]Timelog)
    for _, log := range logs {
        day := log.Start.AsDate()
        grouped[day] = append(grouped[day], log)
    }

    // Get all the unique days, and sort them
    var days []string
    for day := range grouped {
        days = append(days, day)
    }
    sort.Strings(days)

    // Get only the days we want
    if amount+offset < len(days) {
        days = days[len(days)-(amount+offset):]
    }
    if amount < len(days) {
        days = days[:amount]
    }

    // Return a map of logs for the chosen days
    result := make(map[string][]Timelog)
    for _, day := range days {
        result[day] = grouped[day]
    }
    return result
}
```

I was pleased, and filled with a warm, fuzzy love for Go. So simple, and clear. _Finish en klaar._

But yet... I felt like I needed more. I longed for the excitement of mapping, reducing, filtering, sorting, grouping... all in one long line that does _all of the things_. I wanted a little box of tricks that I could use to just chain everything together and make **✨magic✨** happen.

I mean, all I really had to do was:

1. Group the logs into a map
2. Get the map's keys
3. Sort them
4. Take the last/first few
5. Return a map of only those keys

That's like 5 lines of Ruby, right?

So I went to bed that night, wrestling with my feelings of inner turmoil. I couldn't quiet that little inner voice. I knew it had to be possible. Go has function literals, right. Go has reflection. This must be doable. There must be a way to have my cake _and_ eat it.

Well... it turns out there was a way, and I found it. And you just found it too. I woke up early the next morning and, after much `reflect`-ing, I emerged with this thing of beauty:

```go
func groupLogsByDate(logs []Timelog, amount, offset int) map[string][]Timelog {
    grouped := tricks.Slice(logs).
        GroupBy(func(t Timelog) string { return t.Start.AsDate() })
    days := grouped.Keys().
        Sort().
        Last(amount + offset).
        First(amount)
    return grouped.
        Only(days).
        Value().(map[string][]Timelog)
}
```

And then promptly rewrote it like this:

```go
func groupLogsByDate(logs []Timelog, amount, offset int) map[string][]Timelog {
    gb := tricks.Slice(logs).GroupBy(func(t Timelog) string { return t.Start.AsDate() })
    return gb.Only(gb.Keys().Sort().Last(amount + offset).First(amount)).Value().(map[string][]Timelog)
}
```

**🤘YEAAA!🤘** Now that's what I'm talking about! I felt the mad rush of power from chaining all those methods and now I was _truly_ pleased. I slept well that night, knowing I had done a bad thing, but feeling damn good about it.

## So, should I actually use this?

Probably not. 😆

Type safety aside though, to my mind, it's a choice between a declarative vs. imperative style.

The declarative style is more _expressive_. We improve readability by simply reducing the code on the page, keeping things short and to the point. This makes it easier to parse what is intended (vs. what is actually being done).

The imperative style is more _accurate_. Readability is gained from code that is clear and precise (as Go usually is). You can see exactly what is being done, and understand the inner workings of each piece. This usually makes for more efficient code too.

I feel that **tricks** makes it easier to write less, and be more expressive, at the cost of reduced accuracy.

But yes, please only use this for writing tests, or in your pet projects. You don't want to take a dependency on a single package that changes how your code is structured in such a fundamental way. This forces everyone else to learn how some crazy package works just to maintain your code. Rather keep things plain and idiomatic.

Interestingly, there are some nice new features coming in Go 1.8 which do things similar to what I've done here, like [`sort.Slice`](https://tip.golang.org/pkg/sort/#Slice). So there is a balance to be struck between these two styles. Hopefully this package can inspire some people, and maybe more of these tricks will slowly be superseded by conveniences from the Go core.

[Now go and read the API docs please, and make up your mind over there.](https://godoc.org/github.com/aviddiviner/tricks)

## Wishlist

- `slice.Append(...interface{}) TrickSlice`
- `slice.Apply(func(T) T) TrickSlice` (like a `slice.Map` in place, same type)
- `slice.Compact() TrickSlice`
- `slice.Cut(i, j int)`
- `slice.DeepCopy() TrickSlice`
- `slice.Drop` / `DeleteIf` `(func(T) bool) TrickSlice`
- `slice.Expand(i, j int)`
- `slice.Select` / `Reject` `(func(T) bool) TrickSlice` (no reallocating?)
- `slice.Insert(n, ...interface{})` (insert any number of elements)
- `slice.Partition(func(T) bool) (a, b TrickSlice)`
- `slice.Product() float64`
- `slice.Sample(n int) TrickSlice`
- `slice.Shuffle() TrickSlice`
- `slice.Sum() float64`
- `slice.ToMap() TrickMap`
- `slice.Uniq() TrickSlice`
- `slice.Zip(...interface{}) TrickSlice`
- `map.DeepCopy() TrickMap`
- `map.Drop(func(K, V) bool) TrickMap`
- `map.Filter` / `Choose` / `Select` `(func(K, V) bool) TrickMap`
- `map.Merge(map[K]V)`
- Lazy evaluation / enumerators
- Combinatorics (choose, permute)
- https://github.com/golang/go/wiki/SliceTricks `Cut` / `Delete` / `Insert`
