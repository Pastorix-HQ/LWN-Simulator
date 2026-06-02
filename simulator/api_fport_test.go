package simulator

import (
	"testing"

	"github.com/arslab/lwnsimulator/socket"
)

// A typed uplink (e.g. Position Report on FPort 5) must be triggerable on a
// per-message FPort via the control surface, independent of the device's
// configured default FPort.
func TestResolveFPort(t *testing.T) {
	t.Run("returns the override when FPort > 0", func(t *testing.T) {
		got := resolveFPort(socket.NewPayload{FPort: 5})
		if got == nil {
			t.Fatal("expected an FPort override, got nil")
		}
		if *got != 5 {
			t.Fatalf("got FPort %d, want 5", *got)
		}
	})

	t.Run("returns nil when FPort is unset (0 keeps the device default)", func(t *testing.T) {
		if got := resolveFPort(socket.NewPayload{}); got != nil {
			t.Fatalf("expected nil, got %d", *got)
		}
	})

	t.Run("returns nil for out-of-range FPort", func(t *testing.T) {
		if got := resolveFPort(socket.NewPayload{FPort: -1}); got != nil {
			t.Fatalf("expected nil for negative FPort, got %d", *got)
		}
		if got := resolveFPort(socket.NewPayload{FPort: 256}); got != nil {
			t.Fatalf("expected nil for FPort > 255, got %d", *got)
		}
	})
}
