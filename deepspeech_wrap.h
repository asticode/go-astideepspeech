#ifdef __cplusplus
extern "C" {
#endif
    typedef struct TokenMetadata {
        const char* text;
        const unsigned int timestep;
        const float start_time;
    } TokenMetadata;

    typedef struct CandidateTranscript {
        const TokenMetadata* const tokens;
        const unsigned int num_tokens;
        const double confidence;
    } CandidateTranscript;

    typedef struct Metadata {
        const CandidateTranscript* const transcripts;
        const unsigned int num_transcripts;
    } Metadata;

    typedef void* ModelWrapper;
    ModelWrapper* New(const char* aModelPath, int* errorOut);
    void Model_Close(ModelWrapper* w);
    unsigned int Model_BeamWidth(ModelWrapper* w);
    int Model_SetBeamWidth(ModelWrapper* w, unsigned int aBeamWidth);
    int Model_SampleRate(ModelWrapper* w);
    int Model_EnableExternalScorer(ModelWrapper* w, const char* aScorerPath);
    int Model_DisableExternalScorer(ModelWrapper* w);
    int Model_SetScorerAlphaBeta(ModelWrapper* w, float aAlpha, float aBeta);
    char* Model_STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize);
    Metadata* Model_STTWithMetadata(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, unsigned int aNumResults);

    typedef void* StreamWrapper;
    StreamWrapper* Model_NewStream(ModelWrapper* w, int* errorOut);
    void Stream_Discard(StreamWrapper* sw);
    void Stream_FeedAudioContent(StreamWrapper* sw, const short* aBuffer, unsigned int aBufferSize);
    char* Stream_IntermediateDecode(StreamWrapper* sw);
    Metadata* Stream_IntermediateDecodeWithMetadata(StreamWrapper* sw, unsigned int aNumResults);
    char* Stream_Finish(StreamWrapper* sw);
    Metadata* Stream_FinishWithMetadata(StreamWrapper* sw, unsigned int aNumResults);

    const CandidateTranscript* Metadata_Transcripts(Metadata* m);
    unsigned int Metadata_NumTranscripts(Metadata* m);
    void Metadata_Close(Metadata* m);

    const TokenMetadata* CandidateTranscript_Tokens(CandidateTranscript* ct);
    unsigned int CandidateTranscript_NumTokens(CandidateTranscript* ct);
    double CandidateTranscript_Confidence(CandidateTranscript* ct);

    const char* TokenMetadata_Text(TokenMetadata* tm);
    unsigned int TokenMetadata_Timestep(TokenMetadata* tm);
    float TokenMetadata_StartTime(TokenMetadata* tm);

    void FreeString(char* s);
    char* Version();
    char* ErrorCodeToErrorMessage(int aErrorCode);

#ifdef __cplusplus
}
#endif
