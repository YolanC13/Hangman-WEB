package main

import (
	"fmt"
	hangman "hangman/Internals"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

type User struct {
	Username   *string
	Score      *int
	WordStreak *int
}

type GameVariables struct {
	PlayerLives int
	Word        *[]string
	Letters     *string
	UsedLetters *string
	HangmanChar []string
	PlayerScore int
	WordStreak  int
	RandomMsg   string
}

type Leaderboard struct {
	Users []User
}

var userInfo = User{}

func main() {
	fichier := "words.txt"
	*hangman.WordListPtr = hangman.LoadTextFile(fichier)
	InitialiseServer()
}

func InitializeVariables(text string) {
	hangman.Characters = &hangman.HangmanChar
	*hangman.PlayerLives = 9
	for i := 0; i < len(text); i++ {
		if text[i] == ' ' {
			hangman.HangmanChar = append(hangman.HangmanChar, " ")
			*hangman.Letters = append(*hangman.Letters, " ")
		} else {
			hangman.HangmanChar = append(hangman.HangmanChar, strings.ToLower(string(text[i])))
			*hangman.Letters = append(*hangman.Letters, "_")
		}
	}
	*hangman.Characters = hangman.HangmanChar

	//AJOUTE DES LETTRES ALEATOIREMENT
	if len(*hangman.Characters) > 9 {
		for i := 0; i < 2; i++ {
			x := rand.Intn(len(*hangman.Characters))
			y := (*hangman.Characters)[x]
			for i := 0; i < len(*hangman.Characters); i++ {
				if (*hangman.Characters)[i] == y {
					(*hangman.Letters)[i] = y
				}
			}
			hangman.UsedLetters = append(hangman.UsedLetters, (*hangman.Letters)[x])
		}
	} else if len(*hangman.Characters) > 5 {

		x := rand.Intn(len(*hangman.Characters))
		y := (*hangman.Characters)[x]
		for i := 0; i < len(*hangman.Characters); i++ {
			if (*hangman.Characters)[i] == y {
				(*hangman.Letters)[i] = y
			}
		}
		hangman.UsedLetters = append(hangman.UsedLetters, (*hangman.Letters)[x])
	}
}

func checkLetter(letter string) {
	foundLetter := false
	for i := 0; i < len(*hangman.Characters); i++ {
		if letter == (*hangman.Characters)[i] && !slices.Contains(hangman.UsedLetters, letter) {
			(*hangman.Letters)[i] = letter
			foundLetter = true
		}
	}
	if !foundLetter {
		if !slices.Contains(hangman.UsedLetters, letter) {
			*hangman.PlayerLives -= 1
		}
	}
	hangman.UsedLetters = append(hangman.UsedLetters, letter)
}

func IsLetter(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && r != ' ' {
			return false
		}
	}
	return true
}

func InitialiseServer() {
	temp, errTemp := template.ParseGlob("htmlStuff/*.html")
	if errTemp != nil {
		fmt.Printf("Error: %v\n", errTemp)
		return
	}

	UserInfo := User{
		Username:   new(string),
		WordStreak: new(int),
	}

	//Menu principal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/mainMenu", http.StatusSeeOther)
	})

	http.HandleFunc("/mainMenu", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "mainMenu", nil)
	})

	http.HandleFunc("/mainMenu/userForm", func(w http.ResponseWriter, r *http.Request) {
		ResetVariables()
		ResetScore()
		temp.ExecuteTemplate(w, "userForm", nil)
	})

	//Jeu
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		gameVars := GameVariables{
			PlayerLives: *hangman.PlayerLives,
			Word:        hangman.Characters,
			Letters:     new(string),
			UsedLetters: new(string),
			HangmanChar: hangman.HangmanChar,
			PlayerScore: *hangman.PlayerScorePtr,
			WordStreak:  *hangman.WordStreakPtr,
		}

		for i := 0; i < len(*hangman.Letters); i++ {
			*gameVars.Letters += (*hangman.Letters)[i]
			*gameVars.Letters += " "
		}

		if len(hangman.UsedLetters) > 0 {
			for i := 0; i < len(hangman.UsedLetters)-1; i++ {
				*gameVars.UsedLetters += (hangman.UsedLetters)[i] + ", "
			}
			*gameVars.UsedLetters += (hangman.UsedLetters)[len(hangman.UsedLetters)-1]
		} else {
			*gameVars.UsedLetters = "Aucune lettre utilisée"
		}

		if *hangman.PlayerLives < 1 {
			http.Redirect(w, r, "/game/resultat", http.StatusSeeOther)
		} else {
			for j := 0; j < len(*hangman.Letters); j++ {
				if (*hangman.Letters)[j] == "_" {
					break
				} else {
					if j == len(*hangman.Letters)-1 {
						http.Redirect(w, r, "/game/resultat", http.StatusSeeOther)
					}
				}
			}
		}

		temp.ExecuteTemplate(w, "game", gameVars)
	})

	http.HandleFunc("/game/initialisation/first", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			*UserInfo.Username = r.FormValue("pseudo")
			ResetVariables()
			InitializeVariables((*hangman.WordListPtr)[rand.Intn(len(*hangman.WordListPtr))])
			http.Redirect(w, r, "/game", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/game/initialisation", func(w http.ResponseWriter, r *http.Request) {
		ResetVariables()
		InitializeVariables((*hangman.WordListPtr)[rand.Intn(len(*hangman.WordListPtr))])
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	})

	http.HandleFunc("/game/treatment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			var x string = ""
			if len(hangman.HangmanChar) > 1 {
				for i := 0; i < len(hangman.HangmanChar); i++ {
					x += hangman.HangmanChar[i]
				}
			}

			if len(r.FormValue("letter")) == 1 {
				checkLetter(r.FormValue("letter"))
			} else if r.FormValue("letter") == x {
				copy(*hangman.Letters, hangman.HangmanChar)
			} else {
				*hangman.PlayerLives -= 2
				hangman.UsedLetters = append(hangman.UsedLetters, r.FormValue("letter"))
			}
			http.Redirect(w, r, "/game", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/game/resultat", func(w http.ResponseWriter, r *http.Request) {
		gameVars := GameVariables{
			PlayerLives: *hangman.PlayerLives,
			Word:        hangman.Characters,
			UsedLetters: new(string),
			Letters:     new(string),
			HangmanChar: hangman.HangmanChar,
			PlayerScore: *hangman.PlayerScorePtr,
			WordStreak:  *hangman.WordStreakPtr,
			RandomMsg:   "",
		}

		for i := 0; i < len(gameVars.HangmanChar); i++ {
			*gameVars.Letters += (gameVars.HangmanChar)[i]
		}

		userInfo := User{
			Username:   UserInfo.Username,
			WordStreak: UserInfo.WordStreak,
			Score:      hangman.PlayerScorePtr,
		}

		// Le joueur a gagné
		if *hangman.PlayerLives > 0 {
			// Augmente la streak de victoires
			*hangman.WordStreakPtr += 1
			gameVars.WordStreak = *hangman.WordStreakPtr

			// Calcul du score basé sur le mot
			for i := 0; i < len(*hangman.Characters); i++ {
				*hangman.PlayerScorePtr += 10 * (len(*hangman.Characters) - i)
			}

			// Ajoute un bonus pour la série de victoires
			if *hangman.WordStreakPtr > 0 {
				*hangman.PlayerScorePtr += 5 * *hangman.WordStreakPtr
			}

			gameVars.PlayerScore = *hangman.PlayerScorePtr

			// Met à jour les informations utilisateur
			*UserInfo.WordStreak += 1

			// Affiche un message de victoire aléatoire
			messageList := []string{
				"Bravo, vous avez gagné !",
				"Victoire !",
				"Vous avez gagné !",
				"Bien joué !",
			}
			gameVars.RandomMsg = messageList[rand.Intn(len(messageList))]

		} else { // Le joueur a perdu
			// Le joueur a perdu, enregistre le score dans le leaderboard
			if userInfo.Username != nil {
				err := hangman.AddScoreToFile(*userInfo.Username, *userInfo.Score, *userInfo.WordStreak, "leaderboardStat.txt")
				if err != nil {
					log.Printf("Erreur lors de l'ajout du score au fichier leaderboard: %v", err)
				}
			}

			// Réinitialise la série de victoires après la défaite
			*hangman.WordStreakPtr = 0
			*UserInfo.WordStreak = 0

			messageList := []string{
				"Dommage, vous avez perdu...",
				"Désolé, vous avez perdu...",
				"Vous avez perdu...",
				"Vous ferez mieux la prochaine fois...",
			}
			gameVars.RandomMsg = messageList[rand.Intn(len(messageList))]
		}

		// Exécute le template avec les variables du jeu
		err := temp.ExecuteTemplate(w, "resultat", gameVars)
		if err != nil {
			http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		}
	})

	//Leaderboard
	http.HandleFunc("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		// Créer un leaderboard vide
		Leaderboard := Leaderboard{
			Users: []User{}, // Plus besoin d'utiliser un pointeur ici
		}

		// Lire toutes les lignes une seule fois
		lignes := hangman.ReadFileAndReturn()

		// Parcourir chaque ligne
		for _, ligne := range lignes {
			// Séparer les champs par des espaces
			parts := strings.Fields(ligne)
			if len(parts) < 3 {
				// S'il n'y a pas assez d'éléments dans la ligne, ignorer
				fmt.Println("Ligne mal formatée :", ligne)
				continue
			}

			// Extraire les champs user, score et streak
			user := parts[0]
			scoreStr := parts[1]
			streakStr := parts[2]

			// Convertir les scores de string à int
			score, err := strconv.Atoi(scoreStr)
			if err != nil {
				fmt.Printf("Erreur lors de la conversion du score pour %s: %v\n", user, err)
				continue
			}
			streak, err := strconv.Atoi(streakStr)
			if err != nil {
				fmt.Printf("Erreur lors de la conversion du streak pour %s: %v\n", user, err)
				continue
			}

			// Ajouter l'utilisateur au leaderboard
			Leaderboard.Users = append(Leaderboard.Users, User{
				Username:   &user,
				Score:      &score,
				WordStreak: &streak,
			})
		}

		// Exécuter le template et afficher le leaderboard
		err := temp.ExecuteTemplate(w, "leaderboard", Leaderboard)
		if err != nil {
			http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		}
	})

	//Serveur
	http.Handle("/htmlStuff/", http.StripPrefix("/htmlStuff/", http.FileServer(http.Dir("./htmlStuff"))))
	http.ListenAndServe("localhost:8080", nil)
}

func ResetVariables() {
	*hangman.PlayerLives = 9
	hangman.UsedLetters = []string{}
	hangman.HangmanChar = []string{}
	hangman.Letters = &[]string{}
	hangman.Characters = &[]string{}
	hangman.WordListPtr = &[]string{}

	*hangman.WordListPtr = hangman.LoadTextFile("words.txt")
}

func ResetScore() {
	*hangman.PlayerScorePtr = 0
	*hangman.WordStreakPtr = 0
}
