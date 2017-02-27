package gouuid

import "testing"
import "strings"

const TIME = 1257894000000000000

func TestV1(t *testing.T) {
	u1 := NewV1()
	u2 := NewV1()
	if strings.Compare(u1.Format(), u2.Format()) == 1 {
		t.Errorf("uuid order did not match")
	}
}

func TestV4(t *testing.T) {
	u1 := NewV4()
	if u1.Format() == "" {
		t.Errorf("uuid version 4 not created")

	}
}
