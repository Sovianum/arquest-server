package sqldao

import (
	"database/sql"
	"github.com/Sovianum/arquest-server/model"
	"github.com/Sovianum/arquest-server/dao"
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

func NewQuestDAO(db *sql.DB) dao.QuestDAO {
	return &dbQuestDAO{db: db}
}

type dbQuestDAO struct {
	db *sql.DB
}

func (sqldao *dbQuestDAO) GetFinishedQuests(userID int) ([]model.Quest, dao.DBError) {
	return sqldao.getQuests(getFinishedQuest, userID)
}

func (sqldao *dbQuestDAO) GetAllQuests() ([]model.Quest, dao.DBError) {
	return sqldao.getQuests(getAllQuests)
}

func (sqldao *dbQuestDAO) ExistsByID(questID int) (bool, dao.DBError) {
	row := sqldao.db.QueryRow(existQuest, questID)
	cnt := 0
	err := row.Scan(&cnt)
	if err != nil {
		return false, dao.NewCrashDBErr(err)
	}
	return cnt > 0, nil
}

func (sqldao *dbQuestDAO) getQuests(sql string, args ...interface{}) ([]model.Quest, dao.DBError) {
	rows, err := sqldao.db.Query(sql, args...)
	if err != nil {
		return nil, dao.NewCrashDBErr(err)
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
			return nil, dao.NewCrashDBErr(err)
		}
		result = append(result, quest)
	}

	err = rows.Err()
	if err != nil {
		return nil, dao.NewCrashDBErr(err)
	}
	return result, nil
}
