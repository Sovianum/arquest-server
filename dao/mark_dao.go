package dao

import (
	"database/sql"
	"github.com/Sovianum/arquest-server/model"
)

const (
	getUserVotes = `
		SELECT user_id, quest_id, vote FROM QuestUserLink
	`
	markQuest = `
		UPDATE quest_user_link SET mark = $1, marked = TRUE WHERE user_id = $2 AND quest_id = $3
	`
	updateRating = `
		UPDATE quest SET rating = mark_count / (mark_count + 1) * rating + $1 / (mark_count + 1), $1 / (mark_count + 1) WHERE id = $2
	`
	finishQuest = `
		UPDATE quest_user_link SET finished = TRUE WHERE user_id = $1 AND quest_id = $2
	`
)

func NewMarkDAO(db *sql.DB) MarkDAO {
	return &dbMarkDAO{db: db}
}

type MarkDAO interface {
	FinishQuest(userID, questID int) error
	MarkQuest(userID, questID int, mark float32) error
	GetUserMarks(userID int) ([]model.Mark, error)
}

type dbMarkDAO struct {
	db *sql.DB
}

func (dao *dbMarkDAO) FinishQuest(userID, questID int) error {
	_, err := dao.db.Exec(finishQuest, userID, questID)
	return err
}

func (dao *dbMarkDAO) MarkQuest(userID, questID int, mark float32) error {
	if _, err := dao.db.Exec(markQuest, mark, userID, questID); err != nil {
		return err
	}
	_, err := dao.db.Exec(updateRating, mark, questID)
	return err
}

func (dao *dbMarkDAO) GetUserMarks(userID int) ([]model.Mark, error) {
	return dao.getMarks(getUserVotes, userID)
}

func (dao *dbMarkDAO) getMarks(sql string, args ...interface{}) ([]model.Mark, error) {
	var rows, err = dao.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Mark, 0)
	for rows.Next() {
		vote := model.Mark{}
		err = rows.Scan(
			&vote.UserID,
			&vote.QuestID,
			&vote.Mark,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, vote)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}
