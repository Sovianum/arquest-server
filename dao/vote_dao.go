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
		UPDATE QuestUserLink SET mark = $1, marked = TRUE WHERE user_id = $2 AND quest_id = $3
	`
	finishQuest = `
		UPDATE QuestUserLink SET finished = TRUE WHERE user_id = $1
	`
)

func NewVoteDAO(db *sql.DB) VoteDAO {
	return &dbVoteDAO{db: db}
}

type VoteDAO interface {
	FinishQuest(userID, questID int) error
	MarkQuest(userID, questID int, mark float32) error
	GetUserVotes(userID int) ([]model.Vote, error)
}

type dbVoteDAO struct {
	db *sql.DB
}

func (dao *dbVoteDAO) FinishQuest(userID, questID int) error {
	_, err := dao.db.Exec(finishQuest, userID, questID)
	return err
}

func (dao *dbVoteDAO) MarkQuest(userID, questID int, mark float32) error {
	_, err := dao.db.Exec(markQuest, mark, userID, questID)
	return err
}

func (dao *dbVoteDAO) GetUserVotes(userID int) ([]model.Vote, error) {
	return dao.getVotes(getUserVotes, userID)
}

func (dao *dbVoteDAO) getVotes(sql string, args ...interface{}) ([]model.Vote, error) {
	var rows, err = dao.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Vote
	for rows.Next() {
		vote := model.Vote{}
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
