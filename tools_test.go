package toolkit_test

import (
	"github.com/GiTweeker/toolkit"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools toolkit.Tools
	s := testTools.RandomString(10)

	if len(s) != 10 {
		t.Error("wrong length of random string returned")
	}
}
