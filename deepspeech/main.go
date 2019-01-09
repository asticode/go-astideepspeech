package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/asticode/go-astideepspeech"
	"github.com/asticode/go-astilog"
	"github.com/cryptix/wav"
	"github.com/pkg/errors"
)

// Constants
const (
	beamWidth            = 500
	nCep                 = 26
	nContext             = 9
	lmWeight             = 0.75
	validWordCountWeight = 1.85
)

var model    = flag.String("model", "",    "Path to the model (protocol buffer binary file)")
var alphabet = flag.String("alphabet", "", "Path to the configuration file specifying the alphabet used by the network")
var audio    = flag.String("audio", "",    "Path to the audio file to run (WAV format)")
var lm       = flag.String("lm", "",       "Path to the language model binary file")
var trie     = flag.String("trie", "",     "Path to the language model trie file created with native_client/generate_trie")
var version  = flag.Bool("version", false, "Print version and exits")

func main() {
	flag.Parse()

	astilog.FlagInit()

	if *version {
		astideepspeech.PrintVersions()
                return
	}

	if *model == "" || *alphabet == "" || *audio == "" {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	// Initialize DeepSpeech
	m := astideepspeech.New(*model, nCep, nContext, *alphabet, beamWidth)
	defer m.Close()
	if *lm != "" {
		m.EnableDecoderWithLM(*alphabet, *lm, *trie, lmWeight, validWordCountWeight)
	}

	// Stat audio
	i, err := os.Stat(*audio)
	if err != nil {
		astilog.Fatal(errors.Wrapf(err, "stating %s failed", *audio))
	}

	// Open audio
	f, err := os.Open(*audio)
	if err != nil {
		astilog.Fatal(errors.Wrapf(err, "opening %s failed", audio))
	}

	// Create reader
	r, err := wav.NewReader(f, i.Size())
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "creating new reader failed"))
	}

	// Read
	var d []int16
	for {
		// Read sample
		s, err := r.ReadSample()
		if err == io.EOF {
			break
		} else if err != nil {
			astilog.Fatal(errors.Wrap(err, "reading sample failed"))
		}

		// Append
		d = append(d, int16(s))
	}

	// Speech to text
	astilog.Infof("Text: %s", m.SpeechToText(d, uint(len(d)), 16000))
}
