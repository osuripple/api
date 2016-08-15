package common

import (
	"errors"
	"strconv"
	"time"
)

// OsuTimeFormat is the time format for scores in the DB. Can be used with time.Parse etc.
const OsuTimeFormat = "060102150405"

// OsuTime is simply a time.Time, but can be used to convert an
// osu timestamp in the database into a native time.Time.
type OsuTime time.Time

func (s *OsuTime) setTime(t string) error {
	newTime, err := time.Parse(OsuTimeFormat, t)
	if _, ok := err.(*time.ParseError); err != nil && !ok {
		return err
	}
	if err == nil {
		*s = OsuTime(newTime)
	}
	return nil
}

// Scan decodes src into an OsuTime.
func (s *OsuTime) Scan(src interface{}) error {
	if s == nil {
		return errors.New("rippleapi/common: OsuTime is nil")
	}
	switch src := src.(type) {
	case int64:
		return s.setTime(strconv.FormatInt(src, 64))
	case float64:
		return s.setTime(strconv.FormatInt(int64(src), 64))
	case string:
		return s.setTime(src)
	case []byte:
		return s.setTime(string(src))
	case nil:
		// Nothing, leave zero value on timestamp
	default:
		return errors.New("rippleapi/common: unhandleable type")
	}
	return nil
}

// MarshalJSON -> time.Time.MarshalJSON
func (s OsuTime) MarshalJSON() ([]byte, error) {
	return time.Time(s).MarshalJSON()
}

// UnmarshalJSON -> time.Time.UnmarshalJSON
func (s *OsuTime) UnmarshalJSON(x []byte) error {
	t := new(time.Time)
	err := t.UnmarshalJSON(x)
	*s = OsuTime(*t)
	return err
}
