package astideepspeech

/*
#cgo LDFLAGS: -ldeepspeech
#include "deepspeech_wrap.h"
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
// validWordCountWeight The weight (bonus) to give to beams when adding a new valid word to the decoding.
func (m *Model) EnableDecoderWithLM(alphabetConfigPath, lmPath, triePath string, lmWeight, validWordCountWeight float64) {
	C.EnableDecoderWithLM(m.w, C.CString(alphabetConfigPath), C.CString(lmPath), C.CString(triePath), C.float(lmWeight), C.float(validWordCountWeight))
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
func (m *Model) SpeechToText(buffer []int16, bufferSize, sampleRate uint) string {
	return C.GoString(C.STT(m.w, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize), C.uint(sampleRate)))
}

// Stream represent a streaming state
type Stream struct {
        sw	*C.StreamWrapper
}

// SetupStream creates a new audio stream
//
// mw               The DeepSpeech model to use
// preAllocFrames   Number of timestep frames to reserve. One timestep
//                  is equivalent to two window lengths (20ms). If set to
//                  0 we reserve enough frames for 3 seconds of audio (150).
// aSampleRate      The sample-rate of the audio signal.
func SetupStream(mw *Model, preAllocFrames uint, sampleRate uint) *Stream {
	return &Stream{
		sw:	C.SetupStream(mw.w, C.uint(preAllocFrames), C.uint(sampleRate)),
	}
}


// FeedAudioContent Feed audio samples to an ongoing streaming inference.
// aBuffer      An array of 16-bit, mono raw audio samples at the  appropriate sample rate.
// aBufferSize  The number of samples in @p aBuffer.
func (s *Stream) FeedAudioContent(buffer []int16, bufferSize uint) {
	C.FeedAudioContent(s.sw, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize))
}

// IntermediateDecode Compute the intermediate decoding of an ongoing streaming inference.
// This is an expensive process as the decoder implementation isn't
// currently capable of streaming, so it always starts from the beginning
// of the audio.
func (s *Stream) IntermediateDecode() string {
	return C.GoString(C.IntermediateDecode(s.sw))
}

// FinishStream Signal the end of an audio signal to an ongoing streaming
// inference, returns the STT result over the whole audio signal.
func (s *Stream) FinishStream() string {
	return C.GoString(C.FinishStream(s.sw))
}

// Destroy a streaming state without decoding the computed logits. This
// can be used if you no longer need the result of an ongoing streaming
// inference and don't want to perform a costly decode operation.
func (s *Stream) DiscardStream() {
        C.DiscardStream(s.sw);
}

// PrintVersions Print version of this library and of the linked TensorFlow library.
func PrintVersions() {
        C.PrintVersions()
}
