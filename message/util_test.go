package message

import "testing"

func TestOlderVersion(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want bool
	}{
		{"1", "2", true},
		{"1.1", "1.02", true},
		{"1.10", "1.02", false},
		{"1.1", "1.1.1", true},
		{"1.1", "1.1a", false},
		{"1.1", "undefined", false},
	}
	for _, tt := range tests {
		t.Run(tt.a+"<"+tt.b, func(t *testing.T) {
			if got := OlderVersion(tt.a, tt.b); got != tt.want {
				t.Errorf("OlderVersion(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
