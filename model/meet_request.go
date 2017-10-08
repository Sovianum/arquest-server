package model

const (
	StatusPending  = "PENDING"
	StatusAccepted = "ACCEPTED"
	StatusDeclined = "DECLINED"
)

type MeetRequest struct {
	Id          int        `json:"id"`
	RequesterId int        `json:"requester_id"`
	RequestedId int        `json:"requested_id"`
	Time        QuotedTime `json:"time"`
	Status      string     `json:"status"`
}
