package model

import (
	"encoding/json"
)

const (
	PositionRequiredUserId = "\"user_id\" field required"
	PositionRequiredPoint  = "\"position\" field required"
	PositionRequireTime    = "\"time\" field required"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Position struct {
	Id     int        `json:"id"`
	UserId int        `json:"user_id"`
	Point  Point      `json:"point"`
	Time   QuotedTime `json:"time"`
}

func (pos *Position) UnmarshalJSON(data []byte) error {
	var err = checkPresence(
		data,
		[]string{"point"},
		[]string{PositionRequiredPoint},
	)
	if err != nil {
		return err
	}

	type positionAlias Position
	var dest = (*positionAlias)(pos)

	err = json.Unmarshal(data, dest)
	if err != nil {
		return err
	}

	err = pos.Validate()

	return err
}

func (pos *Position) Validate() error {
	return nil
}
