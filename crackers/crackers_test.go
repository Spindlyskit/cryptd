package crackers

import (
	"bytes"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/spindlyskit/cryptd/scorer"
)

type Test struct {
	t                       *testing.T
	name                    string
	cracker                 Cracker
	solveCt, solvePlainText []byte
	solveKey                string
	crackCt, crackPlainText []byte
	crackKey                string
}

var textScorer scorer.TextScorer

func init() {
	log.SetLevel(log.DebugLevel)
	textScorer = scorer.TextScorer{
		Monograms: scorer.LoadMonograms("../assets/english_monograms.txt"),
		Quadgrams: scorer.LoadQuadgrams("../assets/english_quadgrams.txt"),
	}
}

func Suite(test *Test) {
	t := test.t
	cracker := test.cracker

	// Solve with key test
	solveReply := ReplyData{}
	solveOptions := SolveOptions{
		CT:  test.solveCt,
		Key: test.solveKey,
	}
	cracker.Decrypt(solveOptions, &solveReply)
	t.Logf("%s - Solve: got key %s",
		test.name, solveReply.Key)
	if solveReply.Key != test.solveKey {
		t.Errorf("%s - Solve: expected key to be %s but got %s",
			test.name, test.solveKey, solveReply.Key)
	}
	t.Logf("%s - Solve: got plaintext %s",
		test.name, solveReply.PlainText)
	if !bytes.Equal(solveReply.PlainText, test.solvePlainText) {
		t.Errorf("%s - Solve: expected plaintext to be %s but got %s",
			test.name, test.solvePlainText, solveReply.PlainText)
	}

	// Crack without key test
	crackReply := ReplyData{}
	crackOptions := SolveOptions{
		CT: test.crackCt,
	}
	cracker.Decrypt(crackOptions, &crackReply)
	t.Logf("%s - Crack: got plaintext %s",
		test.name, crackReply.PlainText)
	if !bytes.Equal(crackReply.PlainText, test.crackPlainText) {
		t.Errorf("%s - Crack: expected plaintext to be %s but got %s",
			test.name, test.crackPlainText, crackReply.PlainText)
	}
	t.Logf("%s - Crack: got key %s",
		test.name, crackReply.Key)
	if crackReply.Key != test.crackKey {
		t.Errorf("%s - Crack: expected used key to be %s but got %s",
			test.name, test.crackKey, test.crackCt)
	}
}

func TestCaesar(t *testing.T) {
	test := Test{
		t:              t,
		cracker:        &CaesarCracker{Scorer: textScorer},
		name:           "caesar",
		solveCt:        []byte("axeeh phkew"),
		solveKey:       "8",
		solvePlainText: []byte("ifmmp xpsme"),
		crackCt:        []byte("Pm ol ohk hufaopun jvumpkluaphs av zhf, ol dyval pa pu jpwoly, aoha pz, if zv johunpun aol vykly vm aol slaalyz vm aol hswohila, aoha uva h dvyk jvbsk il thkl vba."),
		crackPlainText: []byte("If he had anything confidential to say, he wrote it in cipher, that is, by so changing the order of the letters of the alphabet, that not a word could be made out."),
		crackKey:       "19",
	}
	Suite(&test)
}
