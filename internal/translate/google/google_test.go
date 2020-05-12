package google

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	translator := NewTranslator()
	ret, err := translator.Do("this is one of best things about teching the course. it's easy to see the principles and ideas in practice.", "en", "zh")
	if err != nil {
		t.Errorf("%+v", err)
	}
	fmt.Println(ret)
}
