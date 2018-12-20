package dao

import (
	"github.com/Sovianum/arquest-server/model"
)

type QuestDAO interface {
	GetFinishedQuests(userID int) ([]model.Quest, DBError)
	GetAllQuests() ([]model.Quest, DBError)
	ExistsByID(questID int) (bool, DBError)
}
