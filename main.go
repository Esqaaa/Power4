package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Plateau 6x7 (0=vide, 1=rouge, 2=jaune)
var board [6][7]int
var currentPlayer = 1
var gameOver = false

var tpl *template.Template

func init() {
	var err error
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}
	fmt.Println("✓ Templates chargés avec succès")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Board         [6][7]int
		CurrentPlayer int
		GameOver      bool
		Message       string
	}{
		Board:         board,
		CurrentPlayer: currentPlayer,
		GameOver:      gameOver,
	}

	if gameOver {
		if winner := checkWin(); winner != 0 {
			data.Message = "Joueur " + strconv.Itoa(winner) + " gagne ! 🎉"
		} else {
			data.Message = "Match nul ! 😅"
		}
	} else {
		color := map[int]string{1: "(Rouge)", 2: "(Jaune)"}[currentPlayer]
		data.Message = "Tour du Joueur " + strconv.Itoa(currentPlayer) + " " + color
	}

	err := tpl.ExecuteTemplate(w, "jeu.html", data)
	if err != nil {
		fmt.Println("ERREUR Template:", err)
		http.Error(w, "Erreur template", http.StatusInternalServerError)
	}
}

func checkWin() int {
	for row := 0; row < 6; row++ {
		for col := 0; col <= 3; col++ {
			if board[row][col] != 0 && board[row][col] == board[row][col+1] &&
				board[row][col] == board[row][col+2] && board[row][col] == board[row][col+3] {
				return board[row][col]
			}
		}
	}
	for col := 0; col < 7; col++ {
		for row := 0; row <= 2; row++ {
			if board[row][col] != 0 && board[row][col] == board[row+1][col] &&
				board[row][col] == board[row+2][col] && board[row][col] == board[row+3][col] {
				return board[row][col]
			}
		}
	}
	for row := 0; row <= 2; row++ {
		for col := 0; col <= 3; col++ {
			if board[row][col] != 0 && board[row][col] == board[row+1][col+1] &&
				board[row][col] == board[row+2][col+2] && board[row][col] == board[row+3][col+3] {
				return board[row][col]
			}
		}
	}
	for row := 3; row < 6; row++ {
		for col := 0; col <= 3; col++ {
			if board[row][col] != 0 && board[row][col] == board[row-1][col+1] &&
				board[row][col] == board[row-2][col+2] && board[row][col] == board[row-3][col+3] {
				return board[row][col]
			}
		}
	}
	return 0
}
