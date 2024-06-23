package dbresult

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type GameResult struct {
	ID       int
	DateTime UnixTimestamp
	Score    int
	Moves    int
	Username string
}

type UnixTimestamp time.Time

func (ut *UnixTimestamp) Scan(value interface{}) error {
	t, ok := value.(int64)
	if !ok {
		return fmt.Errorf("could not convert value to int64")
	}
	*ut = UnixTimestamp(time.Unix(t, 0))
	return nil
}

func (ut UnixTimestamp) Value() (driver.Value, error) {
	return time.Time(ut).Unix(), nil
}

func NewGameResult(db *sqlx.DB, score int, moves int, username string) (*GameResult, error) {
	result := &GameResult{
		DateTime: UnixTimestamp(time.Now()),
		Score:    score,
		Moves:    moves,
		Username: username,
	}

	insertQuery := "INSERT INTO game_results (datetime, score, moves, username) VALUES (:datetime, :score, :moves, :username)"
	tx := db.MustBegin()
	tx.NamedExec(insertQuery, result)
	err := tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}
