package notion

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	layoutDate    = "2006-01-02"
	lenLayoutDate = len(layoutDate)
)

// tmpDate is used to unmarshall into this so we can properly adjust the times to be in the right time zone.
type tmpDate struct {
	End      *string `json:"end"`
	Start    string  `json:"start"`
	TimeZone *string `json:"time_zone"`
}

// UnmarshalJSON fulfils json.Unmarshaller.
func (d *Date) UnmarshalJSON(b []byte) error {
	tmp := &tmpDate{}

	err := json.Unmarshal(b, tmp)
	if err != nil {
		return err
	}

	// Time zone information for start and end. Possible values are extracted from the IANA database and they are based on the time zones from Moment.js.
	//
	// When time zone is provided, start and end should not have any UTC offset. In addition, when time zone is provided, start and end cannot be dates without time information.
	//
	// If null, time zone information will be contained in UTC offsets in start and end.
	d.TimeZone = tmp.TimeZone

	loc := time.UTC
	if tmp.TimeZone != nil {
		loc, err = time.LoadLocation(*d.TimeZone)
		if err != nil {
			return err
		}
	}

	d.Start, err = parseTimeOrDate(tmp.Start, loc)
	if err != nil {
		return err
	}

	if tmp.End != nil {
		end, err := parseTimeOrDate(*tmp.End, loc)
		if err != nil {
			return err
		}

		d.End = &end
	}

	return nil
}

func parseTimeOrDate(ts string, loc *time.Location) (time.Time, error) {
	if lenLayoutDate == len(ts) {
		return time.ParseInLocation(layoutDate, ts, loc)
	}

	return time.ParseInLocation(time.RFC3339, ts, loc)
}

func (d Date) String() string {
	if d.End == nil {
		return d.Start.String()
	}

	return fmt.Sprintf("%s-%s", d.Start, d.End)
}
