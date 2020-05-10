package tencent

import (
	"fmt"
	"testing"
)

func TestGetTVTK(t *testing.T) {
	tv, tk, err := GetTVTK()
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Printf("tv:%s\n", tv)
	fmt.Printf("tk:%s\n", tk)
}
