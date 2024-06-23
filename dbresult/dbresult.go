package dbresult

import (
        "database/sql/driver"
        "time"

        "github.com/jmoiron/sqlx"
        _ "github.com/mattn/go-sqlite3"
)

type GameResult struct {
        ID       int
        DateTime UnixTimestamp
	Size     int
        Score    int
        Moves    int
        Username string
}

func (gr *GameResult) Save(db *sqlx.DB) error {
        insertQuery := "INSERT INTO game_results (datetime, size, score, moves, username) VALUES (:datetime, :size, :score, :moves, :username)"
        tx := db.MustBegin()
        tx.NamedExec(insertQuery, gr)
        err := tx.Commit()
        if err != nil {
                return err
        }
        return nil
}

type UnixTimestamp time.Time

func (ut *UnixTimestamp) Scan(value interface{}) error {
        t := value.(int64)
        *ut = UnixTimestamp(time.Unix(t, 0))
        return nil
}

func (ut UnixTimestamp) Value() (driver.Value, error) {
        return time.Time(ut).Unix(), nil
}

func newGameResult(db *sqlx.DB, size int, score int, moves int, username string) error {
        result := &GameResult{
                DateTime: UnixTimestamp(time.Now()),
		Size:     size,
                Score:    score,
                Moves:    moves,
                Username: username,
        }
        err := result.Save(db)
        if err != nil {
                return err
        }
        return nil
}
