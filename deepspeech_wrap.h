#ifdef __cplusplus
extern "C" {
#endif
    typedef struct MetadataItem {
        char* character;
        int timestep;
        float start_time;
    } MetadataItem;

    typedef struct Metadata {
        MetadataItem* items;
        int num_items;
        double probability;
    } Metadata;

    typedef void* ModelWrapper;
    ModelWrapper* New(const char* aModelPath, int aNCep, int aNContext, const char* aAlphabetConfigPath, int aBeamWidth);
    void Close(ModelWrapper* w);
    void EnableDecoderWithLM(ModelWrapper* w, const char* aAlphabetConfigPath, const char* aLMPath, const char* aTriePath, float aLMWeight, float aValidWordCountWeight);
    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, unsigned int aSampleRate);
    Metadata* STTWithMetadata(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, unsigned int aSampleRate);

    typedef void* StreamWrapper;
    StreamWrapper* SetupStream(ModelWrapper* w, unsigned int aPreAllocFrames, unsigned int aSampleRate);
    void DiscardStream(StreamWrapper* sw);
    void FeedAudioContent(StreamWrapper* sw, const short* aBuffer, unsigned int aBufferSize);
    char* IntermediateDecode(StreamWrapper* sw);
    char* FinishStream(StreamWrapper* sw);
    Metadata* FinishStreamWithMetadata(StreamWrapper* sw);

    MetadataItem* Metadata_GetItems(Metadata* m);
    double Metadata_GetProbability(Metadata* m);
    int Metadata_GetNumItems(Metadata* m);

    char* MetadataItem_GetCharacter(MetadataItem* mi);
    int MetadataItem_GetTimestep(MetadataItem* mi);
    float MetadataItem_GetStartTime(MetadataItem* mi);

    void FreeString(char* s);
    void FreeMetadata(Metadata* m);
    void PrintVersions();

#ifdef __cplusplus
}
#endif
