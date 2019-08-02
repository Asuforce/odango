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

func TestValidateEndpoint(t *testing.T) {
	c := Config{}
	tests := []struct {
		i, o string
	}{
		{"", "/deploy/"},
		{"deploy", "/deploy/"},
		{"/deploy", "/deploy/"},
		{"deploy/", "/deploy/"},
	}

	for _, tt := range tests {
		c.Server.Endpoint = tt.i
		c.validateEndpoint()
		if c.Server.Endpoint != tt.o {
			t.Fatalf("want = %s, got = %s", tt.o, tt.i)
		}
	}
}

func TestValidatePort(t *testing.T) {
	c := Config{}
	tests := []struct {
		i, o int
	}{
		{0, 8080},
		{8080, 8080},
		{3000, 3000},
	}

	for _, tt := range tests {
		c.Server.Port = tt.i
		c.validatePort()
		if c.Server.Port != tt.o {
			t.Fatalf("want = %d, got = %d", tt.o, tt.i)
		}
	}
}
