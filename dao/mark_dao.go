package dao

import (
	"database/sql"
	"github.com/Sovianum/arquest-server/model"
)

type MarkDAO interface {
	FinishQuest(userID, questID int) DBError
	MarkQuest(userID, questID int, mark float32) DBError
	GetUserMarks(userID int) ([]model.Mark, DBError)
}

type dbMarkDAO struct {
	db *sql.DB
}
