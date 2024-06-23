package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/eiannone/keyboard"
)

type Output interface {
	send(message string)
}

type ConsoleOutput struct{}

func (co *ConsoleOutput) send(message string) {
	fmt.Println(message)
}

type WebUIOutput struct{}

func (wo *WebUIOutput) send(message string) {
	// ...
}

type BoardGenerator interface {
	CreateBoard(size int) [][]int
}

type DynamicBoard struct{}

func (db *DynamicBoard) CreateBoard(size int) [][]int {
	lines := make([][]int, size)
	for line := range lines {
		lines[line] = make([]int, size)
	}
	return lines
}

type _2048 struct {
	board  [][]int
	size   int
	moves  int
	maxSum int
	output Output
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (g *_2048) Init() {
	db := &DynamicBoard{}
	g.size = 5
	g.board = db.CreateBoard(g.size)
	g.fill()
	g.printBoard()
}

func (g *_2048) fill() bool {
	g.output.send("Press 'ESC' to quit the game. Good Luck!")
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
	if elem > g.maxSum {
		g.maxSum = elem
	}
}

func (g *_2048) printBoard() {
	g.output.send("")
	for _, row := range g.board {
		g.output.send(strings.Replace(strings.Join(strings.Fields(fmt.Sprint(row)), "\t"), "'", "", -1))
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

func main() {
	game := _2048{output: &ConsoleOutput{}}
	game.Init()
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()
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
		}
		if key == keyboard.KeyEsc {
			break
		}
		game.move(direction)
		game.printBoard()
		if !game.addTwo() {
			game.output.send(fmt.Sprintf("\nGAME OVER...\nMaximal sum: %d\nMoves: %d", game.maxSum, game.moves))
			break
		}
	}
}
