package model

type Vote struct {
	UserID  int     `json:"user_id"`
	QuestID int     `json:"quest_id"`
	Rating  float32 `json:"rating"`
}
