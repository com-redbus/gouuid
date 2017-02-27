package gouuid

import "testing"

func TestRng(t *testing.T) {
	_, err := rng(16)
	if err != nil {
		t.Errorf("could not generate random bytes")
	}

}
