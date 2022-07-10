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

const layoutDate = "2006-01-02"

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

	loc, err := time.LoadLocation("Asia/Shanghai")
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

			start := randomTime(t, tt.startLayout == layoutDate, loc)
			end := randomTime(t, tt.endLayout == layoutDate, loc)

			rawJSON := bytes.Buffer{}
			rawJSON.Write([]byte{'{'})

			hasEnd := tt.endLayout != ""

			if hasEnd {
				rawJSON.WriteString(`"end":"`)
				rawJSON.WriteString(end.In(time.UTC).Format(tt.endLayout))
				rawJSON.WriteString(`",`)
			}

			rawJSON.WriteString(`"start":"`)
			rawJSON.WriteString(start.In(time.UTC).Format(tt.startLayout))
			rawJSON.WriteString(`","time_zone":"Asia/Shanghai"}`)

			date := notion.Date{}
			assert.NoError(t, json.Unmarshal(rawJSON.Bytes(), &date))

			assert.Equal(t, start.Format(tt.startLayout), date.Start.Format(tt.startLayout))
			assert.Equal(t, loc, date.Start.Location())

			if hasEnd {
				assert.Equal(t, end.Format(tt.endLayout), date.End.Format(tt.endLayout))
				assert.Equal(t, loc, date.End.Location())
			} else {
				assert.Nil(t, date.End)
			}

			start = start.In(loc)
			end = end.In(loc)

			res := bytes.Buffer{}

			res.WriteString(start.Format(tt.startLayout))

			if hasEnd {
				res.WriteString(" - ")
				res.WriteString(end.Format(tt.endLayout))
			}

			assert.Equal(t, res.String(), date.String(), "String should be equal")

			b, err := json.Marshal(date)
			assert.NoError(t, err)
			assert.Equal(t, rawJSON.String(), string(b))
		})
	}
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
