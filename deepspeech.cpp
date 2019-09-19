#include <stdio.h>
#include <deepspeech.h>

extern "C" {
    class ModelWrapper {
        private:
            ModelState* model;

        public:
            ModelWrapper(const char* aModelPath, int aBeamWidth)
            {
                DS_CreateModel(aModelPath, aBeamWidth, &model);
            }

            ~ModelWrapper()
            {
                DS_FreeModel(model);
            }

            void enableDecoderWithLM(const char* aLMPath, const char* aTriePath, float aLMWeight, float aValidWordCountWeight)
            {
                DS_EnableDecoderWithLM(model, aLMPath, aTriePath, aLMWeight, aValidWordCountWeight);
            }

	    int getModelSampleRate()
	    {
		return DS_GetModelSampleRate(model);
	    }

            char* stt(const short* aBuffer, unsigned int aBufferSize)
            {
                return DS_SpeechToText(model, aBuffer, aBufferSize);
            }

            Metadata* sttWithMetadata(const short* aBuffer, unsigned int aBufferSize)
            {
                return DS_SpeechToTextWithMetadata(model, aBuffer, aBufferSize);
            }

            ModelState* getModel()
            {
                return model;
            }
    };

    ModelWrapper* New(const char* aModelPath, int aBeamWidth)
    {
        return new ModelWrapper(aModelPath, aBeamWidth);
    }
    void Close(ModelWrapper* w)
    {
        delete w;
    }

    void EnableDecoderWithLM(ModelWrapper* w, const char* aLMPath, const char* aTriePath, float aLMWeight, float aValidWordCountWeight)
    {
        w->enableDecoderWithLM(aLMPath, aTriePath, aLMWeight, aValidWordCountWeight);
    }

    int GetModelSampleRate(ModelWrapper* w)
    {
	return w->getModelSampleRate();
    }

    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize)
    {
        return w->stt(aBuffer, aBufferSize);
    }

    Metadata* STTWithMetadata(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize)
    {
        return w->sttWithMetadata(aBuffer, aBufferSize);
    }

    double Metadata_GetConfidence(Metadata* m)
    {
        return m->confidence;
    }

    int Metadata_GetNumItems(Metadata* m)
    {
        return m->num_items;
    }

    MetadataItem* Metadata_GetItems(Metadata* m)
    {
        return m->items;
    }

    char* MetadataItem_GetCharacter(MetadataItem* mi)
    {
        return mi->character;
    }

    int MetadataItem_GetTimestep(MetadataItem* mi)
    {
        return mi->timestep;
    }

    float MetadataItem_GetStartTime(MetadataItem* mi)
    {
        return mi->start_time;
    }

    class StreamWrapper {
        private:
            StreamingState* s;

        public:
            StreamWrapper(ModelWrapper* w)
            {
                DS_CreateStream(w->getModel(), &s);
            }

            ~StreamWrapper()
            {
                DS_FreeStream(s);
            }

            void feedAudioContent(const short* aBuffer, unsigned int aBufferSize)
            {
                DS_FeedAudioContent(s, aBuffer, aBufferSize);
            }

            char* intermediateDecode()
            {
                return DS_IntermediateDecode(s);
            }

            char* finishStream()
            {
                return DS_FinishStream(s);
            }

            Metadata* finishStreamWithMetadata()
            {
                return DS_FinishStreamWithMetadata(s);
            }

            void freeStream()
            {
                DS_FreeStream(s);
            }
    };

    StreamWrapper* CreateStream(ModelWrapper* mw, unsigned int aPreAllocFrames)
    {
        return new StreamWrapper(mw);
    }
    void FreeStream(StreamWrapper* sw)
    {
        delete sw;
    }

    void FeedAudioContent(StreamWrapper* sw, const short* aBuffer, unsigned int aBufferSize)
    {
        sw->feedAudioContent(aBuffer, aBufferSize);
    }

    char* IntermediateDecode(StreamWrapper* sw)
    {
        return sw->intermediateDecode();
    }

    char* FinishStream(StreamWrapper* sw)
    {
        return sw->finishStream();
    }

    Metadata* FinishStreamWithMetadata(StreamWrapper* sw)
    {
        return sw->finishStreamWithMetadata();
    }

    void FreeString(char* s)
    {
        DS_FreeString(s);
    }

    void FreeMetadata(Metadata* m)
    {
        DS_FreeMetadata(m);
    }

    void PrintVersions()
    {
        DS_PrintVersions();
    }
}
