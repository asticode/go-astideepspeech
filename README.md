[![GoReportCard](http://goreportcard.com/badge/github.com/asticode/go-astideepspeech)](http://goreportcard.com/report/github.com/asticode/go-astideepspeech)
[![GoDoc](https://godoc.org/github.com/asticode/go-astideepspeech?status.svg)](https://godoc.org/github.com/asticode/go-astideepspeech)

Golang bindings for Mozilla's/Coqui's [STT](https://github.com/coqui-ai/STT) speech-to-text library.

`astideepspeech` is compatible with version `v1.0.0` of `STT`.

# Installation  
## Install Coqui STT

- fetch an up-to-date `native_client.<your system>.tar.xz` matching your system from ["releases"](https://github.com/coqui-ai/STT/releases/tag/v1.0.0)
- extract its content to /tmp/stt/lib
- copy the downloaded `stt.so` and `coqui-stt.h` files
to directories that are searched by default, e.g. `/usr/local/lib` and
`/usr/local/include`, respectively.

## Install astideepspeech

Run the following command:

    $ go get -u github.com/asticode/go-astideepspeech/...
    
# Example       
## Get the pre-trained model and scorer

Sign up with your email and download the scorer and tflite files from eg https://coqui.ai/english/coqui/v1.0.0-large-vocab
    
## Get the audio files

Run the following commands: 

    $ cd /tmp/deepspeech
    $ wget https://github.com/coqui-ai/STT/releases/download/v1.0.0/audio-1.0.0.tar.gz
    $ tar xvfz audio-1.0.0.tar.gz
    
## Use the client

Run the following commands (make sure `$GOPATH/bin` is in your `$PATH`):

    $ cd /tmp/deepspeech
    $ stt -model stt-1.0.0-model.tflite -scorer stt-1.0.0-model.scorer -audio audio/2830-3980-0043.wav
    
        Text: experience proves this
    
    $ stt -model stt-1.0.0-model.tflite -scorer stt-1.0.0-model.scorer -audio audio/4507-16021-0012.wav
    
        Text: why should one hall on the way
        
    $ stt -model stt-1.0.0-model.tflite -scorer stt-1.0.0-model.scorer -audio audio/8455-210777-0068.wav
    
        Text: your power is sufficient i said
