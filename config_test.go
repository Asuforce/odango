package main

import "testing"

func TestFormatPath(t *testing.T) {
	c := Config{}

	tests := []struct {
		i, o string
	}{
		{"/foo", "/foo/"},
		{"bar/", "/bar/"},
		{"/buz/", "/buz/"},
	}

	for _, tt := range tests {
		if c.formatPath(tt.i) != tt.o {
			t.Fatalf("want = %s, got = %s", tt.o, tt.i)
		}
	}
}
