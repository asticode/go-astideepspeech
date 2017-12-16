package astideepspeech

/*
#cgo LDFLAGS: -ldeepspeech -ldeepspeech_utils -ltensorflow_cc -ltensorflow_framework
#include "deepspeech.h"
*/
import "C"
import "unsafe"

// Model represents a DeepSpeech model
type Model struct {
	alphabetConfigPath string
	beamWidth          int
	modelPath          string
	nCep               int
	nContext           int
	w                  *C.ModelWrapper
}

// New creates a new Model
//
// modelPath          The path to the frozen model graph.
// nCep               The number of cepstrum the model was trained with.
// nContext           The context window the model was trained with.
// alphabetConfigPath The path to the configuration file specifying the alphabet used by the network.
// beamWidth          The beam width used by the decoder. A larger beam width generates better results at the cost of decoding time.
func New(modelPath string, nCep, nContext int, alphabetConfigPath string, beamWidth int) *Model {
	return &Model{
		alphabetConfigPath: alphabetConfigPath,
		beamWidth:          beamWidth,
		modelPath:          modelPath,
		nCep:               nCep,
		nContext:           nContext,
		w:                  C.New(C.CString(modelPath), C.int(nCep), C.int(nContext), C.CString(alphabetConfigPath), C.int(beamWidth)),
	}
}

// Close closes the model properly
func (m *Model) Close() error {
	C.Close(m.w)
	return nil
}

// EnableDecoderWithLM enables decoding using beam scoring with a KenLM language model.
//
// alphabetConfigPath   The path to the configuration file specifying the alphabet used by the network.
// lmPath 	        The path to the language model binary file.
// triePath 	        The path to the trie file build from the same vocabulary as the language model binary.
// lmWeight 	        The weight to give to language model results when scoring.
// wordCountWeight      The weight (penalty) to give to beams when increasing the word count of the decoding.
// validWordCountWeight The weight (bonus) to give to beams when adding a new valid word to the decoding.
func (m *Model) EnableDecoderWithLM(alphabetConfigPath, lmPath, triePath string, lmWeight, wordCountWeight, validWordCountWeight float64) {
	C.EnableDecoderWithLM(m.w, C.CString(alphabetConfigPath), C.CString(lmPath), C.CString(triePath), C.float(lmWeight), C.float(wordCountWeight), C.float(validWordCountWeight))
}

// sliceHeader represents a slice header
type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

// SpeechToText uses the DeepSpeech model to perform Speech-To-Text.
// buffer     A 16-bit, mono raw audio signal at the appropriate sample rate.
// bufferSize The number of samples in the audio signal.
// sampleRate The sample-rate of the audio signal.
// TODO Make sure the C string is cleaned properly
func (m *Model) SpeechToText(buffer []int16, bufferSize, sampleRate int) string {
	return C.GoString(C.STT(m.w, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize), C.int(sampleRate)))
}
