# tricks

An entirely unidiomatic approach to working with maps and slices in Go.

## What is this?

Originally started as an exercise in teaching myself Go reflection, I got a little carried away... and created something _awesome_. This goes out to all those Gophers with a yearning in their hearts for the good old days of chaining long strings of methods together (à la Ruby).

### Back-story

So, there I was one day, merrily coding in Go. Feeling so productive, and happy with my life, I examined a piece of code I had just written. I had some logs which I'd read off disk, and I wanted to group them by date, and only take the last few days (with maybe an offset to paginate them). So I looked at my code:

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

I was _pleased_, and filled with a warm, fuzzy love for Go. So simple, and clear. Done.

But yet... I didn't feel complete. I felt like I wanted more. I longed for the excitement of mapping, filtering, sorting, reducing, grouping... maybe all in one long line that does _all of the things_. I wanted a little box of tricks that I could use to just chain everything together and make **✨magic✨** happen.

I mean, all I really had to do was:

1. Group the logs into a map
1. Get the map's keys
1. Sort them
1. Take the last/first few
1. Return a map with those keys

That should be 5 lines of code, right? I mean, it used to be that way... in _Ruby_.

_“No!”_ I told myself. _“This is not Ruby! This is a grown-up language. Used by grown-ups. For big, serious, grown-up things!”_ ... _“Go is this way for a **reason**.”_

I went to bed that night, wrestling with my feelings of inner turmoil. I couldn't quiet that nagging inner voice. I knew it had to be possible. Go has function literals. Go has reflection. There must be a way to have my cake _and_ eat it.

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
    grouped := tricks.Slice(logs).GroupBy(func(t Timelog) string { return t.Start.AsDate() })
    return grouped.Only(grouped.Keys().Sort().Last(amount + offset).First(amount)).Value().(map[string][]Timelog)
}
```

**🤘YEAAA!🤘** Now that's what I'm talking about! I felt the mad rush of power from chaining all those methods and now, I was _truly pleased_. I slept well that night, knowing I had done a bad thing, but still, feeling damn good about it.
