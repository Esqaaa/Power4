package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
<<<<<<< HEAD
=======
	"strings"
	"time"
>>>>>>> 36d237ff6c4854d7f6e3b9d19470850ec838ae6d
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
