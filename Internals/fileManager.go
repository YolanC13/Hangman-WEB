package hangman

import (
	"bufio"
	"os"
)

func LoadTextFile(fileToLoad string) []string {
	//Section 1
	file, err := os.Open(fileToLoad)
	if err != nil {
		return nil
	}
	defer file.Close()

	r := bufio.NewReader(file)
	words := []string{}
	// Section 2
	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
			words = append(words, string(line))
		}
		if err != nil {
			return words
		}
	}
}

func FileExists(filename string) bool {

	info, err := os.Stat(filename)

	if os.IsNotExist(err) {

		return false

	}

	return !info.IsDir()

}
