#include <stdio.h>
#include <deepspeech.h>

extern "C" {
    class ModelWrapper {
        private:
            ModelState* model;

        public:
            ModelWrapper(const char* aModelPath)
            {
                DS_CreateModel(aModelPath, &model);
            }

            ~ModelWrapper()
            {
                DS_FreeModel(model);
            }

            unsigned int getModelBeamWidth()
            {
                return DS_GetModelBeamWidth(model);
            }

            int setModelBeamWidth(unsigned int aBeamWidth)
            {
                return DS_SetModelBeamWidth(model, aBeamWidth);
            }

            int getModelSampleRate()
            {
                return DS_GetModelSampleRate(model);
            }

            int enableExternalScorer(const char* aScorerPath)
            {
                return DS_EnableExternalScorer(model, aScorerPath);
            }

            int disableExternalScorer()
            {
                return DS_DisableExternalScorer(model);
            }

            int setScorerAlphaBeta(float aAlpha, float aBeta)
            {
                return DS_SetScorerAlphaBeta(model, aAlpha, aBeta);
            }

            char* stt(const short* aBuffer, unsigned int aBufferSize)
            {
                return DS_SpeechToText(model, aBuffer, aBufferSize);
            }

            Metadata* sttWithMetadata(const short* aBuffer, unsigned int aBufferSize, unsigned int aNumResults)
            {
                return DS_SpeechToTextWithMetadata(model, aBuffer, aBufferSize, aNumResults);
            }

            ModelState* getModel()
            {
                return model;
            }
    };

    ModelWrapper* New(const char* aModelPath)
    {
        return new ModelWrapper(aModelPath);
    }
    void Close(ModelWrapper* w)
    {
        delete w;
    }

    unsigned int GetModelBeamWidth(ModelWrapper* w)
    {
        return w->getModelBeamWidth();
    }

    int SetModelBeamWidth(ModelWrapper* w, unsigned int aBeamWidth)
    {
        return w->setModelBeamWidth(aBeamWidth);
    }

    int GetModelSampleRate(ModelWrapper* w)
    {
        return w->getModelSampleRate();
    }

    int EnableExternalScorer(ModelWrapper* w, const char* aScorerPath)
    {
        return w->enableExternalScorer(aScorerPath);
    }

    int DisableExternalScorer(ModelWrapper* w)
    {
        return w->disableExternalScorer();
    }

    int SetScorerAlphaBeta(ModelWrapper* w, float aAlpha, float aBeta)
    {
        return w->setScorerAlphaBeta(aAlpha, aBeta);
    }

    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize)
    {
        return w->stt(aBuffer, aBufferSize);
    }

    Metadata* STTWithMetadata(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, unsigned int aNumResults)
    {
        return w->sttWithMetadata(aBuffer, aBufferSize, aNumResults);
    }

    const CandidateTranscript* Metadata_GetTranscripts(Metadata* m)
    {
        return m->transcripts;
    }

    unsigned int Metadata_GetNumTranscripts(Metadata* m)
    {
        return m->num_transcripts;
    }

    const TokenMetadata* CandidateTranscript_GetTokens(CandidateTranscript* ct)
    {
        return ct->tokens;
    }

    int CandidateTranscript_GetNumTokens(CandidateTranscript* ct)
    {
        return ct->num_tokens;
    }

    double CandidateTranscript_GetConfidence(CandidateTranscript* ct)
    {
        return ct->confidence;
    }

    const char* TokenMetadata_GetText(TokenMetadata* tm)
    {
        return tm->text;
    }

    unsigned int TokenMetadata_GetTimestep(TokenMetadata* tm)
    {
        return tm->timestep;
    }

    float TokenMetadata_GetStartTime(TokenMetadata* tm)
    {
        return tm->start_time;
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

            Metadata* intermediateDecodeWithMetadata(unsigned int aNumResults)
            {
                return DS_IntermediateDecodeWithMetadata(s, aNumResults);
            }

            char* finishStream()
            {
                return DS_FinishStream(s);
            }

            Metadata* finishStreamWithMetadata(unsigned int aNumResults)
            {
                return DS_FinishStreamWithMetadata(s, aNumResults);
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

    Metadata* IntermediateDecodeWithMetadata(StreamWrapper* sw, unsigned int aNumResults)
    {
        return sw->intermediateDecodeWithMetadata(aNumResults);
    }

    char* FinishStream(StreamWrapper* sw)
    {
        return sw->finishStream();
    }

    Metadata* FinishStreamWithMetadata(StreamWrapper* sw, unsigned int aNumResults)
    {
        return sw->finishStreamWithMetadata(aNumResults);
    }

    void FreeString(char* s)
    {
        DS_FreeString(s);
    }

    void FreeMetadata(Metadata* m)
    {
        DS_FreeMetadata(m);
    }

    char* Version()
    {
        return DS_Version();
    }

    char* ErrorCodeToErrorMessage(int aErrorCode)
    {
        return DS_ErrorCodeToErrorMessage(aErrorCode);
    }
}
