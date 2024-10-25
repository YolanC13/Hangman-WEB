package hangman

import (
	"bufio"
	"fmt"
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

func AddScoreToFile(username string, score, streak int, filepath string) error {
	// Ouvrir le fichier en mode append (ajout)
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier : %v", err)
	}
	defer file.Close()

	// Formatage de la ligne à écrire dans le fichier
	line := fmt.Sprintf("%s %d %d\n", username, score, streak)

	// Écrire la ligne dans le fichier
	_, err = file.WriteString(line)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture dans le fichier : %v", err)
	}

	fmt.Println("Score ajouté avec succès au leaderboard")
	return nil
}
