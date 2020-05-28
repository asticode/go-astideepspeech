// package astideepspeech provides bindings for Mozilla's DeepSpeech speech-to-text library.
package astideepspeech

/*
#cgo CXXFLAGS: -std=c++11
#cgo LDFLAGS: -ldeepspeech
#include "deepspeech_wrap.h"
#include "stdlib.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Model provides an interface to a trained DeepSpeech model.
type Model struct {
	w *C.ModelWrapper
}

// New creates a new Model.
// modelPath is the path to the frozen model graph.
func New(modelPath string) (*Model, error) {
	cModelPath := C.CString(modelPath)
	defer C.free(unsafe.Pointer(cModelPath))

	var ret C.int
	w := C.New(cModelPath, &ret)
	if ret != 0 {
		C.Close(w)
		return nil, errorFromCode(ret)
	}
	return &Model{w}, nil
}

// Close frees associated resources and destroys the model object.
func (m *Model) Close() {
	C.Close(m.w)
}

// GetModelBeamWidth returns the beam width value used by the model.
// If SetModelBeamWidth was not called before, it will return the default
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
	cScorerPath := C.CString(scorerPath)
	defer C.free(unsafe.Pointer(cScorerPath))
	return errorFromCode(C.EnableExternalScorer(m.w, cScorerPath))
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
func (m *Model) SpeechToText(buffer []int16) (string, error) {
	hdr := (*sliceHeader)(unsafe.Pointer(&buffer))
	str := C.STT(m.w, (*C.short)(unsafe.Pointer(hdr.Data)), C.uint(hdr.Len))
	if str == nil {
		return "", errors.New("conversion failed")
	}
	defer C.FreeString(str)
	return C.GoString(str), nil
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
func (m *Metadata) Close() {
	C.FreeMetadata((*C.Metadata)(unsafe.Pointer(m)))
}

// SpeechToTextWithMetadata uses the DeepSpeech model to convert speech to text and
// output results including metadata.
//
// buffer is a 16-bit, mono raw audio signal at the appropriate sample rate (matching what the model was trained on).
// numResults is the maximum number of CandidateTranscript structs to return. Returned value might be smaller than this.
// If an error is not returned, the returned metadata's Close method must be called later to free resources.
func (m *Model) SpeechToTextWithMetadata(buffer []int16, numResults uint) (*Metadata, error) {
	hdr := (*sliceHeader)(unsafe.Pointer(&buffer))
	md := (*Metadata)(unsafe.Pointer(C.STTWithMetadata(
		m.w, (*C.short)(unsafe.Pointer(hdr.Data)), C.uint(hdr.Len), C.uint(numResults))))
	if md == nil {
		return nil, errors.New("conversion failed")
	}
	return md, nil
}

// Stream represents a streaming inference state.
type Stream struct {
	sw *C.StreamWrapper
}

// CreateStream creates a new streaming inference state.
// m is the DeepSpeech model to use.
// If an error is not returned, exactly one of the returned stream's FinishStream,
// FinishStreamWithMetadata, or FreeStream methods must be called later to free resources.
func CreateStream(m *Model) (*Stream, error) {
	var ret C.int
	sw := C.CreateStream(m.w, &ret)
	if ret != 0 {
		C.FreeStream(sw)
		return nil, errorFromCode(ret)
	}
	return &Stream{sw}, nil
}

// FeedAudioContent feeds audio samples to an ongoing streaming inference.
// buffer is an array of 16-bit, mono raw audio samples at the appropriate sample rate
// (matching what the model was trained on).
func (s *Stream) FeedAudioContent(buffer []int16) {
	hdr := (*sliceHeader)(unsafe.Pointer(&buffer))
	C.FeedAudioContent(s.sw, (*C.short)(unsafe.Pointer(hdr.Data)), C.uint(hdr.Len))
}

// IntermediateDecode computes the intermediate decoding of an ongoing streaming inference.
// This is an expensive process as the decoder implementation isn't
// currently capable of streaming, so it always starts from the beginning
// of the audio.
func (s *Stream) IntermediateDecode() (string, error) {
	// DS_IntermediateDecode isn't documented as returning null, but future-proofing this seems safer.
	str := C.IntermediateDecode(s.sw)
	if str == nil {
		return "", errors.New("decoding failed")
	}
	defer C.FreeString(str)
	return C.GoString(str), nil
}

// IntermediateDecodeWithMetadata computes the intermediate decoding of an
// ongoing streaming inference, returning results including metadata.
// numResults is the number of candidate transcripts to return.
// If an error is not returned, the metadata's Close method must be called.
func (s *Stream) IntermediateDecodeWithMetadata(numResults uint) (*Metadata, error) {
	md := (*Metadata)(unsafe.Pointer(C.IntermediateDecodeWithMetadata(s.sw, C.uint(numResults))))
	if md == nil {
		return nil, errors.New("decoding failed")
	}
	return md, nil
}

// FinishStream computes the final decoding of an ongoing streaming inference and returns the result.
// This signals the end of an ongoing streaming inference.
func (s *Stream) FinishStream() (string, error) {
	// DS_FinishStream isn't documented as returning null, but future-proofing this seems safer.
	str := C.FinishStream(s.sw)
	if str == nil {
		return "", errors.New("decoding failed")
	}
	defer C.FreeString(str)
	return C.GoString(str), nil
}

// FinishStreamWithMetadata computes the final decoding of an ongoing streaming inference and returns
// results including metadata. This signals the end of an ongoing streaming inference.
// If an error is not returned, the metadata's Close method must be called.
func (s *Stream) FinishStreamWithMetadata(numResults uint) (*Metadata, error) {
	md := (*Metadata)(unsafe.Pointer(C.FinishStreamWithMetadata(s.sw, C.uint(numResults))))
	if md == nil {
		return nil, errors.New("decoding failed")
	}
	return md, nil
}

// FreeStream destroys a streaming state without decoding the computed logits.
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
