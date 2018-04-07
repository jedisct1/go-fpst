package fpst

import (
	"testing"
)

func TestFPST(t *testing.T) {
	fpst := New()
	fpst = fpst.Insert([]byte("test"), 42)

	k, v := fpst.StartsWithExistingKey([]byte("test"))
	if k == nil || v != 42 {
		panic("not found/wrong value")
	}

	k, v = fpst.StartsWithExistingKey([]byte("testimonial"))
	if k == nil || v != 42 {
		panic("not found/wrong value")
	}

	k, v = fpst.StartsWithExistingKey([]byte("tes"))
	if k != nil || v != nil {
		panic("found/value")
	}
}
