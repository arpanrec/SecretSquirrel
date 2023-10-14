package storage

import (
	"testing"
)

func TestAbs(t *testing.T) {
	// encrypt plain text message using public key
	ss := "Hello World"
	encryptMessage(&ss)
	t.Log(ss)
	decryptMessage(&ss)
	t.Log(ss)
}
