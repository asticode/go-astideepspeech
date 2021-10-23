#include <stdio.h>
#include <coqui-stt.h>

extern "C" {
    class ModelWrapper {
        private:
            ModelState* model;

        public:
            ModelWrapper(const char* aModelPath, int *errorOut)
            {
                model = nullptr;
                *errorOut = STT_CreateModel(aModelPath, &model);
            }

            ~ModelWrapper()
            {
                if (model) {
                    STT_FreeModel(model);
                    model = nullptr;
                }
            }

            unsigned int beamWidth()
            {
                return STT_GetModelBeamWidth(model);
            }

            int setBeamWidth(unsigned int aBeamWidth)
            {
                return STT_SetModelBeamWidth(model, aBeamWidth);
            }

            int sampleRate()
            {
                return STT_GetModelSampleRate(model);
            }

            int enableExternalScorer(const char* aScorerPath)
            {
                return STT_EnableExternalScorer(model, aScorerPath);
            }

            int disableExternalScorer()
            {
                return STT_DisableExternalScorer(model);
            }

            int setScorerAlphaBeta(float aAlpha, float aBeta)
            {
                return STT_SetScorerAlphaBeta(model, aAlpha, aBeta);
            }

            char* stt(const short* aBuffer, unsigned int aBufferSize)
            {
                return STT_SpeechToText(model, aBuffer, aBufferSize);
            }

            Metadata* sttWithMetadata(const short* aBuffer, unsigned int aBufferSize, unsigned int aNumResults)
            {
                return STT_SpeechToTextWithMetadata(model, aBuffer, aBufferSize, aNumResults);
            }

            ModelState* getModel()
            {
                return model;
            }
    };

    ModelWrapper* New(const char* aModelPath, int* errorOut)
    {
        auto mw = new ModelWrapper(aModelPath, errorOut);
        if (*errorOut != STT_ERR_OK) {
            delete mw;
            mw = nullptr;
        }
        return mw;
    }
    void Model_Close(ModelWrapper* w)
    {
        delete w;
    }

    unsigned int Model_BeamWidth(ModelWrapper* w)
    {
        return w->beamWidth();
    }

    int Model_SetBeamWidth(ModelWrapper* w, unsigned int aBeamWidth)
    {
        return w->setBeamWidth(aBeamWidth);
    }

    int Model_SampleRate(ModelWrapper* w)
    {
        return w->sampleRate();
    }

    int Model_EnableExternalScorer(ModelWrapper* w, const char* aScorerPath)
    {
        return w->enableExternalScorer(aScorerPath);
    }

    int Model_DisableExternalScorer(ModelWrapper* w)
    {
        return w->disableExternalScorer();
    }

    int Model_SetScorerAlphaBeta(ModelWrapper* w, float aAlpha, float aBeta)
    {
        return w->setScorerAlphaBeta(aAlpha, aBeta);
    }

    char* Model_STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize)
    {
        return w->stt(aBuffer, aBufferSize);
    }

    Metadata* Model_STTWithMetadata(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, unsigned int aNumResults)
    {
        return w->sttWithMetadata(aBuffer, aBufferSize, aNumResults);
    }

    const CandidateTranscript* Metadata_Transcripts(Metadata* m)
    {
        return m->transcripts;
    }

    unsigned int Metadata_NumTranscripts(Metadata* m)
    {
        return m->num_transcripts;
    }

    void Metadata_Close(Metadata* m)
    {
        STT_FreeMetadata(m);
    }

    const TokenMetadata* CandidateTranscript_Tokens(CandidateTranscript* ct)
    {
        return ct->tokens;
    }

    int CandidateTranscript_NumTokens(CandidateTranscript* ct)
    {
        return ct->num_tokens;
    }

    double CandidateTranscript_Confidence(CandidateTranscript* ct)
    {
        return ct->confidence;
    }

    const char* TokenMetadata_Text(TokenMetadata* tm)
    {
        return tm->text;
    }

    unsigned int TokenMetadata_Timestep(TokenMetadata* tm)
    {
        return tm->timestep;
    }

    float TokenMetadata_StartTime(TokenMetadata* tm)
    {
        return tm->start_time;
    }

    class StreamWrapper {
        private:
            StreamingState* s;

        public:
            StreamWrapper(ModelWrapper* w, int* errorOut)
            {
                s = nullptr;
                *errorOut = STT_CreateStream(w->getModel(), &s);
            }

            ~StreamWrapper()
            {
                if (s) {
                    STT_FreeStream(s);
                    s = nullptr;
                }
            }

            void feedAudioContent(const short* aBuffer, unsigned int aBufferSize)
            {
                STT_FeedAudioContent(s, aBuffer, aBufferSize);
            }

            char* intermediateDecode()
            {
                return STT_IntermediateDecode(s);
            }

            Metadata* intermediateDecodeWithMetadata(unsigned int aNumResults)
            {
                return STT_IntermediateDecodeWithMetadata(s, aNumResults);
            }

            char* finish()
            {
                // STT_FinishStream frees the supplied state pointer.
                char* res = STT_FinishStream(s);
                s = nullptr;
                return res;
            }

            Metadata* finishWithMetadata(unsigned int aNumResults)
            {
                // STT_FinishStreamWithMetadata frees the supplied state pointer.
                Metadata* m = STT_FinishStreamWithMetadata(s, aNumResults);
                s = nullptr;
                return m;
            }

            void discard()
            {
                STT_FreeStream(s);
                s = nullptr;
            }
    };

    StreamWrapper* Model_NewStream(ModelWrapper* mw, int* errorOut)
    {
        auto sw = new StreamWrapper(mw, errorOut);
        if (*errorOut != STT_ERR_OK) {
            delete sw;
            sw = nullptr;
        }
        return sw;
    }
    void Stream_Discard(StreamWrapper* sw)
    {
        sw->discard();
        delete sw;
    }

    void Stream_FeedAudioContent(StreamWrapper* sw, const short* aBuffer, unsigned int aBufferSize)
    {
        sw->feedAudioContent(aBuffer, aBufferSize);
    }

    char* Stream_IntermediateDecode(StreamWrapper* sw)
    {
        return sw->intermediateDecode();
    }

    Metadata* Stream_IntermediateDecodeWithMetadata(StreamWrapper* sw, unsigned int aNumResults)
    {
        return sw->intermediateDecodeWithMetadata(aNumResults);
    }

    char* Stream_Finish(StreamWrapper* sw)
    {
        char* str = sw->finish();
        delete sw;
        return str;
    }

    Metadata* Stream_FinishWithMetadata(StreamWrapper* sw, unsigned int aNumResults)
    {
        Metadata* m = sw->finishWithMetadata(aNumResults);
        delete sw;
        return m;
    }

    void FreeString(char* s)
    {
        STT_FreeString(s);
    }

    char* Version()
    {
        return STT_Version();
    }

    char* ErrorCodeToErrorMessage(int aErrorCode)
    {
        return STT_ErrorCodeToErrorMessage(aErrorCode);
    }
}
