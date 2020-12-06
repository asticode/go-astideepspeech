[![GoReportCard](http://goreportcard.com/badge/github.com/asticode/go-astideepspeech)](http://goreportcard.com/report/github.com/asticode/go-astideepspeech)
[![GoDoc](https://godoc.org/github.com/asticode/go-astideepspeech?status.svg)](https://godoc.org/github.com/asticode/go-astideepspeech)

Golang bindings for Mozilla's [DeepSpeech](https://github.com/mozilla/DeepSpeech) speech-to-text library.

`astideepspeech` is compatible with version `v0.9.0` of `DeepSpeech`.

# Installation
## Install DeepSpeech

- fetch an up-to-date `native_client.<your system>.tar.xz` matching your system from DeepSpeech's ["releases"](https://github.com/mozilla/DeepSpeech/releases/tag/v0.9.0)
- extract its content to /tmp/deepspeech/lib
- download `deepspeech.h` from https://github.com/mozilla/DeepSpeech/raw/v0.9.0/native_client/deepspeech.h
- copy it to /tmp/deepspeech/include
- export CGO_LDFLAGS="-L/tmp/deepspeech/lib/"
- export CGO_CXXFLAGS="-I/tmp/deepspeech/include/"
- export LD_LIBRARY_PATH=/tmp/deepspeech/lib/:$LD_LIBRARY_PATH

Alternatively, copy the downloaded `libdeepspeech.so` and `deepspeech.h` files
to directories that are searched by default, e.g. `/usr/local/lib` and
`/usr/local/include`, respectively.

## Install astideepspeech

Run the following command:

    $ go get -u github.com/asticode/go-astideepspeech/...
    
# Example
## Get the pre-trained model and scorer

Run the following commands:

    $ mkdir /tmp/deepspeech
    $ cd /tmp/deepspeech
    $ wget https://github.com/mozilla/DeepSpeech/releases/download/v0.9.0/deepspeech-0.9.0-models.pbmm
    $ wget https://github.com/mozilla/DeepSpeech/releases/download/v0.9.0/deepspeech-0.9.0-models.scorer
    
## Get the audio files

Run the following commands:

    $ cd /tmp/deepspeech
    $ wget https://github.com/mozilla/DeepSpeech/releases/download/v0.9.0/audio-0.9.0.tar.gz
    $ tar xvfz audio-0.9.0.tar.gz
    
## Use the client

Run the following commands (make sure `$GOPATH/bin` is in your `$PATH`):

    $ cd /tmp/deepspeech
    $ deepspeech -model deepspeech-0.9.0-models.pbmm -scorer deepspeech-0.9.0-models.scorer -audio audio/2830-3980-0043.wav
    
        Text: experience proves this
    
    $ deepspeech -model deepspeech-0.9.0-models.pbmm -scorer deepspeech-0.9.0-models.scorer -audio audio/4507-16021-0012.wav
    
        Text: why should one hall on the way
        
    $ deepspeech -model deepspeech-0.9.0-models.pbmm -scorer deepspeech-0.9.0-models.scorer -audio audio/8455-210777-0068.wav
    
        Text: your power is sufficient i said
