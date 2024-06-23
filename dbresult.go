// package dbresult

// import (
//  "database/sql"
//  "fmt"
//  "time"

//  "github.com/jmoiron/sqlx"
//  _ "github.com/mattn/go-sqlite3"
// )

// type GameResult struct {
//   ID       int            db:"id"
//   DateTime UnixTimestamp  db:"datetime"
//   Score    int            db:"score"
//   Moves    int            db:"moves"
//   Username string         db:"username"
// }