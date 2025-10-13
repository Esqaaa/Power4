package main

import "fmt"

const (
	LIGNES  = 6
	COLONNE = 7
)

type plateau [LIGNES][COLONNE]int

func main() {
	var p plateau
	fmt.Println("Bienvenue au Puissance 4 !")
	afficherPlateau(p)
}

func afficherPlateau(p plateau) {
	fmt.Println(" 1  2  3  4  5  6  7 ")
	fmt.Println("---------------------")

	for l := 0; l < LIGNES; l++ {
		for c := 0; c < COLONNE; c++ {
			jeton := "."
			if p[l][c] == 1 {
				jeton = "X" // Si la case contient 1, c'est le joueur 1
			} else if p[l][c] == 2 {
				jeton = "O" // Si la case contient 2, c'est le joueur 2
			}

			// Ici, on affiche la variable 'jeton' (qui est bien une chaîne de caractères)
			fmt.Printf(" %s ", jeton)
		}
		fmt.Println() // On passe à la ligne suivante après chaque rangée
	}
	fmt.Println("---------------------")
}
