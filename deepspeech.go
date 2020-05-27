package astideepspeech

/*
#cgo LDFLAGS: -ldeepspeech
#include "deepspeech_wrap.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Model provides an interface to a trained DeepSpeech model.
type Model struct {
	modelPath string
	w         *C.ModelWrapper
}

// New creates a new Model.
// modelPath is the path to the frozen model graph.
func New(modelPath string) *Model {
	return &Model{
		modelPath: modelPath,
		w:         C.New(C.CString(modelPath)),
	}
}

// Close frees associated resources and destroys the model object.
func (m *Model) Close() error {
	C.Close(m.w)
	return nil
}

// GetModelBeamWidth returns the beam width value used by the model.
// If SetModelBeamWidth was not called before, will return the default
// value loaded from the model file.
func (m *Model) GetModelBeamWidth() uint {
	return uint(C.GetModelBeamWidth(m.w))
}

// SetModelBeamWidth sets the beam width value used by the model.
// A larger beam width value generates better results at the cost of decoding time.
func (m *Model) SetModelBeamWidth(width uint) error {
	return errorFromCode(C.SetModelBeamWidth(m.w, C.uint(width)))
}

// GetModelSampleRate returns the sample rate that was used to produce the model file.
func (m *Model) GetModelSampleRate() int {
	return int(C.GetModelSampleRate(m.w))
}

// EnableExternalScorer enables decoding using an external scorer.
// scorerPath is the path to the external scorer file.
func (m *Model) EnableExternalScorer(scorerPath string) error {
	return errorFromCode(C.EnableExternalScorer(m.w, C.CString(scorerPath)))
}

// DisableExternalScorer disables decoding using an external scorer.
func (m *Model) DisableExternalScorer() error {
	return errorFromCode(C.DisableExternalScorer(m.w))
}

// SetScorerAlphaBeta sets hyperparameters alpha and beta of the external scorer.
// alpha is the language model weight. beta is the word insertion weight.
func (m *Model) SetScorerAlphaBeta(alpha, beta float32) error {
	return errorFromCode(C.SetScorerAlphaBeta(m.w, C.float(alpha), C.float(beta)))
}

// sliceHeader represents a slice header
type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

// SpeechToText uses the DeepSpeech model to convert speech to text.
// buffer is 16-bit, mono raw audio signal at the appropriate sample rate (matching what the model was trained on).
// bufferSize is the number of samples in the audio signal.
func (m *Model) SpeechToText(buffer []int16, bufferSize uint) string {
	str := C.STT(m.w, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize))
	defer C.FreeString(str)
	retval := C.GoString(str)
	return retval
}

// TokenMetadata stores text of an individual token, along with its timing information.
type TokenMetadata C.struct_TokenMetadata

// Text returns the text corresponding to this token.
func (tm *TokenMetadata) Text() string {
	return C.GoString(C.TokenMetadata_GetText((*C.TokenMetadata)(unsafe.Pointer(tm))))
}

// Timestep returns the position of the token in units of 20ms.
func (tm *TokenMetadata) Timestep() uint {
	return uint(C.TokenMetadata_GetTimestep((*C.TokenMetadata)(unsafe.Pointer(tm))))
}

// StartTime returns the position of the token in seconds.
func (tm *TokenMetadata) StartTime() float32 {
	return float32(C.TokenMetadata_GetStartTime((*C.TokenMetadata)(unsafe.Pointer(tm))))
}

// CandidateTranscript is a single transcript computed by the model,
// including a confidence value and the metadata for its constituent tokens.
type CandidateTranscript C.struct_CandidateTranscript

func (ct *CandidateTranscript) NumTokens() uint {
	return uint(C.CandidateTranscript_GetNumTokens((*C.CandidateTranscript)(unsafe.Pointer(ct))))
}

func (ct *CandidateTranscript) Tokens() []TokenMetadata {
	numTokens := uint(C.CandidateTranscript_GetNumTokens((*C.CandidateTranscript)(unsafe.Pointer(ct))))
	allTokens := C.CandidateTranscript_GetTokens((*C.CandidateTranscript)(unsafe.Pointer(ct)))
	return (*[1 << 30]TokenMetadata)(unsafe.Pointer(allTokens))[:numTokens:numTokens]
}

// Confidence returns the approximated confidence value for this transcript.
// This is roughly the sum of the acoustic model logit values for each timestep/character that
// contributed to the creation of this transcript.
func (ct *CandidateTranscript) Confidence() float64 {
	return float64(C.CandidateTranscript_GetConfidence((*C.CandidateTranscript)(unsafe.Pointer(ct))))
}

// Metadata holds an array of CandidateTranscript objects computed by the model.
type Metadata C.struct_Metadata

func (m *Metadata) NumTranscripts() uint {
	return uint(C.Metadata_GetNumTranscripts((*C.Metadata)(unsafe.Pointer(m))))
}

func (m *Metadata) Transcripts() []CandidateTranscript {
	numTranscripts := int32(C.Metadata_GetNumTranscripts((*C.Metadata)(unsafe.Pointer(m))))
	allTranscripts := C.Metadata_GetTranscripts((*C.Metadata)(unsafe.Pointer(m)))
	return (*[1 << 30]CandidateTranscript)(unsafe.Pointer(allTranscripts))[:numTranscripts:numTranscripts]
}

// Close frees the Metadata structure properly.
func (m *Metadata) Close() error {
	C.FreeMetadata((*C.Metadata)(unsafe.Pointer(m)))
	return nil
}

// SpeechToTextWithMetadata uses the DeepSpeech model to convert speech to text and
// output results including metadata.
// buffer is a 16-bit, mono raw audio signal at the appropriate sample rate (matching what the model was trained on).
// bufferSize is the number of samples in the audio signal.
// numResults is the maximum number of CandidateTranscript structs to return. Returned value might be smaller than this.
func (m *Model) SpeechToTextWithMetadata(buffer []int16, bufferSize, numResults uint) *Metadata {
	return (*Metadata)(unsafe.Pointer(C.STTWithMetadata(m.w, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize), C.uint(numResults))))
}

// Stream represent a streaming state
type Stream struct {
	sw *C.StreamWrapper
}

// CreateStream creates a new audio stream
//
// mw               The DeepSpeech model to use
func CreateStream(mw *Model) *Stream {
	return &Stream{
		sw: C.CreateStream(mw.w),
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

// IntermediateDecodeWithMetadata computes the intermediate decoding of an
// ongoing streaming inference, returning results including metadata.
// numResults is the number of candidate transcripts to return.
func (s *Stream) IntermediateDecodeWithMetadata(numResults uint) *Metadata {
	return (*Metadata)(unsafe.Pointer(C.IntermediateDecodeWithMetadata(s.sw, C.uint(numResults))))
}

// FinishStream Signal the end of an audio signal to an ongoing streaming
// inference, returns the STT result over the whole audio signal.
func (s *Stream) FinishStream() string {
	str := C.FinishStream(s.sw)
	defer C.FreeString(str)
	retval := C.GoString(str)
	return retval
}

// FinishStreamWithMetadata Signal the end of an audio signal to an ongoing streaming
// inference, returns extended metadata.
func (s *Stream) FinishStreamWithMetadata(numResults uint) *Metadata {
	return (*Metadata)(unsafe.Pointer(C.FinishStreamWithMetadata(s.sw, C.uint(numResults))))
}

// DiscardStream Destroy a streaming state without decoding the computed logits.
// This can be used if you no longer need the result of an ongoing streaming
// inference and don't want to perform a costly decode operation.
func (s *Stream) FreeStream() {
	C.FreeStream(s.sw)
}

// Version returns the version of the DeepSpeech C library.
// The returned version is a semantic version (SemVer 2.0.0).
func Version() string {
	str := C.Version()
	defer C.FreeString(str)
	return C.GoString(str)
}

// errorFromCode converts an error code returned by DeepSpeech into a Go error.
// Returns nil if code is equal to zero, indicating success.
func errorFromCode(code C.int) error {
	if code == 0 {
		return nil
	}
	str := C.ErrorCodeToErrorMessage(code)
	defer C.FreeString(str)
	return errors.New(C.GoString(str))
}
