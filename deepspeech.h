#ifdef __cplusplus
extern "C" {
#endif

    typedef void* ModelWrapper;
    ModelWrapper* New(const char* aModelPath, int aNCep, int aNContext, const char* aAlphabetConfigPath, int aBeamWidth);
    void Close(ModelWrapper* w);
    void EnableDecoderWithLM(ModelWrapper* w, const char* aAlphabetConfigPath, const char* aLMPath, const char* aTriePath, float aLMWeight, float aWordCountWeight, float aValidWordCountWeight);
    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, int aSampleRate);

#ifdef __cplusplus
}
#endif