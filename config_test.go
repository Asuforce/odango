package main

import "testing"

func TestFormatPath(t *testing.T) {

	tests := []struct {
		i, o string
	}{
		{"/foo", "/foo/"},
		{"bar/", "/bar/"},
		{"/buz/", "/buz/"},
		{"hoge", "/hoge/"},
	}

	for _, tt := range tests {
		if formatPath(tt.i) != tt.o {
			t.Fatalf("want = %s, got = %s", tt.o, tt.i)
		}
	}
}
