package cluster

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/validation"
)

func TestAddressToHostnameOverride(t *testing.T) {
	tests := []struct {
		address  string
		expected string
	}{
		{"2001:db8::1", "2001-db8-1"},
		{"fe80::1", "fe80-1"},
		{"::1", "1"},
		{"192.168.1.1", "192.168.1.1"},
		{"node1.example.com", "node1.example.com"},
	}

	for _, tc := range tests {
		t.Run(tc.address, func(t *testing.T) {
			got := addressToHostnameOverride(tc.address)
			if got != tc.expected {
				t.Fatalf("addressToHostnameOverride(%q) = %q, want %q", tc.address, got, tc.expected)
			}
			if errs := validation.IsDNS1123Subdomain(got); len(errs) > 0 {
				t.Fatalf("hostname %q is not a valid DNS1123 subdomain: %v", got, errs)
			}
		})
	}
}