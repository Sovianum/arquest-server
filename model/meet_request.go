package model

const (
	STATUS_PENDING  = "STATUS_PENDING"
	STATUS_ACCEPTED = "STATUS_ACCEPTED"
	STATUS_DECLINED = "STATUS_DECLINED"
)

type MeetRequest struct {
	Id          int        `json:"id"`
	RequesterId int        `json:"requester_id"`
	RequestedId int        `json:"requested_id"`
	Time        QuotedTime `json:"time"`
	Status      string
}
