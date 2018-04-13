package model

type Mark struct {
	UserID  int     `json:"user_id"`
	QuestID int     `json:"quest_id"`
	Mark    float32 `json:"mark"`
}
