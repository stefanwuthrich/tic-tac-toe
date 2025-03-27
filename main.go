package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Constants for players and empty cell
const (
	playerX   = "X"
	playerO   = "O" // Computer will be 'O'
	emptyCell = " "
)

// Board type
type board [3][3]string

// --- Game Logic Functions ---

// initializeBoard creates an empty board
func initializeBoard() board {
	var b board
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			b[i][j] = emptyCell
		}
	}
	return b
}

// printBoard displays the current state of the board
func printBoard(b board) {
	fmt.Println("\n-------------")
	for i := 0; i < 3; i++ {
		fmt.Printf("| %s | %s | %s |\n", b[i][0], b[i][1], b[i][2])
		fmt.Println("-------------")
	}
	fmt.Println()
}

// getAvailableMoves returns a list of coordinates for empty cells
func getAvailableMoves(b board) [][2]int {
	var moves [][2]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] == emptyCell {
				moves = append(moves, [2]int{i, j})
			}
		}
	}
	return moves
}

// checkWin checks if a given player has won
func checkWin(b board, player string) bool {
	// Check rows
	for i := 0; i < 3; i++ {
		if b[i][0] == player && b[i][1] == player && b[i][2] == player {
			return true
		}
	}
	// Check columns
	for j := 0; j < 3; j++ {
		if b[0][j] == player && b[1][j] == player && b[2][j] == player {
			return true
		}
	}
	// Check diagonals
	if b[0][0] == player && b[1][1] == player && b[2][2] == player {
		return true
	}
	if b[0][2] == player && b[1][1] == player && b[2][0] == player {
		return true
	}
	return false
}

// isBoardFull checks if there are no empty cells left
func isBoardFull(b board) bool {
	return len(getAvailableMoves(b)) == 0
}

// evaluateBoard assigns a score for the Minimax algorithm
// +10 for Computer (O) win, -10 for Player (X) win, 0 for draw/ongoing
func evaluateBoard(b board) int {
	if checkWin(b, playerO) {
		return 10 // Computer wins
	}
	if checkWin(b, playerX) {
		return -10 // Player wins
	}
	return 0 // Draw or ongoing
}

// --- Minimax Algorithm ---

// minimax function implements the core logic
// isMaximizing determines if it's Computer's (O) turn (maximize score) or Player's (X) turn (minimize score)
func minimax(b board, depth int, isMaximizing bool) int {
	score := evaluateBoard(b)

	// Base cases: Terminal states (win/loss/draw)
	if score == 10 || score == -10 || isBoardFull(b) {
		return score
	}

	availableMoves := getAvailableMoves(b)

	if isMaximizing { // Computer's turn (O) - maximize the score
		bestScore := math.MinInt32
		for _, move := range availableMoves {
			r, c := move[0], move[1]
			b[r][c] = playerO // Make the move
			currentScore := minimax(b, depth+1, false) // Recurse for opponent's turn
			b[r][c] = emptyCell                         // Undo the move (backtrack)
			if currentScore > bestScore {
				bestScore = currentScore
			}
		}
		return bestScore
	} else { // Player's turn (X) - minimize the score
		bestScore := math.MaxInt32
		for _, move := range availableMoves {
			r, c := move[0], move[1]
			b[r][c] = playerX // Make the move
			currentScore := minimax(b, depth+1, true) // Recurse for computer's turn
			b[r][c] = emptyCell                        // Undo the move (backtrack)
			if currentScore < bestScore {
				bestScore = currentScore
			}
		}
		return bestScore
	}
}

// findBestMove determines the optimal move for the computer (Player O)
func findBestMove(b board) (int, int) {
	bestScore := math.MinInt32
	bestMove := [2]int{-1, -1}
	availableMoves := getAvailableMoves(b)

	// Handle the very first move randomly for variability, otherwise Computer always picks the same start
	if len(availableMoves) == 9 {
		rand.Seed(time.Now().UnixNano())
		move := availableMoves[rand.Intn(len(availableMoves))]
		return move[0], move[1]
	}


	for _, move := range availableMoves {
		r, c := move[0], move[1]
		b[r][c] = playerO // Try the move
		moveScore := minimax(b, 0, false) // Evaluate the move (next turn is minimizing player)
		b[r][c] = emptyCell              // Undo the move

		// Update best move if current move is better
		if moveScore > bestScore {
			bestScore = moveScore
			bestMove = move
		}
	}

	if bestMove[0] == -1 { // Should not happen in Tic Tac Toe if logic is correct, but as fallback
		fmt.Println("Error: No best move found, picking random available move.")
		if len(availableMoves) > 0 {
			rand.Seed(time.Now().UnixNano())
			move := availableMoves[rand.Intn(len(availableMoves))]
			return move[0], move[1]
		}
		return -1,-1 // Indicate error or no moves left if truly stuck
	}


	return bestMove[0], bestMove[1]
}

// --- Player Interaction ---

// getPlayerMove prompts the human player for their move and validates it
func getPlayerMove(b board) (int, int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your move (number 1-9): ")
		inputStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		inputStr = strings.TrimSpace(inputStr)
		moveNum, err := strconv.Atoi(inputStr)

		if err != nil || moveNum < 1 || moveNum > 9 {
			fmt.Println("Invalid input. Please enter a number between 1 and 9.")
			continue
		}

		// Convert 1-9 to 0-indexed row and column
		row := (moveNum - 1) / 3
		col := (moveNum - 1) % 3

		if b[row][col] != emptyCell {
			fmt.Println("Cell already taken. Choose an empty cell.")
			continue
		}

		return row, col
	}
}

// --- Main Game Loop ---

func main() {
	fmt.Println("Welcome to Go Tic Tac Toe!")
	fmt.Println("You are Player X, Computer is Player O.")
	fmt.Println("Enter a number (1-9) corresponding to the cell:")
	fmt.Println("-------------")
	fmt.Println("| 1 | 2 | 3 |")
	fmt.Println("-------------")
	fmt.Println("| 4 | 5 | 6 |")
	fmt.Println("-------------")
	fmt.Println("| 7 | 8 | 9 |")
	fmt.Println("-------------")


	gameBoard := initializeBoard()
	currentPlayer := playerX // Human starts

	for {
		printBoard(gameBoard)

		var row, col int
		if currentPlayer == playerX {
			// Player's turn
			row, col = getPlayerMove(gameBoard)
			gameBoard[row][col] = playerX
		} else {
			// Computer's turn
			fmt.Println("Computer's turn (O)...")
			row, col = findBestMove(gameBoard)
			if row == -1 { // Error case from findBestMove
				fmt.Println("Computer couldn't determine a move. Game ends unexpectedly.")
				break;
			}
			gameBoard[row][col] = playerO
			fmt.Printf("Computer chose cell %d\n", row*3+col+1)
		}

		// Check for game end conditions
		if checkWin(gameBoard, currentPlayer) {
			printBoard(gameBoard)
			if currentPlayer == playerX {
				fmt.Println("Congratulations! You (X) win!")
			} else {
				fmt.Println("Computer (O) wins!")
			}
			break // End game
		}

		if isBoardFull(gameBoard) {
			printBoard(gameBoard)
			fmt.Println("It's a draw!")
			break // End game
		}

		// Switch player
		if currentPlayer == playerX {
			currentPlayer = playerO
		} else {
			currentPlayer = playerX
		}
	}

	fmt.Println("Game Over.")
}
