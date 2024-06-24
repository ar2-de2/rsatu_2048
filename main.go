package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ar2-de2/rsatu_2048/dbresult"

	"github.com/eiannone/keyboard"
	"github.com/jmoiron/sqlx"
)

const IOTIMEOUT = 6

type _2048 struct {
        board [][]int
        size  int
        moves int
        score int
}

type DynamicBoard struct{}

func (db *DynamicBoard) CreateBoard(size int) [][]int {
	lines := make([][]int, size)
	for line := range lines {
		lines[line] = make([]int, size)
	}
	return lines
}

type Direction int

const (
	None Direction = iota
	Up
	Down
	Left
	Right
)

func (g *_2048) init(size int, io IO) {
	db := &DynamicBoard{}
	g.size = size
	g.board = db.CreateBoard(g.size)
	g.fill()
	g.printBoard(io)
}

func (g *_2048) fill() bool {
	line, col := rand.Intn(g.size-1), rand.Intn(g.size-1)
	firstTwo := false
	for !firstTwo {
		if g.board[line][col] == 0 {
			g.board[line][col] = 2
			firstTwo = true
		}
		if !g.checkZeroFindMaxSum() {
			return false
		}
		line, col = rand.Intn(g.size-1), rand.Intn(g.size-1)
	}
	return true
}

func (g *_2048) checkZeroFindMaxSum() bool {
	for _, line := range g.board {
		for _, elem := range line {
			if elem == 0 {
				return true
			}
			g.findMaxSum(elem)
		}
	}
	return false
}

func (g *_2048) findMaxSum(elem int) {
	if elem > g.score {
		g.score = elem
	}
}

func (g *_2048) printBoard(io IO) {
	io.send("", IOTIMEOUT)
	for _, row := range g.board {
		io.send(strings.Replace(strings.Join(strings.Fields(fmt.Sprint(row)), "\t"), "'", "", -1), IOTIMEOUT)
	}
}

func (g *_2048) move(direction Direction) {
	switch direction {
	case Up:
		g.reverseBoard(1)
		g.moveLeft()
		g.reverseBoard(3)
	case Down:
		g.reverseBoard(3)
		g.moveLeft()
		g.reverseBoard(1)
	case Left:
		g.moveLeft()
	case Right:
		g.reverseBoard(2)
		g.moveLeft()
		g.reverseBoard(2)
	}
}

func (g *_2048) reverseBoard(times int) {
	for time := 0; time < times; time++ {
		db := &DynamicBoard{}
		newBoard := db.CreateBoard(g.size)
		iter := len(g.board) - 1
		for line := range newBoard {
			for col := range newBoard[line] {
				newBoard[line][col] = g.board[col][iter]
			}
			iter--
		}
		g.board = newBoard
	}
}

func (g *_2048) moveLeft() {
	for line := range g.board {
		isActive := true
		for isActive {
			isActive = false
			for col := 0; col < len(g.board)-1; col++ {
				if g.board[line][col] == 0 && g.board[line][col+1] != 0 {
					g.board[line][col], g.board[line][col+1] = g.board[line][col+1], g.board[line][col]
					isActive = true
					continue
				} else if g.board[line][col] == g.board[line][col+1] && g.board[line][col] != 0 {
					g.board[line][col] *= 2
					g.board[line][col+1] = 0
					isActive = true
					continue
				}
			}
		}
	}
}

func (g *_2048) addTwo() bool {
	if !g.checkZeroFindMaxSum() {
		return false
	}
	zeroLCList := g.getZeroLCList()
	random := rand.Intn(len(zeroLCList))
	g.board[zeroLCList[random][0]][zeroLCList[random][1]] = 2
	return true
}

func (g *_2048) getZeroLCList() [][]int {
	var zeroLCList [][]int
	for lineID, line := range g.board {
		for colID, elem := range line {
			if elem == 0 {
				zeroLCList = append(zeroLCList, []int{lineID, colID})
			}
		}
	}
	return zeroLCList
}

type IO interface {
	send(message string, secTimeout int) error
	receive(secTimeout int) (string, error)
	directionDetector(db *sqlx.DB, game *_2048, username string) bool
}

type ioConsole struct{}

func (ioc *ioConsole) send(message string, secTimeout int) error {
	fmt.Println(message)
	return nil
}

func (ioc *ioConsole) receive(secTimeout int) (string, error) {
	for {
		reader := bufio.NewReader(os.Stdin)
		ch := make(chan string, 1)
		go func() {
			text, _ := reader.ReadString('\n')
			ch <- strings.TrimSpace(text)
		}()
		select {
		case <-time.After(time.Duration(secTimeout) * time.Hour):
			return "", errors.New("input timed out")
		case text := <-ch:
			if text == "exit" {
				os.Exit(1)
			}
			return text, nil
		}
	}
}

func (ioc *ioConsole) directionDetector(db *sqlx.DB, game *_2048, username string) bool {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		os.Exit(1)
	}()
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		var direction Direction
		switch key {
		case 65517:
			direction = Up
		case 65515:
			direction = Left
		case 65516:
			direction = Down
		case 65514:
			direction = Right
		case keyboard.KeyEsc:
			return false
		case keyboard.KeyCtrlC, keyboard.KeyCtrlD:
			os.Exit(1)
		default:
			direction = None
		}
		if direction == None {
			continue
		}
		if !gamePostRunner(game, db, ioc, direction, username) {
			break
		}
	}
	return true
}

type ioWebUI struct{}

func (iow *ioWebUI) send(message string) {
	// ...
}

func gamePreRunner(db *sqlx.DB, io IO) bool {
	io.send("\nPlease enter your username (min. 1 and max. 4 symbols): ", IOTIMEOUT)
	username, err := io.receive(IOTIMEOUT)
	if len(username) < 1 || len(username) > 4 || err != nil {
		io.send("Invalid name! Please try again...", IOTIMEOUT)
		return false
	}
	io.send("Please enter the size of the board ('x*x', where 'x' is between 3 and 5): ", IOTIMEOUT)
	stringSize, err := io.receive(IOTIMEOUT)
	size, err := strconv.Atoi(stringSize)
	if err != nil || size < 3 || size > 5 {
		io.send("Invalid size! Please try again...", IOTIMEOUT)
		return false
	}
	io.send("Press 'ESC' to start a new game. Use the arrows to make your moves. Good Luck!", IOTIMEOUT)
	game := _2048{}
	game.init(size, io)
	return io.directionDetector(db, &game, username)
}

func gamePostRunner(game *_2048, db *sqlx.DB, io IO, direction Direction, username string) bool {
	game.move(direction)
	game.moves++
	game.printBoard(io)
	if !game.addTwo() {
		io.send(fmt.Sprintf("\nGAME OVER...\nScore (maximal sum): %d\nMoves: %d", game.score, game.moves), IOTIMEOUT)
		if err := dbresult.NewGameResult(db, game.size, game.score, game.moves, username); err != nil {
			panic(err)
		}
		results, err := dbresult.Top3Results(db, game.size)
		if err != nil {
			panic(err)
		}
		io.send(fmt.Sprintf("\nTop 3 Results for '%s*%s' game:\n%s", strconv.Itoa(game.size), strconv.Itoa(game.size), results), IOTIMEOUT)
		return false
	}
	return true
}

func main() {
	db, err := sqlx.Open("sqlite3", "sqlite/game.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS game_results (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                datetime INTEGER,
		size INTEGER,
                score INTEGER,
                moves INTEGER,
                username TEXT
        )`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}
	for {
		io := &ioConsole{}
		gamePreRunner(db, io)
	}
}
