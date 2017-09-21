package model

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

func (pos *Position) UnmarshalJSON([]byte) error {

}

func (pos *Position) Validate() error {
	return nil
}

