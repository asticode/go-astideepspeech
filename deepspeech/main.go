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

// Constants
const (
	beamWidth            = 500
	lmWeight             = 0.75
	validWordCountWeight = 1.85
)

var model = flag.String("model", "", "Path to the model (protocol buffer binary file)")
var audio = flag.String("audio", "", "Path to the audio file to run (WAV format)")
var lm = flag.String("lm", "", "Path to the language model binary file")
var trie = flag.String("trie", "", "Path to the language model trie file created with native_client/generate_trie")
var version = flag.Bool("version", false, "Print version and exits")
var extended = flag.Bool("extended", false, "Use extended metadata")

func metadataToString(m *astideepspeech.Metadata) string {
	retval := ""
	for _, item := range m.Items() {
		retval += item.Character()
	}
	return retval
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *version {
		astideepspeech.PrintVersions()
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
	m := astideepspeech.New(*model, beamWidth)
	defer m.Close()
	if *lm != "" {
		m.EnableDecoderWithLM(*lm, *trie, lmWeight, validWordCountWeight)
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

	output := ""
	// Speech to text
	if *extended {
		metadata := m.SpeechToTextWithMetadata(d, uint(len(d)))
		defer metadata.Close()
		output = metadataToString(metadata)
	} else {
		output = m.SpeechToText(d, uint(len(d)))
	}

	log.Printf("Text: %s\n", output)
}
