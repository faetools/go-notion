package notion_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/faetools/go-notion/pkg/notion"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tzName = "Asia/Shanghai"

const (
	layoutDate = "2006-01-02"
	layoutTime = "2006-01-02T15:04:05.000Z07:00"
)

func randomTime(t *testing.T, isDate bool, loc *time.Location) time.Time {
	t.Helper()

	ts := time.Time{}
	fuzz.New().Fuzz(&ts)

	if isDate {
		loc = time.UTC
	}

	// make sure it is in the right location
	ts = ts.In(loc)

	// remove nanoseconds
	ts = ts.Add(-1 * time.Duration(ts.Nanosecond()))

	// remove hour, minute and second
	if isDate {
		ts = ts.Add(-1 * time.Duration(ts.Hour()) * time.Hour)
		ts = ts.Add(-1 * time.Duration(ts.Minute()) * time.Minute)
		ts = ts.Add(-1 * time.Duration(ts.Second()) * time.Second)
	}

	return ts
}

func TestDate(t *testing.T) {
	t.Parallel()

	shanghai, err := time.LoadLocation(tzName)
	require.NoError(t, err)

	for _, tt := range []struct {
		name   string
		layout string
		hasEnd bool
	}{
		{"only start, as time", time.RFC3339, false},
		{"start and end, as time", time.RFC3339, true},
		{"only start, as date", layoutDate, false},
		{"start and end, as dates", layoutDate, true},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			start := randomTime(t, tt.layout == layoutDate, shanghai)
			end := randomTime(t, tt.layout == layoutDate, shanghai)

			rawJSON := bytes.Buffer{}
			rawJSON.WriteString(`{"end":`)

			if tt.hasEnd {
				rawJSON.WriteString(`"`)
				rawJSON.WriteString(end.Format(tt.layout))
				rawJSON.WriteString(`"`)
			} else {
				rawJSON.WriteString("null")
			}

			rawJSON.WriteString(`,"start":"`)
			rawJSON.WriteString(start.Format(tt.layout))
			rawJSON.WriteString(`","time_zone":"Asia/Shanghai"}`)

			date := notion.Date{}
			assert.NoError(t, json.Unmarshal(rawJSON.Bytes(), &date))

			b, err := json.Marshal(date)
			assert.NoError(t, err)
			assert.Equal(t, rawJSON.String(), string(b))

			assertSameTime(t, date.Start, start)

			if tt.layout != layoutDate {
				assert.Equal(t, shanghai, date.Start.Location())
			}

			if tt.hasEnd {
				assert.True(t, date.End.Equal(end))

				if tt.layout != layoutDate {
					assert.Equal(t, shanghai, date.End.Location())
				}
			} else {
				assert.Nil(t, date.End)
			}

			// String

			res := bytes.Buffer{}

			res.WriteString(start.Format(tt.layout))

			if tt.hasEnd {
				date.End = &end

				res.WriteString(" - ")
				res.WriteString(end.Format(tt.layout))
			}

			assert.Equal(t, res.String(), date.String(), "String should be equal")
		})
	}
}

func assertSameTime(t *testing.T, a, b time.Time) {
	t.Helper()

	diff := a.Sub(b)
	if diff < 0 {
		diff = -1 * diff
	}

	// for some reason it gets off to up to one minute
	assert.Less(t, diff, time.Minute,
		"not the same time: %s vs. %s", a, b)
}

func TestDate_Errors(t *testing.T) {
	t.Parallel()

	date := notion.Date{}

	assert.EqualError(t, json.Unmarshal([]byte{'{'}, &date),
		"unexpected end of JSON input")

	assert.EqualError(t, json.Unmarshal([]byte(`{"time_zone":"foo"}`), &date),
		"unknown time zone foo")

	assert.EqualError(t, json.Unmarshal([]byte(`{"start":"foo"}`), &date),
		`parsing time "foo" as "2006-01-02T15:04:05Z07:00": cannot parse "foo" as "2006"`)

	assert.EqualError(t, json.Unmarshal([]byte(`{"start":"2022-07-10T19:01:52Z","end":"foo"}`), &date),
		`parsing time "foo" as "2006-01-02T15:04:05Z07:00": cannot parse "foo" as "2006"`)
}

func TestDate_UnusualTime(t *testing.T) {
	t.Parallel()

	z, err := time.LoadLocation(tzName)
	assert.NoError(t, err)

	ts := time.Date(1900, 1, 1, 0, 0, 3, 0, z)

	d := notion.Date{Start: ts, TimeZone: &tzName}

	b, err := json.Marshal(d)
	assert.NoError(t, err)

	assert.Equal(t, `{"end":null,"start":"1900-01-01T00:00:03+08:05","time_zone":"Asia/Shanghai"}`, string(b))

	newD := notion.Date{}
	assert.NoError(t, json.Unmarshal(b, &newD))

	assert.Equal(t, d.TimeZone, newD.TimeZone)
	assert.Equal(t, z, newD.Start.Location())

	assert.True(t, d.Start.Equal(newD.Start), "time is off by %s", newD.Start.Sub(d.Start))
}
