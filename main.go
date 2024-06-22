package main

import (
 "fmt"
 "math/rand"
 "strings"

 "github.com/eiannone/keyboard"
)

type dir int

const (
 up dir = iota
 down
 left
 right
)

var (
 board  [4][4]int
 moves  int
 maxSum int
)

func init() {
 fill()
 printBoard()
}

func main() {
 if err := keyboard.Open(); err != nil {
  panic(err)
 }
 defer keyboard.Close()

 for {
  _, key, err := keyboard.GetKey()
  if err != nil {
   panic(err)
  }

  moves++
  var direction dir
  switch key {
  case 65517:
   direction = up
  case 65515:
   direction = left
  case 65516:
   direction = down
  case 65514:
   direction = right
  }

  if key == keyboard.KeyEsc {
   break
  }

  move(direction)
  printBoard()

  if !addTwo() {
   fmt.Println("\nGAME OVER...\nMaximal sum:", maxSum, "\nMoves:", moves)
   break
  }
 }
}

func printBoard() {
 fmt.Println()
 for _, row := range board {
  fmt.Println(strings.Replace(strings.Join(strings.Fields(fmt.Sprint(row)), "\t"), "'", "", -1))
 }
}

func fill() bool {
 fmt.Println("Press 'ESC' to quit the game. Good Luck!")
 line, col := rand.Intn(3), rand.Intn(3)
 firstTwo := false

 for !firstTwo {
  if board[line][col] == 0 {
   board[line][col] = 2
   firstTwo = true
  }

  if !checkZeroFindMaxSum() {
   return false
  }

  line, col = rand.Intn(3), rand.Intn(3)
 }

 return true
}

func addTwo() bool {
 if !checkZeroFindMaxSum() {
  return false
 }

 zeroLCList := getZeroLCList()
 random := rand.Intn(len(zeroLCList))
 board[zeroLCList[random][0]][zeroLCList[random][1]] = 2

 return true
}

func checkZeroFindMaxSum() bool {
 for _, line := range board {
  for _, elem := range line {
   if elem == 0 {
    return true
   }
   findMaxSum(elem)
  }
 }
 return false
}

func findMaxSum(elem int) {
 if elem > maxSum {
  maxSum = elem
 }
}

func getZeroLCList() [][]int {
 var zeroLCList [][]int

 for lineID, line := range board {
  for colID, elem := range line {
   if elem == 0 {
    zeroLCList = append(zeroLCList, []int{lineID, colID})
   }
  }
 }

 return zeroLCList
}

func move(direction dir) {
 switch direction {
 case up:
  reverseBoard(1)
  moveLeft()
  reverseBoard(3)
 case down:
  reverseBoard(3)
  moveLeft()
  reverseBoard(1)
 case left:
  moveLeft()
 case right:
  reverseBoard(2)
  moveLeft()
  reverseBoard(2)
 }
}

func reverseBoard(times int) {
 for time := 0; time < times; time++ {
  newBoard := [4][4]int{}
  iter := len(board) - 1

  for line := range newBoard {
   for col := range newBoard[line] {
    newBoard[line][col] = board[col][iter]
   }
   iter--
  }

  board = newBoard
 }
}

func moveLeft() {
 for col := range board {
  isActive := true

  for isActive {
   isActive = false

   for line := 0; line < len(board)-1; line++ {
    if board[col][line] == 0 && board[col][line+1] != 0 {
     board[col][line], board[col][line+1] = board[col][line+1], board[col][line]
     isActive = true
     continue
    } else if board[col][line] == board[col][line+1] && board[col][line] != 0 {
     board[col][line] *= 2
     board[col][line+1] = 0
     isActive = true
     continue
    }
   }
  }
 }
}
