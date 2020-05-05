package voice

import "testing"

func TestExtractor_ExtractAudio(t *testing.T) {
	extractor := NewExtractor("16000")
	extractor.Valid()
	extractor.ExtractAudio("/Users/apple/go/src/github.com/JerryZhou343/ClosedCaption/bin/1-1-application.mp4",
		"/Users/apple/go/src/github.com/JerryZhou343/ClosedCaption/bin/1-1-application.mp3")
}
