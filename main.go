package main

import (
	"html/template"
	"math/rand"
	"net/http"
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

func init() {
	var err error
	rand.Seed(time.Now().UnixNano())
	tpl, err = template.New("").Funcs(template.FuncMap{
		"add1": func(i int) int { return i + 1 },
	}).ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}
}

// Route accueil
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

// Route initGame : affiche la page de saisie des noms
func initHandler(w http.ResponseWriter, r *http.Request) {
	var j1, j2 string

	// Si on arrive depuis une redirection (POST -> GET), garder les noms
	if nomJoueur1 != "" {
		j1 = nomJoueur1
	}
	if nomJoueur2 != "" {
		j2 = nomJoueur2
	}

	data := struct {
		Joueur1 string
		Joueur2 string
		Error   string
	}{
		Joueur1: j1,
		Joueur2: j2,
		Error:   "Veuillez saisir les noms des joueurs avant de commencer !",
	}

	tpl.ExecuteTemplate(w, "initGame.html", data)
}
