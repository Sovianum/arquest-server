package model

import (
	"fmt"
	"time"
)

type QuotedTime time.Time

func (t *QuotedTime) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Format("2006-01-02T15:04:05")
	stamp := fmt.Sprintf("\"%v\"", ts)

	return []byte(stamp), nil
}

func (t *QuotedTime) UnmarshalJSON(b []byte) error {
	var layout = "2006-01-02T15:04:05"

	var inputS = string(b)
	var ts, err = time.Parse(layout, inputS[1:len(inputS)-1]) // slicing removes quotes

	if err != nil {
		return err
	}

	*t = QuotedTime(ts)
	return nil
}
