package api

import (
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

type APIClient struct {
	gClient graphql.Client
}

type authedTransport struct {
	wrapped      http.RoundTripper
	apiKey       string
	acceptHeader string
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+t.apiKey)
	req.Header.Set("Accept", t.acceptHeader)
	req.Header.Set("X-Github-Next-Global-ID", "1")
	return t.wrapped.RoundTrip(req)
}

func NewAPIClient(apiKey string) *APIClient {
	c := graphql.NewClient("https://api.github.com/graphql", &http.Client{
		Transport: &authedTransport{
			apiKey:       apiKey,
			wrapped:      http.DefaultTransport,
			acceptHeader: "application/json",
		},
	})

	return &APIClient{
		gClient: c,
	}
}
