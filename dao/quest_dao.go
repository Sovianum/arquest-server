package dao

import (
	"database/sql"
	"github.com/Sovianum/arquest-server/model"
)

const (
	getAllQuests     = `SELECT id, name, description, rating FROM quest`
	getFinishedQuest = `
		SELECT q.id id, q.name name, q.description description, q.rating rating
		FROM 
			quest q 
			JOIN quest_user_link link ON q.id = link.quest_id 
		WHERE link.user_id = $1 link.finished
	`
)

func NewQuestDAO(db *sql.DB) QuestDAO {
	return &dbQuestDAO{db: db}
}

type QuestDAO interface {
	GetFinishedQuests(userID int) ([]model.Quest, error)
	GetAllQuests() ([]model.Quest, error)
}

type dbQuestDAO struct {
	db *sql.DB
}

func (dao *dbQuestDAO) GetFinishedQuests(userID int) ([]model.Quest, error) {
	return dao.getQuests(getFinishedQuest, userID)
}

func (dao *dbQuestDAO) GetAllQuests() ([]model.Quest, error) {
	return dao.getQuests(getAllQuests)
}

func (dao *dbQuestDAO) getQuests(sql string, args ...interface{}) ([]model.Quest, error) {
	var rows, err = dao.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Quest, 0)
	for rows.Next() {
		quest := model.Quest{}
		err = rows.Scan(
			&quest.ID,
			&quest.Name,
			&quest.Description,
			&quest.Rating,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, quest)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}
