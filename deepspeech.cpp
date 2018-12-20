#include <stdio.h>
#include <deepspeech.h>

extern "C" {
    class ModelWrapper {
        private:
            ModelState* model;

        public:
            ModelWrapper(const char* aModelPath, int aNCep, int aNContext, const char* aAlphabetConfigPath, int aBeamWidth)
            {
                DS_CreateModel(aModelPath, aNCep, aNContext, aAlphabetConfigPath, aBeamWidth, &model);
            }

            ~ModelWrapper()
            {
                DS_DestroyModel(model);
            }

            void enableDecoderWithLM(const char* aAlphabetConfigPath, const char* aLMPath, const char* aTriePath, float aLMWeight, float aValidWordCountWeight)
            {
                DS_EnableDecoderWithLM(model, aAlphabetConfigPath, aLMPath, aTriePath, aLMWeight, aValidWordCountWeight);
            }

            char* stt(const short* aBuffer, unsigned int aBufferSize, int aSampleRate)
            {
                return DS_SpeechToText(model, aBuffer, aBufferSize, aSampleRate);
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
    void EnableDecoderWithLM(ModelWrapper* w, const char* aAlphabetConfigPath, const char* aLMPath, const char* aTriePath, float aLMWeight, float aValidWordCountWeight)
    {
        w->enableDecoderWithLM(aAlphabetConfigPath, aLMPath, aTriePath, aLMWeight, aValidWordCountWeight);
    }
    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize, int aSampleRate)
    {
        return w->stt(aBuffer, aBufferSize, aSampleRate);
    }
}
