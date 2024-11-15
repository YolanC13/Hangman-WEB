package hangman

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
