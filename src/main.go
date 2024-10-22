package main

import (
	"fmt"
	hangman "hangman/Internals"
	"math/rand"
	"net/http"
	"slices"
	"strings"
	"text/template"
)

type UserInfo struct {
	Username   *string
	WordStreak *int
}

type GameVariables struct {
	PlayerLives int
	Word        *[]string
	Letters     *string
	UsedLetters []string
	HangmanChar []string
	Winner      *bool
	PlayerScore int
	WordStreak  int
}

type Leaderboard struct {
	Users *[]string
}

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
		*hangman.PlayerLives -= 1
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

	UserInfo := UserInfo{
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

	http.HandleFunc("/mainMenu/newGame", func(w http.ResponseWriter, r *http.Request) {
		ResetScore()
		temp.ExecuteTemplate(w, "newGame", nil)
	})

	//Jeu
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		gameVars := GameVariables{
			PlayerLives: *hangman.PlayerLives,
			Word:        hangman.Characters,
			Letters:     new(string),
			UsedLetters: hangman.UsedLetters,
			HangmanChar: hangman.HangmanChar,
			PlayerScore: *hangman.PlayerScorePtr,
			WordStreak:  *hangman.WordStreakPtr,
		}

		for i := 0; i < len(*hangman.Letters); i++ {
			*gameVars.Letters += (*hangman.Letters)[i]
			*gameVars.Letters += " "
		}

		if *hangman.PlayerLives == 0 {
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
			UsedLetters: hangman.UsedLetters,
			HangmanChar: hangman.HangmanChar,
			Winner:      hangman.Winner,
			PlayerScore: *new(int),
			WordStreak:  *hangman.WordStreakPtr,
		}

		if *hangman.PlayerLives <= 0 {
			*hangman.Winner = false
		} else {
			*hangman.Winner = true

			//Calcul du Word Streak
			*hangman.WordStreakPtr += 1
			gameVars.WordStreak = *hangman.WordStreakPtr

			//Calcul du score
			for i := 0; i < len(*hangman.Characters); i++ {
				*hangman.PlayerScorePtr += 1 * (len(*hangman.Characters) - i) * 10
				if *hangman.WordStreakPtr != 0 {
					*hangman.PlayerScorePtr += 5 * *hangman.WordStreakPtr
				}
			}
			gameVars.PlayerScore = *hangman.PlayerScorePtr

			*UserInfo.WordStreak += 1
		}

		gameVars.Winner = hangman.Winner
		temp.ExecuteTemplate(w, "resultat", gameVars)
	})

	//Leaderboard
	http.HandleFunc("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		Leaderboard := Leaderboard{
			Users: &[]string{},
		}

		for i := 0; i < len(hangman.ReadFileAndReturn()); i++ {
			ligne := hangman.ReadFileAndReturn()
			mots := strings.Fields(strings.Join(ligne, " "))

			*Leaderboard.Users = append(*Leaderboard.Users, mots[i]+" "+mots[i+1])
		}

		temp.ExecuteTemplate(w, "leaderboard", Leaderboard)
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
	fichier := "words.txt"
	*hangman.WordListPtr = hangman.LoadTextFile(fichier)
}

func ResetScore() {
	*hangman.PlayerScorePtr = 0
	*hangman.WordStreakPtr = 0
}
