[![GoReportCard](http://goreportcard.com/badge/github.com/asticode/go-astideepspeech)](http://goreportcard.com/report/github.com/asticode/go-astideepspeech)
[![GoDoc](https://godoc.org/github.com/asticode/go-astideepspeech?status.svg)](https://godoc.org/github.com/asticode/go-astideepspeech)

Golang bindings for Mozilla's [DeepSpeech](https://github.com/mozilla/DeepSpeech) speech-to-text library.

As of now, `astideepspeech` is only compatible with version `v0.4.0` of `DeepSpeech`.

# Installation
## Install DeepSpeech

- fetch an uptodate `native_client.tar.xz` matching your system, from DeepSpeech's ["releases"](https://github.com/mozilla/DeepSpeech/releases/tag/v0.4.0)
- download `deepspeech.h` from https://github.com/mozilla/DeepSpeech/raw/v0.4.0/native_client/deepspeech.h
- extract/copy both to /tmp/ds/
- export CGO_LDFLAGS="-L/tmp/ds/"
- export CGO_CXXFLAGS="-I/tmp/ds/"

## Install astideepspeech

Run the following command:

    $ go get -u github.com/asticode/go-astideepspeech/...
    
# Example
## Get the pre-trained model

Run the following commands:

    $ mkdir /tmp/deepspeech
    $ cd /tmp/deepspeech
    $ wget https://github.com/mozilla/DeepSpeech/releases/download/v0.4.0/deepspeech-0.4.0-models.tar.gz
    $ tar xvfz deepspeech-0.4.0-models.tar.gz
    
## Get the audio files

Run the following commands:

    $ cd /tmp/deepspeech
    $ wget https://github.com/mozilla/DeepSpeech/releases/download/v0.4.0/audio-0.4.0.tar.gz
    $ tar xvfz audio-0.4.0.tar.gz
    
## Use the client

Run the following commands (make sure `$GOPATH/bin` is in your `$PATH`):

    $ cd /tmp/deepspeech
    $ deepspeech models/output_graph.pb audio/2830-3980-0043.wav models/alphabet.txt models/lm.binary models/trie
    
        Text: experience proves this
    
    $ deepspeech models/output_graph.pb audio/4507-16021-0012.wav models/alphabet.txt models/lm.binary models/trie
    
        Text: why should one halt on the way
        
    $ deepspeech models/output_graph.pb audio/8455-210777-0068.wav models/alphabet.txt models/lm.binary models/trie
    
        Text: your power is sufficient i said
