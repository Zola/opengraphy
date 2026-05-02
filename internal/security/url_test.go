package security

import (
	"net"
	"testing"
)

func TestNormalizeAddsHTTPS(t *testing.T) {
	normalized, u, err := Normalize("example.com/path")
	if err != nil {
		t.Fatal(err)
	}
	if normalized != "https://example.com/path" {
		t.Fatalf("normalized = %q", normalized)
	}
	if u.Scheme != "https" {
		t.Fatalf("scheme = %q", u.Scheme)
	}
}

func TestBlockedIPs(t *testing.T) {
	cases := []string{"127.0.0.1", "10.0.0.2", "172.16.0.1", "192.168.1.1", "::1", "fc00::1", "fe80::1"}
	for _, tc := range cases {
		if !IsBlockedIP(net.ParseIP(tc)) {
			t.Fatalf("expected %s to be blocked", tc)
		}
	}
}

func TestPublicIPAllowed(t *testing.T) {
	if IsBlockedIP(net.ParseIP("8.8.8.8")) {
		t.Fatal("8.8.8.8 should be allowed")
	}
}
