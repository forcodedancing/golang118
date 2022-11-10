package golang118

import (
	ed25519std "crypto/ed25519"
	"testing"

	"golang.org/x/crypto/ed25519"
)

func Test_GoWork(t *testing.T) {
	public, private, _ := ed25519std.GenerateKey(nil)

	message := []byte("test message")
	sig := ed25519.Sign(private, message)
	if !ed25519.Verify(public, message, sig) {
		t.Errorf("valid signature rejected")
	}
}
