package random

import (
	"testing"
)

func TestString(t *testing.T) {
	t.Log(String(16))
}

func TestInt(t *testing.T) {
	t.Log(Int(1, 10))
}

func TestFloat(t *testing.T) {
	f, err := Float(16)
	if err != nil {
		t.Error(err)
	}
	t.Log(f)
	t.Logf("%.16f", f)
}
