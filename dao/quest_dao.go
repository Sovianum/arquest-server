package dao

import (
	"database/sql"
	"github.com/Sovianum/arquest-server/model"
)

const (
	getAllQuests     = `SELECT id, name, description, rating FROM quest`
	getFinishedQuest = `
		SELECT q.id AS id, q.name AS name, q.description AS description, q.rating AS rating
		FROM 
			quest AS q 
			JOIN quest_user_link AS link ON q.id = link.quest_id 
		WHERE link.user_id = $1 AND link.completed
	`
	existQuest = `SELECT count(*) FROM quest WHERE id = $1`
)

func NewQuestDAO(db *sql.DB) QuestDAO {
	return &dbQuestDAO{db: db}
}

type QuestDAO interface {
	GetFinishedQuests(userID int) ([]model.Quest, DBError)
	GetAllQuests() ([]model.Quest, DBError)
	ExistsByID(questID int) (bool, DBError)
}

type dbQuestDAO struct {
	db *sql.DB
}

func (dao *dbQuestDAO) GetFinishedQuests(userID int) ([]model.Quest, DBError) {
	return dao.getQuests(getFinishedQuest, userID)
}

func (dao *dbQuestDAO) GetAllQuests() ([]model.Quest, DBError) {
	return dao.getQuests(getAllQuests)
}

func (dao *dbQuestDAO) ExistsByID(questID int) (bool, DBError) {
	row := dao.db.QueryRow(existQuest, questID)
	cnt := 0
	err := row.Scan(&cnt)
	if err != nil {
		return false, NewCrashDBErr(err)
	}
	return cnt > 0, nil
}

func (dao *dbQuestDAO) getQuests(sql string, args ...interface{}) ([]model.Quest, DBError) {
	rows, err := dao.db.Query(sql, args...)
	if err != nil {
		return nil, NewCrashDBErr(err)
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
			return nil, NewCrashDBErr(err)
		}
		result = append(result, quest)
	}

	err = rows.Err()
	if err != nil {
		return nil, NewCrashDBErr(err)
	}
	return result, nil
}
