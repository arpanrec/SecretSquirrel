package encryption

import (
	"testing"
)

func TestAbs(t *testing.T) {
	// encrypt plain text message using public key
	ss := "Hello World"
	EncryptMessage(&ss)
	t.Log(ss)
	DecryptMessage(&ss)
	t.Log(ss)
}
