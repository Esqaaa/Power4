package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Variable globale
var tpl *template.Template

// Grille dynamique
var board [][]int
var rows, cols int
var numBlocks int

// État du jeu
var gameOver = false
var currentDifficulty string
var turnCount int
var gravityNormal = true // true = normale, false = inversée

// Tout ce qui concerne les joueurs
var winner int
var currentPlayer = 1
var nomJoueur1, nomJoueur2 string
var joueur1Wins, joueur2Wins int

// Initialisation de la grille selon la difficulté choisie
func initBoard(difficulty string) {
	switch difficulty {
	case "basique":
		rows, cols = 6, 7
		numBlocks = 0
	case "easy":
		rows, cols = 6, 7
		numBlocks = 0
	case "normal":
		rows, cols = 6, 9
		numBlocks = 0
	case "hard":
		rows, cols = 7, 8
		numBlocks = 0
	case "chaos":
		rows = rand.Intn(4) + 6      // 6 à 9 lignes
		cols = rand.Intn(4) + 6      // 6 à 9 colonnes
		numBlocks = rand.Intn(6) + 3 // 3 à 8 blocs
	case "blocfou":
		rows, cols = 6, 7
		numBlocks = 0
	default:
		rows, cols = 6, 7
		numBlocks = 0
	}

	// Création de la grille vide
	board = make([][]int, rows)
	for i := range board {
		board[i] = make([]int, cols)
	}
}

// Place un nombre de blocs aléatoires sur le plateau
func PlaceRandomBlocs(nb int) {
	placed := 0
	for placed < nb {
		row := rand.Intn(rows)
		col := rand.Intn(cols)
		if board[row][col] == 0 {
			board[row][col] = 3 // 3 = bloc fixe
			placed++
		}
	}
}
