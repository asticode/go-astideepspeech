package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/asticode/go-astideepspeech"
	"github.com/cryptix/wav"
)

const (
	beamWidth = 500

	// These are optimal hyperparameter values found with respect to the LibriSpeech
	// clean dev corpus per https://github.com/mozilla/DeepSpeech/releases/tag/v0.7.1.
	alpha = 0.931289039105002
	beta  = 1.1834137581510284
)

var model = flag.String("model", "", "Path to the model (protocol buffer binary file)")
var audio = flag.String("audio", "", "Path to the audio file to run (WAV format)")
var scorer = flag.String("scorer", "", "Path to the external scorer")
var version = flag.Bool("version", false, "Print version and exits")
var extended = flag.Bool("extended", false, "Use extended metadata")
var maxResults = flag.Uint("max-results", 5, "Maximum number of results when -extended is true")

func metadataToStrings(m *astideepspeech.Metadata) []string {
	results := make([]string, 0, m.NumTranscripts())
	for _, tr := range m.Transcripts() {
		var res string
		for _, tok := range tr.Tokens() {
			res += tok.Text()
		}
		res += fmt.Sprintf(" [%0.3f]", tr.Confidence())
		results = append(results, res)
	}
	return results
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *version {
		println(astideepspeech.Version())
		return
	}

	if *model == "" || *audio == "" {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	// Initialize DeepSpeech
	m := astideepspeech.New(*model)
	defer m.Close()

	if err := m.SetModelBeamWidth(beamWidth); err != nil {
		log.Fatal("Failed setting beam width: ", err)
	}
	if *scorer != "" {
		if err := m.EnableExternalScorer(*scorer); err != nil {
			log.Fatal("Failed enabling external scorer: ", err)
		}
		if err := m.SetScorerAlphaBeta(alpha, beta); err != nil {
			log.Fatal("Failed setting scorer hyperparameters: ", err)
		}
	}

	// Stat audio
	i, err := os.Stat(*audio)
	if err != nil {
		log.Fatal(fmt.Errorf("stating %s failed: %w", *audio, err))
	}

	// Open audio
	f, err := os.Open(*audio)
	if err != nil {
		log.Fatal(fmt.Errorf("opening %s failed: %w", *audio, err))
	}

	// Create reader
	r, err := wav.NewReader(f, i.Size())
	if err != nil {
		log.Fatal(fmt.Errorf("creating new reader failed: %w", err))
	}

	// Read
	var d []int16
	for {
		// Read sample
		s, err := r.ReadSample()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(fmt.Errorf("reading sample failed: %w", err))
		}

		// Append
		d = append(d, int16(s))
	}

	// Speech to text
	var results []string
	if *extended {
		metadata := m.SpeechToTextWithMetadata(d, uint(len(d)), *maxResults)
		defer metadata.Close()
		results = metadataToStrings(metadata)
	} else {
		results = []string{m.SpeechToText(d, uint(len(d)))}
	}
	for _, res := range results {
		log.Println("Text: ", res)
	}

}
