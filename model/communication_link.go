package model

type CommLink struct {
	Id   int        `json:"id"`
	Id1  int        `json:"id_1"`
	Id2  int        `json:"id_2"`
	Time QuotedTime `json:"time"`
}
