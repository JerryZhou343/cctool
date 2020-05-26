package utils

import (
	"testing"
)

func TestMillisDurationConv(t *testing.T) {
	cases := []struct {
		val    int64
		expect string
	}{
		{
			val:    340,
			expect: "00:00:00,340",
		},
		{
			val:    2365,
			expect: "00:00:02,365",
		},
	}
	for _, itr := range cases {
		if itr.expect != MillisDurationConv(itr.val) {
			t.Errorf("case val %d expect %s failed", itr.val, itr.expect)
		}
	}
}
