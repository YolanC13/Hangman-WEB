package hangman

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var HangmanChar []string
var UsedLetters []string

var lives = 9
var PlayerLives *int = &lives
var Characters = &HangmanChar
var Letters = &UsedLetters

var ASCIIArts map[string]string
var ASCIIArtsPtr = &ASCIIArts

var WordList []string
var WordListPtr = &WordList

var fileImported string
var FileImportedPtr = &fileImported

var PlayerScore = 0
var PlayerScorePtr = &PlayerScore

var WordStrk = 0
var WordStreakPtr = &WordStrk

var isWinning bool = false
var Winner = &isWinning

func AsciiArtsInit() {
	// Initialise ASCIIArts
	*ASCIIArtsPtr = make(map[string]string)

	// Lis la liste des fichiers du dossier AsciiArt
	files, err := os.ReadDir("./AsciiArt")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier AsciiArt")
		fmt.Println(err)
		listDir()
	}

	// Associe chaque fichier à son contenu dans la map AsciiArt
	for _, file := range files {
		asciiArt, err := ioutil.ReadFile("./AsciiArt/" + file.Name())
		if err != nil {
			fmt.Println("Erreur lors de la lecture du fichier " + file.Name())
			fmt.Println(err)
			continue
		}
		(*ASCIIArtsPtr)[file.Name()] = string(asciiArt)
	}
}

func listDir() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
