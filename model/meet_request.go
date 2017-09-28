package model

type CommRequest struct {
	Id          int        `json:"id"`
	RequesterId int        `json:"requester_id"`
	RequestedId int        `json:"requested_id"`
	Time        QuotedTime `json:"time"`
}
