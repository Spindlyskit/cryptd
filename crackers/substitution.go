package crackers

import (
	"math"
	"math/rand"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spindlyskit/cryptd/scorer"
	"github.com/spindlyskit/cryptd/util"
)

const iterations = 5
const internalIterations = 1000

type singleCrackResult struct {
	key     []byte
	fitness float64
}

// SubstitutionCracker provides methods to crack monoalphabetic substitution ciphers
type SubstitutionCracker struct {
	Scorer scorer.TextScorer
}

// Crack cracks a substitution cipher without a key
func (cracker *SubstitutionCracker) Crack(ct []byte) ([]byte, string, error) {
	var usedKey []byte
	bestFitness := math.Inf(-1)

	var wg sync.WaitGroup
	var results [iterations]*singleCrackResult

	for i := 0; i < iterations; i++ {
		results[i] = &singleCrackResult{}
		wg.Add(1)
		go cracker.crackSingle(ct, &wg, results[i])
	}

	// Block until we have results
	wg.Wait()

	for _, result := range results {
		if result.fitness > bestFitness {
			usedKey = result.key
			bestFitness = result.fitness
		}
	}

	return cracker.decode(ct, usedKey), string(usedKey), nil
}

// Decrypt decrypts a substitution cipher with the given options
// Exported by rpc
func (cracker *SubstitutionCracker) Decrypt(options SolveOptions, reply *ReplyData) error {
	return nil
}

// Solve solves a substitution cipher with the given key
func (cracker *SubstitutionCracker) Solve(ct []byte, k string) ([]byte, error) {
	return []byte{}, nil
}

func (cracker *SubstitutionCracker) decode(ct []byte, key []byte) []byte {
	var decoded []byte
	// Normalize the cipher text so we can apply the key to it
	ct = util.NormalizeB(ct)

	for _, char := range ct {
		charIndex := char - 'A'
		decoded = append(decoded, key[charIndex])
	}

	return decoded
}

func (cracker *SubstitutionCracker) randKey() []byte {
	key := make([]byte, 26)
	for i, v := range rand.Perm(26) {
		key[v] = util.Alphabet[i]
	}
	log.Debugf("Generated substitution key %s", key)
	return key
}

func (cracker SubstitutionCracker) swap(text []byte) {
	a := rand.Intn(len(text) - 1)
	b := rand.Intn(len(text) - 1)

	text[a], text[b] = text[b], text[a]
}

func (cracker *SubstitutionCracker) crackSingle(ct []byte, wg *sync.WaitGroup, result *singleCrackResult) {
	bestKey := cracker.randKey()
	decoded := cracker.decode(ct, bestKey)
	bestFitness := cracker.Scorer.QuadgramScore(string(decoded))

	workingKey := make([]byte, len(bestKey))
	copy(workingKey, bestKey)
	for i := 0; i < internalIterations; i++ {
		cracker.swap(workingKey)

		decoded = cracker.decode(ct, workingKey)
		workingFitness := cracker.Scorer.QuadgramScore(string(decoded))
		if workingFitness > bestFitness {
			i = 0
			bestFitness = workingFitness
			copy(bestKey, workingKey)
		} else {
			// Revert key
			copy(workingKey, bestKey)
		}
	}
	result.key = bestKey
	result.fitness = bestFitness
	wg.Done()
}
