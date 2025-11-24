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

// Route start : noms, difficulté, blocs
func startHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/init", http.StatusSeeOther)
		return
	}

	// Récupération des noms
	if j1 := r.FormValue("joueur1"); j1 != "" {
		nomJoueur1 = j1
	}
	if j2 := r.FormValue("joueur2"); j2 != "" {
		nomJoueur2 = j2
	}

	// Vérifier que les noms sont bien remplis
	if nomJoueur1 == "" || nomJoueur2 == "" {
		http.Redirect(w, r, "/init", http.StatusSeeOther)
		return
	}

	//  Nettoyage du champ difficulté
	currentDifficulty = strings.TrimSpace(r.FormValue("difficulty"))
	nbBlocks, _ := strconv.Atoi(r.FormValue("blocks"))

	// Initialisation de la grille
	initBoard(currentDifficulty)

	// Reset de l’état du jeu
	currentPlayer = 1
	gameOver = false
	turnCount = 0

	// Gravité toujours normale au début
	gravityNormal = true

	// Placement des blocs
	if nbBlocks > 0 {
		PlaceRandomBlocs(nbBlocks)
	}

	// Redirection vers le jeu
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

// Route jeu
func gameHandler(w http.ResponseWriter, r *http.Request) {
	winner = checkWin()
	gameOver = winner != 0 || isDraw()

	data := struct {
		Joueur1       string
		Joueur2       string
		Joueur1Wins   int
		Joueur2Wins   int
		CurrentPlayer int
		Board         [][]int
		GameOver      bool
		Message       string
		Cols          []int
		Difficulty    string
		Winner        int
		GravityNormal bool
		Turn          int
	}{
		Joueur1:       nomJoueur1,
		Joueur2:       nomJoueur2,
		Joueur1Wins:   joueur1Wins,
		Joueur2Wins:   joueur2Wins,
		CurrentPlayer: currentPlayer,
		Board:         board,
		GameOver:      gameOver,
		Cols:          make([]int, cols),
		Difficulty:    currentDifficulty,
		Winner:        winner,
		GravityNormal: gravityNormal,
		Turn:          turnCount,
	}

	for i := 0; i < cols; i++ {
		data.Cols[i] = i
	}

	if !gameOver {
		current := map[int]string{1: nomJoueur1, 2: nomJoueur2}[currentPlayer]
		data.Message = "Tour de " + current
	} else {
		if winner != 0 {
			if winner == 1 {
				data.Message = nomJoueur1 + " a gagné ! 🎉"
			} else {
				data.Message = nomJoueur2 + " a gagné ! 🎉"
			}
		} else {
			data.Message = "Match nul ! 😅"
		}
	}

	if nomJoueur1 == "" || nomJoueur2 == "" {
		http.Redirect(w, r, "/init", http.StatusSeeOther)
		return
	}

	if err := tpl.ExecuteTemplate(w, "jeu.html", data); err != nil {
		http.Error(w, "Erreur de rendu du template", http.StatusInternalServerError)
	}
}

// Route play
func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode POST seulement", http.StatusMethodNotAllowed)
		return
	}

	colStr := r.FormValue("column")
	col, err := strconv.Atoi(colStr)
	if err != nil || col < 0 || col >= cols {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}

	if !gameOver {
		if dropPiece(col) { // retourne true si la pièce a bien été placée
			currentPlayer = 3 - currentPlayer
			turnCount++

			// Changer la gravité toutes les 6 tours (sauf basique)
			if currentDifficulty != "basique" && turnCount%6 == 0 {
				gravityNormal = !gravityNormal
			}

			// Mode BlocFou : ajout d’un bloc aléatoire
			if currentDifficulty == "blocfou" {
				addRandomBlock()
			}

			winner = checkWin()
			gameOver = winner != 0 || isDraw()
			if winner == 1 {
				joueur1Wins++
			} else if winner == 2 {
				joueur2Wins++
			}
		}
	}

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

// dropPiece retourne true si placement réussi
func dropPiece(col int) bool {
	if gravityNormal {
		for row := rows - 1; row >= 0; row-- {
			if board[row][col] == 0 {
				board[row][col] = currentPlayer
				return true
			}
		}
	} else {
		for row := 0; row < rows; row++ {
			if board[row][col] == 0 {
				board[row][col] = currentPlayer
				return true
			}
		}
	}
	return false // colonne pleine
}

func addRandomBlock() {
	emptyCells := []struct{ row, col int }{}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if board[r][c] == 0 {
				emptyCells = append(emptyCells, struct{ row, col int }{r, c})
			}
		}
	}
	if len(emptyCells) > 0 {
		cell := emptyCells[rand.Intn(len(emptyCells))]
		board[cell.row][cell.col] = 3 // 3 = bloc fixe
	}
}

func checkWin() int {
	// Horizontal
	for row := 0; row < rows; row++ {
		for col := 0; col <= cols-4; col++ {
			p := board[row][col]
			if p != 0 && p == board[row][col+1] && p == board[row][col+2] && p == board[row][col+3] {
				return p
			}
		}
	}
	// Vertical
	for col := 0; col < cols; col++ {
		for row := 0; row <= rows-4; row++ {
			p := board[row][col]
			if p != 0 && p == board[row+1][col] && p == board[row+2][col] && p == board[row+3][col] {
				return p
			}
		}
	}
	// Diagonale descendante
	for row := 0; row <= rows-4; row++ {
		for col := 0; col <= cols-4; col++ {
			p := board[row][col]
			if p != 0 && p == board[row+1][col+1] && p == board[row+2][col+2] && p == board[row+3][col+3] {
				return p
			}
		}
	}
	// Diagonale montante
	for row := 3; row < rows; row++ {
		for col := 0; col <= cols-4; col++ {
			p := board[row][col]
			if p != 0 && p == board[row-1][col+1] && p == board[row-2][col+2] && p == board[row-3][col+3] {
				return p
			}
		}
	}
	return 0
}

func isDraw() bool {
	for col := 0; col < cols; col++ {
		if board[0][col] == 0 {
			return false
		}
	}
	return true
}

// Reset du plateau
func resetBoard() {
	for i := range board {
		for j := range board[i] {
			board[i][j] = 0
		}
	}
	currentPlayer = 1
	gameOver = false

	// Reinitialisation des bonus
	turnCount = 0
	gravityNormal = true
}

// Route revanche : reset du plateau mais conserve noms et difficulté
func rematchHandler(w http.ResponseWriter, r *http.Request) {
	// Réinitialiser la grille et l'état sans toucher aux noms ni à la difficulté
	initBoard(currentDifficulty)
	currentPlayer = 1
	gameOver = false
	turnCount = 0

	// Gravité selon la difficulté
	switch currentDifficulty {
	case "basique":
		gravityNormal = true
	case "easy":
		gravityNormal = true
	case "normal", "hard", "chaos", "blocfou":
		gravityNormal = false
	default:
		gravityNormal = true
	}

	// Remettre les blocs si nécessaire
	if numBlocks > 0 {
		PlaceRandomBlocs(numBlocks)
	}

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

// Route reset
func resetHandler(w http.ResponseWriter, r *http.Request) {
	resetBoard()
	nomJoueur1 = ""
	nomJoueur2 = ""
	joueur1Wins = 0
	joueur2Wins = 0
	currentPlayer = 1
	gameOver = false
	currentDifficulty = ""
	winner = 0
	numBlocks = 0
	http.Redirect(w, r, "/init", http.StatusSeeOther)
}

func main() {
	// Routes
	http.HandleFunc("/init", initHandler)
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/reset", resetHandler)
	http.HandleFunc("/rematch", rematchHandler)
	http.HandleFunc("/", indexHandler)

	// Servir le CSS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("Styles"))))

	fmt.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
