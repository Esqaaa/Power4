# 🎮 Projet Puissance 4 Web

Bienvenue dans le dépôt **Projet Puissance 4 Web**, un jeu de **Puissance 4** jouable directement dans le navigateur avec HTML, CSS et Go.

> 🔴⚪ Connectez 4 jetons et remportez la partie !

---

## 🚀 Présentation

Ce projet est un **jeu web interactif** pour deux joueurs, avec plusieurs modes de difficulté et fonctionnalités uniques :

### Fonctionnalités principales

* 🖥️ Interface web simple et responsive
* 🔴⚪ Jeu pour **2 joueurs** sur le même écran
* ⚡ **Logique serveur en Go (Goland)**
* 🎨 Design coloré et dynamique avec HTML/CSS
* 💨 **Modes de gravité** : normale ou inversée (les pions tombent de bas en haut)
* 🧱 **Mode Bloc Fou** : des blocs aléatoires apparaissent à chaque tour
* 🎲 **Mode Chaos** : grille et nombre de blocs aléatoires
* 📊 **Différentes difficultés** : Easy, Normal, Hard, Chaos, BlocFou

### Difficultés et paramètres

| Mode    | Grille    | Nombre de blocs            |
| ------- | --------- | -------------------------- |
| Easy    | 6x7       | 3                          |
| Normal  | 6x9       | 5                          |
| Hard    | 7x8       | 7                          |
| Chaos   | aléatoire | aléatoire                  |
| BlocFou | 6x7       | apparaissent à chaque tour |

---

## 🛠️ Installation

### 1. Cloner le dépôt

```bash
git clone https://github.com/Esqaaa/Power4.git
cd puissance4-web
```

### 2. Installer les dépendances Go

```bash
go mod tidy
```

### 3. Lancer le serveur

```bash
go run main.go
```
ou 
```bash
go run power4
```

### 4. Jouer

Ouvrez votre navigateur et allez sur `http://localhost:8080`.

---

## 👨‍💻 Types de langages

* Langages : **Go (Goland), HTML, CSS**
* Serveur web : Go standard net/http
* Templates : Go HTML templates
* Frontend : HTML + CSS pour le plateau et l’interface utilisateur

---

## 🧾 Licence

Projet développé à titre éducatif — © 2025 **Equipe Puissance 4 Web**.

---

## 🏫 Crédits

Projet réalisé dans le cadre d’un projet au sein d'Ynov Campus Strasbourg B1 Info-Cybersécurité 2025/2026 par SCHMALTZ Hugo et SCHMITT Gabriel
