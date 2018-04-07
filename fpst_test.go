package fpst

import (
	"testing"
)

func TestFPST(t *testing.T) {
	fpst := New()
	fpst = fpst.Insert([]byte("test"), 42)
	fpst = fpst.Insert([]byte("brouette"), 42)
	fpst = fpst.Insert([]byte("champignon"), 42)
	fpst = fpst.Insert([]byte("champagne"), 42)

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

	k, v = fpst.StartsWithExistingKey([]byte("brouette"))
	if k == nil || v != 42 {
		panic("not found/wrong value")
	}

	k, v = fpst.StartsWithExistingKey([]byte("champignon"))
	if k == nil || v != 42 {
		panic("not found/wrong value")
	}

	k, v = fpst.StartsWithExistingKey([]byte("champagne"))
	if k == nil || v != 42 {
		panic("not found/wrong values")
	}
}
