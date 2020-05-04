package tencent

import "testing"

func TestDo(t *testing.T) {
	translator := NewTranslator()
	translator.Do("today I'd like to talk about GO", "en", "zh")
}
