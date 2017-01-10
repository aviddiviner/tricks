# tricks

An entirely unidiomatic approach to working with maps and slices in Go.

## Back-story

Originally created as an exercise in teaching myself Go reflection, I got a little carried away... and created something _beautiful_. This goes out to all those Gophers with a yearning in their hearts for the good old days of chaining long lists of methods together (à la Ruby).

### Compare

There I was one day, merrily coding in Go. Feeling so productive, and happy with my life, I reflected on a piece of code I had just written. I had some logs which I'd read off disk, and I wanted to group them by date, and only take the last few days (with maybe an offset to paginate them). So I looked at my code:

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

I was pleased, and filled with a warm, fuzzy love for Go. So simple, and easy. Done.

But yet... I didn't feel complete. I felt like I wanted more. I wanted that rush of mapping, filtering, sorting, reducing, grouping... all in one long line that _does it all_. I wanted a little box of tricks that I could use to just chain everything together and make ✨_magic_✨ happen.

I mean, all I was really doing was:

1. Group the logs into a map
1. Get the map's keys
1. Sort them
1. Take the last/first few
1. Return a map with those keys

That should be 5 lines of code, right? I mean, it used to be that way... in _Ruby_.

But _“No!”_ I told myself. _“This is not Ruby! This is a grown-up language. Used by adults. For big, important things!”_ I went to bed that night, wrestling with feelings of inner turmoil. I couldn't quiet that nagging inner voice. I knew it had to be possible. There must be a way.

Well, there was a way, and I found it. I got up early that next morning, and after much `reflect`-ing, I emerged with this:

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
    grouped := tricks.Slice(logs).GroupBy(func(t Timelog) string { return t.Start.AsDate() })
    return grouped.Only(grouped.Keys().Sort().Last(amount + offset).First(amount)).Value().(map[string][]Timelog)
}
```

YEAAA! Now that's what I'm talking about! I felt the mad rush of power from chaining all those methods and I was _pleased_. I slept well that night, knowing I had done a bad thing, but still, feeling so damn good about it.
