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

var testTimezone = "Asia/Shanghai"

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

	shanghai, err := time.LoadLocation(testTimezone)
	require.NoError(t, err)

	for _, tt := range []struct {
		name        string
		startLayout string
		endLayout   string
	}{
		{"only start, as time", time.RFC3339, ""},
		{"start and end, as time", time.RFC3339, time.RFC3339},
		{"only start, as date", layoutDate, ""},
		{"start and end, as dates", layoutDate, layoutDate},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			start := randomTime(t, tt.startLayout == layoutDate, shanghai)
			end := randomTime(t, tt.endLayout == layoutDate, shanghai)

			rawJSON := bytes.Buffer{}
			rawJSON.WriteString(`{"end":`)

			hasEnd := tt.endLayout != ""

			if hasEnd {
				rawJSON.WriteString(`"`)
				rawJSON.WriteString(end.Format(tt.endLayout))
				rawJSON.WriteString(`"`)
			} else {
				rawJSON.WriteString("null")
			}

			rawJSON.WriteString(`,"start":"`)
			rawJSON.WriteString(start.Format(tt.startLayout))
			rawJSON.WriteString(`","time_zone":"Asia/Shanghai"}`)

			date := notion.Date{}
			assert.NoError(t, json.Unmarshal(rawJSON.Bytes(), &date))

			b, err := json.Marshal(date)
			assert.NoError(t, err)
			assert.Equal(t, rawJSON.String(), string(b))

			assertSameTime(t, date.Start, start)

			if tt.startLayout != layoutDate {
				assert.Equal(t, shanghai, date.Start.Location())
			}

			if hasEnd {
				assert.True(t, date.End.Equal(end))

				if tt.endLayout != layoutDate {
					assert.Equal(t, shanghai, date.End.Location())
				}
			} else {
				assert.Nil(t, date.End)
			}

			// String

			res := bytes.Buffer{}

			res.WriteString(start.Format(tt.startLayout))

			if hasEnd {
				date.End = &end

				res.WriteString(" - ")
				res.WriteString(end.Format(tt.endLayout))
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

	assert.EqualError(t, date.UnmarshalJSON([]byte{'{'}),
		"unexpected end of JSON input")

	assert.EqualError(t, date.UnmarshalJSON([]byte(`{"time_zone":"foo"}`)),
		"unknown time zone foo")

	assert.EqualError(t, date.UnmarshalJSON([]byte(`{"start":"foo"}`)),
		`parsing time "foo" as "2006-01-02T15:04:05Z07:00": cannot parse "foo" as "2006"`)

	assert.EqualError(t, date.UnmarshalJSON([]byte(`{"start":"2022-07-10T19:01:52Z","end":"foo"}`)),
		`parsing time "foo" as "2006-01-02T15:04:05Z07:00": cannot parse "foo" as "2006"`)
}
