package main

import (
	"html/template"
	"net/http"
	"strconv"
)

const (
	empty   = "-"
	player1 = "X"
	player2 = "O"
)

var (
	board = [3][3]string{
		{empty, empty, empty},
		{empty, empty, empty},
		{empty, empty, empty},
	}
	currentPlayer = player1
	tmpl          = template.Must(template.ParseFiles("game.html"))
)

func checkWin() string {
	// check rows
	for _, row := range board {
		if row[0] != empty && row[0] == row[1] && row[1] == row[2] {
			return row[0]
		}
	}

	// check columns
	for i := 0; i < 3; i++ {
		if board[0][i] != empty && board[0][i] == board[1][i] && board[1][i] == board[2][i] {
			return board[0][i]
		}
	}

	// check diagonals
	if board[0][0] != empty && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return board[0][0]
	}
	if board[0][2] != empty && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return board[0][2]
	}

	// check draw
	for _, row := range board {
		for _, cell := range row {
			if cell == empty {
				return ""
			}
		}
	}
	return "draw"
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	row, err := strconv.Atoi(r.Form.Get("row"))
	if err != nil {
		http.Error(w, "Invalid row value", http.StatusBadRequest)
		return
	}
	col, err := strconv.Atoi(r.Form.Get("col"))
	if err != nil {
		http.Error(w, "Invalid col value", http.StatusBadRequest)
		return
	}
	if board[row][col] != empty {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}
	board[row][col] = currentPlayer
	winner := checkWin()
	if winner != "" {
		tmpl.Execute(w, struct {
			Board         [3][3]string
			Winner        string
			CurrentPlayer string
		}{
			Board:         board,
			Winner:        winner,
			CurrentPlayer: currentPlayer,
		})
		return
	}
	if currentPlayer == player1 {
		currentPlayer = player2
	} else {
		currentPlayer = player1
	}
	tmpl.Execute(w, struct {
		Board         [3][3]string
		Winner        string
		CurrentPlayer string
	}{
		Board:         board,
		Winner:        "",
		CurrentPlayer: currentPlayer,
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, struct {
			Board         [3][3]string
			Winner        string
			CurrentPlayer string
		}{
			Board:         board,
			Winner:        "",
			CurrentPlayer: currentPlayer,
		})
	})
	http.HandleFunc("/play", handlePlay)
	http.ListenAndServe(":5757", nil)
}