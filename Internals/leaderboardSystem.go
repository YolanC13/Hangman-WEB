package hangman

import (
	"bufio"
	"log"
	"os"
)

func ReadFileAndReturn() []string {
	// Ouvre le fichier
	file, err := os.Open("leaderboardStat.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Utilise un scanner pour lire le fichier ligne par ligne
	ligne := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Récupère la ligne
		ligne = append(ligne, scanner.Text())
	}

	// Gère les erreurs de lecture
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ligne
}
