package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.roundTripFunc != nil {
		return m.roundTripFunc(req)
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}, nil
}

func TestAuthedTransport_RoundTrip(t *testing.T) {
	tests := []struct {
		name                        string
		apiKey                      string
		acceptHeader                string
		expectedAuth                string
		expectedAccept              string
		expectedXGithubNextGlobalID string
	}{
		{
			name:                        "All headers set correctly",
			apiKey:                      "test-api-key",
			acceptHeader:                "application/vnd.github.v4.graphql",
			expectedAuth:                "bearer test-api-key",
			expectedAccept:              "application/vnd.github.v4.graphql",
			expectedXGithubNextGlobalID: "1",
		},
		{
			name:                        "Empty API key",
			apiKey:                      "",
			acceptHeader:                "application/json",
			expectedAuth:                "bearer ",
			expectedAccept:              "application/json",
			expectedXGithubNextGlobalID: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					if auth := req.Header.Get("Authorization"); auth != tt.expectedAuth {
						t.Errorf("Expected Authorization header '%s', got '%s'", tt.expectedAuth, auth)
					}
					if accept := req.Header.Get("Accept"); accept != tt.expectedAccept {
						t.Errorf("Expected Accept header '%s', got '%s'", tt.expectedAccept, accept)
					}
					if globalID := req.Header.Get("X-Github-Next-Global-ID"); globalID != tt.expectedXGithubNextGlobalID {
						t.Errorf("Expected X-Github-Next-Global-ID header '%s', got '%s'", tt.expectedXGithubNextGlobalID, globalID)
					}
					return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}, nil
				},
			}

			at := &authedTransport{
				wrapped:      mockTransport,
				apiKey:       tt.apiKey,
				acceptHeader: tt.acceptHeader,
			}

			req, _ := http.NewRequest("GET", "http://example.com", nil)
			_, err := at.RoundTrip(req)
			if err != nil {
				t.Errorf("RoundTrip returned an error: %v", err)
			}
		})
	}
}
