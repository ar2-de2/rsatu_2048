package main

import (
        "bufio"
        "fmt"
        "math/rand"
        "os"
        "os/signal"
        "strconv"
        "strings"
        "syscall"

        "github.com/ar2-de2/rsatu_2048/dbresult"

        "github.com/eiannone/keyboard"
        "github.com/jmoiron/sqlx"
)

type DynamicBoard struct{}

func (db *DynamicBoard) CreateBoard(size int) [][]int {
        lines := make([][]int, size)
        for line := range lines {
                lines[line] = make([]int, size)
        }
        return lines
}

type _2048 struct {
        board [][]int
        size  int
        moves int
        score int
        io    IO
}

type Direction int

const (
        None Direction = iota
        Up
        Down
        Left
        Right
)

func (g *_2048) Init(size int) {
        db := &DynamicBoard{}
        g.size = size
        g.board = db.CreateBoard(g.size)
        g.fill()
        g.printBoard()
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

func (g *_2048) printBoard() {
        g.io.send("")
        for _, row := range g.board {
                g.io.send(strings.Replace(strings.Join(strings.Fields(fmt.Sprint(row)), "\t"), "'", "", -1))
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
        for col := range g.board {
                isActive := true
                for isActive {
                        isActive = false
                        for line := 0; line < len(g.board)-1; line++ {
                                if g.board[col][line] == 0 && g.board[col][line+1] != 0 {
                                        g.board[col][line], g.board[col][line+1] = g.board[col][line+1], g.board[col][line]
                                        isActive = true
                                        continue
                                } else if g.board[col][line] == g.board[col][line+1] && g.board[col][line] != 0 {
                                        g.board[col][line] *= 2
                                        g.board[col][line+1] = 0
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
        send(message string)
        receive() string
        runGame(db *sqlx.DB)
}

type ioConsole struct{}

func (ioc *ioConsole) send(message string) {
        fmt.Println(message)
}

func (ioc *ioConsole) receive() string {
        reader := bufio.NewReader(os.Stdin)
        text, _ := reader.ReadString('\n')
        text = strings.TrimSpace(text)
        if text == "exit" {
                os.Exit(1)
        }
        return text
}

func (ioc *ioConsole) runGame(db *sqlx.DB) {
        game := _2048{}
        game.io.send("\nPlease enter your username: ")
        username := game.io.receive()
        game.io.send("Please enter the size of the board ('x'*'x', where 'x' is between 4 and 6): ")
        size, err := strconv.Atoi(game.io.receive())
        if err != nil || size < 4 || size > 6 {
                game.io.send("Invalid size. Please try again...")
                return
        }
        game.io.send("Press 'ESC' to start a new game. Use the arrows to make your moves. Good Luck!")
        game.Init(size)
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
                game.moves++
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
                        return
                case keyboard.KeyCtrlC, keyboard.KeyCtrlD:
                        os.Exit(1)
                default:
                        direction = None
                }
                if direction != None {
                        game.move(direction)
                        game.printBoard()
                }
                if !game.addTwo() {
                        game.io.send(fmt.Sprintf("\nGAME OVER...\nScore (maximal sum): %d\nMoves: %d", game.score, game.moves))
                        if err := dbresult.newGameResult(db, game.score, game.moves, username); err != nil {
                                panic(err)
                        }
                        break
                }
        }
}

type ioWebUI struct{}

func (iow *ioWebUI) send(message string) {
        // ...
}

func gameRunner(db *sqlx.DB, io IO) {
        game := &ioConsole{}
        game.runGame(db)
}

func main() {
        db, err := sqlx.Open("sqlite3", "game.db")
        if err != nil {
                fmt.Println("Error opening database:", err)
                return
        }
        defer db.Close()
        _, err = db.Exec(`CREATE TABLE IF NOT EXISTS game_results (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                datetime INTEGER,
                score INTEGER,
                moves INTEGER,
                username TEXT
        )`)
        if err != nil {
                fmt.Println("Error creating table:", err)
                return
        }
        for {
                game := &ioConsole{}
                game.runGame(db)
        }
}
