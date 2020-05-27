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
    ModelWrapper* New(const char* aModelPath);
    void Close(ModelWrapper* w);
    unsigned int GetModelBeamWidth(ModelWrapper* w);
    int SetModelBeamWidth(ModelWrapper* w, unsigned int aBeamWidth);
    int GetModelSampleRate(ModelWrapper* w);
    int EnableExternalScorer(ModelWrapper* w, const char* aScorerPath);
    int DisableExternalScorer(ModelWrapper* w);
    int SetScorerAlphaBeta(ModelWrapper* w, float aAlpha, float aBeta);
    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize);
    Metadata* STTWithMetadata(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, unsigned int aNumResults);

    typedef void* StreamWrapper;
    StreamWrapper* CreateStream(ModelWrapper* w);
    void FreeStream(StreamWrapper* sw);
    void FeedAudioContent(StreamWrapper* sw, const short* aBuffer, unsigned int aBufferSize);
    char* IntermediateDecode(StreamWrapper* sw);
    Metadata* IntermediateDecodeWithMetadata(StreamWrapper* sw, unsigned int aNumResults);
    char* FinishStream(StreamWrapper* sw);
    Metadata* FinishStreamWithMetadata(StreamWrapper* sw, unsigned int aNumResults);

    const CandidateTranscript* Metadata_GetTranscripts(Metadata* m);
    unsigned int Metadata_GetNumTranscripts(Metadata* m);

    const TokenMetadata* CandidateTranscript_GetTokens(CandidateTranscript* ct);
    unsigned int CandidateTranscript_GetNumTokens(CandidateTranscript* ct);
    double CandidateTranscript_GetConfidence(CandidateTranscript* ct);

    const char* TokenMetadata_GetText(TokenMetadata* tm);
    unsigned int TokenMetadata_GetTimestep(TokenMetadata* tm);
    float TokenMetadata_GetStartTime(TokenMetadata* tm);

    void FreeString(char* s);
    void FreeMetadata(Metadata* m);
    char* Version();
    char* ErrorCodeToErrorMessage(int aErrorCode);

#ifdef __cplusplus
}
#endif
