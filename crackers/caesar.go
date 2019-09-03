package crackers

import (
	"errors"
	"strconv"

	"github.com/spindlyskit/cryptd/scorer"
)

// CaesarCracker provides methods for cracking caesar encoded text
type CaesarCracker struct {
	Scorer scorer.TextScorer
}

// Crack a caesar cipher without a key
func (cracker *CaesarCracker) Crack(ct []byte) ([]byte, string, error) {
	decoded := []byte{}
	var usedKey byte
	fitness := 0.0

	for i := byte(0); i < 26; i++ {
		maybeDecoded := cracker.decode(ct, i)
		maybeFitness := cracker.Scorer.QuadgramScore(string(maybeDecoded))
		if len(decoded) == 0 || maybeFitness > fitness {
			decoded = maybeDecoded
			usedKey = i
			fitness = maybeFitness
		}
	}

	return decoded, strconv.Itoa(int(usedKey)), nil
}

// Decrypt decrypts a caesar cipher with the given options
// Exported by rpc
func (cracker *CaesarCracker) Decrypt(options SolveOptions, reply *ReplyData) error {
	var plaintext []byte
	var usedKey string
	var err error
	if options.Key == "" {
		plaintext, usedKey, err = cracker.Crack(options.CT)
	} else {
		usedKey = options.Key
		plaintext, err = cracker.Solve(options.CT, options.Key)
	}
	if err != nil {
		return err
	}
	*reply = ReplyData{
		Key:       usedKey,
		PlainText: plaintext,
	}
	return nil
}

// Solve solves a caesar cipher with the given key
func (cracker *CaesarCracker) Solve(ct []byte, k string) ([]byte, error) {
	iKey, err := strconv.Atoi(k)
	if err != nil {
		return []byte{}, errors.New("caesar: invalid key")
	}
	key := byte(iKey)
	if key < 0 || key > 25 {
		return []byte{}, errors.New("caesar: key must be between 0 and 25")
	}
	return cracker.decode(ct, key), nil
}

func (cracker *CaesarCracker) decode(ct []byte, key byte) []byte {
	decoded := []byte{}
	for _, char := range []byte(ct) {
		decoded = append(decoded, shift(char, key))
	}

	return decoded
}

func shift(char byte, shift byte) byte {
	var a, z byte

	switch {
	case 'a' <= char && char <= 'z':
		a, z = 'a', 'z'
	case 'A' <= char && char <= 'Z':
		a, z = 'A', 'Z'
	default:
		return char
	}
	return (char-a+shift)%(z-a+1) + a
}
