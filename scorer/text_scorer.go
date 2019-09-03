package scorer

import (
	"bufio"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const monogramLength = 26
const quadgramLength = 26 * 26 * 26 * 26 // 26^4

// TextScorer provides methods for scoring text
type TextScorer struct {
	Monograms [monogramLength]float64
	Quadgrams [quadgramLength]float64
}

// Takes an ngram and calculates its index in a sorted array
// Text must be in the form of a slice of uppercase alphabetic runes
func offset(text []byte, offset, length int) int {
	sum := 0
	for i := 0; i < length; i++ {
		c := text[i+offset]
		sum *= 26
		sum += int(c - 'A')
	}
	return sum
}

func loadNgrams(path string, processNgram func([]byte, float64)) {
	if log.GetLevel() == log.TraceLevel {
		wd, err := os.Getwd()
		if err == nil {
			log.Tracef("Current working directory is %s", wd)
		} else {
			log.Traceln("Failed to fetch working directory")
		}
	}

	file, err := os.Open(path)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	defer file.Close()

	if err != nil {
		log.Fatalf("failed opening file: %s", path)
	}

	sum := float64(0)

	for scanner.Scan() {
		ngramRaw := scanner.Text()
		ngram := []byte(ngramRaw)
		scanner.Scan()
		ngramScore, err := strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			log.Fatalf("failed reading score for ngram %s: %s is not a valid float", ngramRaw, scanner.Text())
		}
		sum += ngramScore
		processNgram(ngram, ngramScore)
	}
}

// LoadMonograms loads monograms into memory from a given file path
// kills the process if the file cannot be read
// quadgrams are formatted as an array of float32
func LoadMonograms(path string) [monogramLength]float64 {
	log.Debugln("Loading monograms from", path)
	var monograms [monogramLength]float64
	loadNgrams(path, func(monogram []byte, monogramScore float64) {
		monograms[offset(monogram, 0, 1)] = math.Log10(monogramScore)
	})

	return monograms
}

// LoadQuadgrams loads quadgrams into memory from a given file path
// kills the process if the file cannot be read
// quadgrams are formatted as an array of float32
func LoadQuadgrams(path string) [quadgramLength]float64 {
	log.Debugln("Loading quadgrams from", path)
	var quadgrams [quadgramLength]float64
	loadNgrams(path, func(quadgram []byte, quadgramScore float64) {
		quadgrams[offset(quadgram, 0, 4)] = math.Log10(quadgramScore)
	})

	return quadgrams
}

// QuadgramScore guesses how likely a text is to be english based on its quadgrams
// a higher score is better
func (s *TextScorer) QuadgramScore(str string) float64 {
	log.Tracef("Quadgram scoring %s\n", str)
	fitness := 0.0
	str = strings.ToUpper(str)
	// TODO Move this regex to a property on the TextScorer struct
	reg, err := regexp.Compile("[^A-Z]+")
	if err != nil {
		log.Fatalf("Error compiling regex: %s", err)
	}
	str = reg.ReplaceAllString(str, "")

	for i := 1; i < len(str)-3; i++ {
		fitness += s.Quadgrams[offset([]byte(str), i, 4)]
	}

	return fitness
}
