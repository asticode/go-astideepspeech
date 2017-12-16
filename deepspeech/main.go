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
	lmWeight             = 1.75
	wordCountWeight      = 1.00
	validWordCountWeight = 1.00
)

func main() {
	// Parse flags
	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, `Usage: deepspeech MODEL_PATH AUDIO_PATH ALPHABET_PATH [LM_PATH] [TRIE_PATH]
  MODEL_PATH:          Path to the model (protocol buffer binary file)
  AUDIO_PATH:          Path to the audio file to run (must be a .wav file)"
  ALPHABET_PATH:       Path to the configuration file specifying the alphabet used by the network."
  LM_PATH(Optional):   Path to the language model binary file.
  TRIE_PATH(Optional): Path to the language model trie file created with native_client/generate_trie.`)
	}
	flag.Parse()
	astilog.FlagInit()

	// Invalid number of args
	if len(os.Args) < 4 || len(os.Args) > 7 {
		flag.Usage()
		return
	}

	// Initialize DeepSpeech
	m := astideepspeech.New(os.Args[1], nCep, nContext, os.Args[3], beamWidth)
	defer m.Close()
	if len(os.Args) > 5 {
		m.EnableDecoderWithLM(os.Args[3], os.Args[4], os.Args[5], lmWeight, wordCountWeight, validWordCountWeight)
	}

	// Stat audio
	i, err := os.Stat(os.Args[2])
	if err != nil {
		astilog.Fatal(errors.Wrapf(err, "stating %s failed", os.Args[2]))
	}

	// Open audio
	f, err := os.Open(os.Args[2])
	if err != nil {
		astilog.Fatal(errors.Wrapf(err, "opening %s failed", os.Args[2]))
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
	astilog.Infof("Text: %s", m.SpeechToText(d, len(d), 16000))
}
