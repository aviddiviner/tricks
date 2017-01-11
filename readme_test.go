package tricks_test

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/aviddiviner/tricks"
	"github.com/stretchr/testify/assert"
)

type Timelog struct {
	Activity string
	Start    TimelogTS
}

type TimelogTS struct{ time.Time }

func (t TimelogTS) AsDate() string {
	return fmt.Sprintf(`%04d-%02d-%02d`, t.Year(), t.Month(), t.Day())
}

func groupLogsOld(logs []Timelog, amount, offset int) map[string][]Timelog {
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

func groupLogsNew(logs []Timelog, amount, offset int) map[string][]Timelog {
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

// Check that I'm not talking crap in the README, and my examples there actually work.
func TestReadmeBackstoryExample(t *testing.T) {
	logs := []Timelog{
		Timelog{"Eating", TimelogTS{time.Date(2016, 10, 5, 9, 5, 9, 0, time.UTC)}},
		Timelog{"Coding", TimelogTS{time.Date(2016, 10, 5, 9, 29, 10, 0, time.UTC)}},
		Timelog{"Resting", TimelogTS{time.Date(2016, 10, 5, 9, 54, 36, 0, time.UTC)}},
		Timelog{"Meeting", TimelogTS{time.Date(2016, 10, 6, 10, 0, 22, 0, time.UTC)}},
		Timelog{"Coding", TimelogTS{time.Date(2016, 10, 6, 10, 6, 42, 0, time.UTC)}},
		Timelog{"Eating", TimelogTS{time.Date(2016, 10, 7, 8, 38, 24, 0, time.UTC)}},
		Timelog{"Coding", TimelogTS{time.Date(2016, 10, 7, 8, 48, 28, 0, time.UTC)}},
		Timelog{"Meeting", TimelogTS{time.Date(2016, 10, 7, 10, 0, 0, 0, time.UTC)}},
		Timelog{"Coding", TimelogTS{time.Date(2016, 10, 8, 10, 12, 13, 0, time.UTC)}},
		Timelog{"Surfing", TimelogTS{time.Date(2016, 10, 9, 11, 12, 26, 0, time.UTC)}},
		Timelog{"Eating", TimelogTS{time.Date(2016, 10, 10, 9, 5, 44, 0, time.UTC)}},
		Timelog{"Surfing", TimelogTS{time.Date(2016, 10, 10, 9, 20, 52, 0, time.UTC)}},
		Timelog{"Talking", TimelogTS{time.Date(2016, 10, 10, 9, 39, 36, 0, time.UTC)}},
		Timelog{"Meeting", TimelogTS{time.Date(2016, 10, 10, 10, 0, 46, 0, time.UTC)}},
		Timelog{"Coding", TimelogTS{time.Date(2016, 10, 11, 10, 9, 18, 0, time.UTC)}},
		Timelog{"Resting", TimelogTS{time.Date(2016, 10, 11, 10, 52, 35, 0, time.UTC)}},
		Timelog{"Coding", TimelogTS{time.Date(2016, 10, 11, 11, 0, 42, 0, time.UTC)}},
	}
	assert.EqualValues(t, groupLogsOld(logs, 0, 0), groupLogsNew(logs, 0, 0))
	assert.EqualValues(t, groupLogsOld(logs, 0, 5), groupLogsNew(logs, 0, 5))
	assert.EqualValues(t, groupLogsOld(logs, 5, 5), groupLogsNew(logs, 5, 5))
	assert.EqualValues(t, groupLogsOld(logs, 10, 5), groupLogsNew(logs, 10, 5))
	assert.EqualValues(t, groupLogsOld(logs, 50, 50), groupLogsNew(logs, 50, 50))
}
