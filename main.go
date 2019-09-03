package main

import (
	"math"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spindlyskit/cryptd/scorer"
)

// SEMVER is the current version of cryptd
const SEMVER = "0.0.0"

func main() {
	log.SetOutput(os.Stdout)

	debug := pflag.BoolP("trace", "t", false, "Prints low level debug information")
	port := pflag.IntP("port", "p", 570, "Specifies the port to listen on")
	verbose := pflag.BoolP("verbose", "v", false, "Prints verbose output")
	version := pflag.Bool("version", false, "Prints the cryptd version and exits")

	pflag.Parse()

	log.Infof("Cryptd %s\n", SEMVER)
	if *version {
		return
	}

	if *debug {
		log.SetLevel(log.TraceLevel)
		log.Debug("Using trace logging")
	} else if *verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Using verbose logging")
	}

	log.Infoln("cryptd rpc server starting on port", *port)

	textScorer := scorer.TextScorer{
		Monograms: scorer.LoadMonograms("assets/english_monograms.txt"),
		Quadgrams: scorer.LoadQuadgrams("assets/english_quadgrams.txt"),
	}

	log.Trace(textScorer.Quadgrams[0] == math.Log10(6705))

	// TODO Start an rpc server that serves crackers
}
