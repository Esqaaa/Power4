package main

import (
	"fmt"
	"html/template"
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
