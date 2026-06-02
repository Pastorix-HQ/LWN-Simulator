package simulator

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/arslab/lwnsimulator/socket"
)

// The C37 heartbeat is a binary 11-byte 0x04 frame that cannot ride a UTF-8
// JSON string, so the control surface must accept it base64-encoded.
func TestResolvePayloadBytes(t *testing.T) {
	heartbeat := []byte{0x04, 0x64, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00}

	t.Run("raw payload is unchanged when Base64 is false", func(t *testing.T) {
		got, err := resolvePayloadBytes(socket.NewPayload{Payload: "hello", Base64: false})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "hello" {
			t.Fatalf("got %q, want %q", got, "hello")
		}
	})

	t.Run("base64 payload decodes to binary bytes", func(t *testing.T) {
		enc := base64.StdEncoding.EncodeToString(heartbeat)
		got, err := resolvePayloadBytes(socket.NewPayload{Payload: enc, Base64: true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(got, heartbeat) {
			t.Fatalf("got %v, want %v", got, heartbeat)
		}
		if len(got) != 11 || got[0] != 0x04 {
			t.Fatalf("decoded frame is not a valid 11-byte 0x04 heartbeat: %v", got)
		}
	})

	t.Run("invalid base64 returns an error", func(t *testing.T) {
		if _, err := resolvePayloadBytes(socket.NewPayload{Payload: "!!!not-base64", Base64: true}); err == nil {
			t.Fatal("expected an error for invalid base64, got nil")
		}
	})
}
