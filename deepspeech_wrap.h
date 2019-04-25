#ifdef __cplusplus
extern "C" {
#endif

    typedef void* ModelWrapper;
    ModelWrapper* New(const char* aModelPath, int aNCep, int aNContext, const char* aAlphabetConfigPath, int aBeamWidth);
    void Close(ModelWrapper* w);
    void EnableDecoderWithLM(ModelWrapper* w, const char* aAlphabetConfigPath, const char* aLMPath, const char* aTriePath, float aLMWeight, float aValidWordCountWeight);
    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, unsigned int aSampleRate);

    typedef void* StreamWrapper;
    StreamWrapper* SetupStream(ModelWrapper* w, unsigned int aPreAllocFrames, unsigned int aSampleRate);
    void DiscardStream(StreamWrapper* sw);
    void FeedAudioContent(StreamWrapper* sw, const short* aBuffer, unsigned int aBufferSize);
    char* IntermediateDecode(StreamWrapper* sw);
    char* FinishStream(StreamWrapper* sw);

    void FreeString(char* s);
    void PrintVersions();

#ifdef __cplusplus
}
#endif
