package main

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

func TestMain_ConnectionRefused(t *testing.T) {
	cmd := exec.Command("go", "run", "./...")
	out, _ := cmd.CombinedOutput()

	actual := string(out)
	expected := "UNKNOWN - url and token are required"

	if !strings.Contains(actual, expected) {
		t.Error("\nActual: ", actual, "\nExpected: ", expected)
	}
}

type IntegrationTest struct {
	name     string
	server   *httptest.Server
	args     []string
	expected string
}

func TestMainCmd(t *testing.T) {
	tests := []IntegrationTest{
		{
			name: "invalid-json",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"foo}`))
			})),
			args:     []string{"run", "./...", "-T", "test", "--url"},
			expected: "UNKNOWN - could not decode JSON from body",
		},
		{
			name: "empty-response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"data" : [ {"id": "c86dc437", "name": "test1"}], "pagination": { "nextCursor": null, "totalItems": 2}}`))
			})),
			args:     []string{"run", "./...", "-T", "test", "--url"},
			expected: "WARNING -",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.server.Close()

			cmd := exec.Command("go", append(test.args, test.server.URL)...)
			out, _ := cmd.CombinedOutput()

			actual := string(out)

			if !strings.Contains(actual, test.expected) {
				t.Error("\nActual: ", actual, "\nExpected: ", test.expected)
			}
		})
	}
}
