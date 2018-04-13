package dao

import (
	"database/sql"
	"github.com/Sovianum/arquest-server/model"
	"net/http"
)

const (
	getUserVotes = `
		SELECT user_id, quest_id, mark FROM quest_user_link WHERE user_id = $1
	`
	markQuest = `
		UPDATE quest_user_link SET mark = $1, marked = TRUE WHERE user_id = $2 AND quest_id = $3
	`
	updateRating = `
		UPDATE quest SET rating = mark_count / (mark_count + 1) * rating + $1 / (mark_count + 1), mark_count = $1 / (mark_count + 1) WHERE id = $2
	`
	finishQuest = `
		INSERT INTO quest_user_link (user_id, quest_id, completed) VALUES ($1, $2, TRUE) ON CONFLICT ON CONSTRAINT ux_user_id_quest_id DO UPDATE SET completed = TRUE 
	`
)

func NewMarkDAO(db *sql.DB) MarkDAO {
	return &dbMarkDAO{db: db}
}

type MarkDAO interface {
	FinishQuest(userID, questID int) DBError
	MarkQuest(userID, questID int, mark float32) DBError
	GetUserMarks(userID int) ([]model.Mark, DBError)
}

type dbMarkDAO struct {
	db *sql.DB
}

func (dao *dbMarkDAO) FinishQuest(userID, questID int) DBError {
	result, err := dao.db.Exec(finishQuest, userID, questID)
	if err != nil {
		return NewCrashDBErr(err)
	}
	return getResultErr(result)
}

func (dao *dbMarkDAO) MarkQuest(userID, questID int, mark float32) DBError {
	if r, err := dao.db.Exec(markQuest, mark, userID, questID); err != nil {
		return NewCrashDBErr(err)
	} else if rErr := getResultErr(r); rErr != nil {
		return rErr
	}

	if r, err := dao.db.Exec(updateRating, mark, questID); err != nil {
		return NewCrashDBErr(err)
	} else if rErr := getResultErr(r); rErr != nil {
		return rErr
	}
	return nil
}

func (dao *dbMarkDAO) GetUserMarks(userID int) ([]model.Mark, DBError) {
	marks, err := dao.getMarks(getUserVotes, userID)
	return marks, NewCrashDBErr(err)
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

func getResultErr(r sql.Result) DBError {
	if affected, err := r.RowsAffected(); err != nil {
		return NewCrashDBErr(err)
	} else if affected == 0 {
		return NewDBErr(http.StatusNotFound, "quest not found")
	}
	return nil
}
