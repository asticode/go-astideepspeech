#include <stdio.h>
#include <DeepSpeech/deepspeech.h>

extern "C" {
    class ModelWrapper {
        private:
            DeepSpeech::Model* model;
        public:
            ModelWrapper(const char* aModelPath, int aNCep, int aNContext, const char* aAlphabetConfigPath, int aBeamWidth)
            {
                model = new DeepSpeech::Model(aModelPath, aNCep, aNContext, aAlphabetConfigPath, aBeamWidth);
            }
            void enableDecoderWithLM(const char* aAlphabetConfigPath, const char* aLMPath, const char* aTriePath, float aLMWeight, float aWordCountWeight, float aValidWordCountWeight)
            {
                model->enableDecoderWithLM(aAlphabetConfigPath, aLMPath, aTriePath, aLMWeight, aWordCountWeight, aValidWordCountWeight);
            }
            char* stt(const short* aBuffer, unsigned int aBufferSize, int aSampleRate)
            {
                return model->stt(aBuffer, aBufferSize, aSampleRate);
            }
    };

    ModelWrapper* New(const char* aModelPath, int aNCep, int aNContext, const char* aAlphabetConfigPath, int aBeamWidth)
    {
        return new ModelWrapper(aModelPath, aNCep, aNContext, aAlphabetConfigPath, aBeamWidth);
    }
    void Close(ModelWrapper* w)
    {
        delete w;
    }
    void EnableDecoderWithLM(ModelWrapper* w, const char* aAlphabetConfigPath, const char* aLMPath, const char* aTriePath, float aLMWeight, float aWordCountWeight, float aValidWordCountWeight)
    {
        w->enableDecoderWithLM(aAlphabetConfigPath, aLMPath, aTriePath, aLMWeight, aWordCountWeight, aValidWordCountWeight);
    }
    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, int aSampleRate)
    {
        return w->stt(aBuffer, aBufferSize, aSampleRate);
    }
}