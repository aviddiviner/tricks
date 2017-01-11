# tricks

An entirely unidiomatic approach to working with maps and slices in Go.

## What is this?

Originally started as an exercise in teaching myself Go reflection, I got a little carried away... and created something _awesome_. This goes out to all those Gophers with a yearning in their hearts for the good old days of chaining long strings of methods together (à la Ruby).

### Show me examples!

Sure. The best place to start is probably [the docs](//godoc.org/github.com/aviddiviner/tricks), but here's some sexy code to admire:

```go
animals := []string{"dog", "cat", "bear", "cow", "bull", "pig", "iguana"}
bearCow := tricks.Slice(animals).
    Map(strings.ToUpper).
    Last(5).
    First(2).
    Value().([]string) // [BEAR COW]
```

```go
// TODO: More, and better examples.
```

## Why did you do this?

**(The back-story.)**

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

## Should I actually use this?

Probably not.

If you're working on a big project, or if someone else has to maintain your code, certainly not. Having said that though... there's nothing _technically_ wrong with what I've done here. And you must admit, it _is_ pretty awesome. So, if you like writing less code, looking cool in front of your friends, and building things _super fast_ with your amazing new magic code skills, then you know what to do.

I use it for my pet projects, so I don't mind either way!

### I still can't decide

In all seriousness... to my mind, the goal is readability (or parse-ability / understandability). The trade-off here is between _expressiveness_ vs. _accuracy_.

On the side of _accuracy_, readability is gained from code that is clear and precise (as Go usually is). You can see exactly what is being done, and understand the inner workings of each piece. This makes for efficient code too.

On the side of _expressiveness_, readability is improved by simply reducing the code on the page, and keeping things short and to the point. This makes it easier to parse what is intended (vs. what is actually being done).

**tricks** makes it easier to write less code, and get your point across more succinctly. Yes, in all that brevity you may be obscuring over some important details. On balance though, I prefer the approach of less code, even if it means I'm copying an extra slice here or there.

### I still can't decide!

Which saying do you prefer? Pick one:

A. _“The devil's in the details.”_
B. _“Less is more.”_

If you chose option A, then move along. If you chose option B, then `go get` 'em.
