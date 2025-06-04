package loader

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/glup3/trendingrepos/internal/api"
)

type searchResponse struct {
	Repos []api.GitHubRepo
	Err   error
}

type apiClientMock struct {
	mu        sync.Mutex
	responses []searchResponse
	count     int
}

func (m *apiClientMock) SearchRepos(ctx context.Context, queryArgs api.QueryArgs) ([]api.GitHubRepo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.count >= len(m.responses) {
		return nil, fmt.Errorf(
			"apiClientMock: SearchRepos called more times than expected (call %d, expected %d)",
			m.count+1,
			len(m.responses),
		)
	}

	resp := m.responses[m.count]
	m.count++
	return resp.Repos, resp.Err
}

func TestLoadRepos(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses []searchResponse
		expectedLen   int
		expectedError bool
	}{
		{
			name: "all unique",
			mockResponses: []searchResponse{
				{Repos: buildRepos(t, 100, "a"), Err: nil},
				{Repos: buildRepos(t, 100, "s"), Err: nil},
				{Repos: buildRepos(t, 100, "d"), Err: nil},
				{Repos: buildRepos(t, 100, "f"), Err: nil},
				{Repos: buildRepos(t, 100, "g"), Err: nil},
				{Repos: buildRepos(t, 100, "h"), Err: nil},
				{Repos: buildRepos(t, 100, "j"), Err: nil},
				{Repos: buildRepos(t, 100, "k"), Err: nil},
				{Repos: buildRepos(t, 100, "l"), Err: nil},
				{Repos: buildRepos(t, 100, ";"), Err: nil},
			},
			expectedLen:   1000,
			expectedError: false,
		},
		{
			name: "with duplicates",
			mockResponses: []searchResponse{
				{Repos: buildRepos(t, 100, "ditto"), Err: nil},
				{Repos: buildRepos(t, 100, "ditto"), Err: nil},
				{Repos: nil, Err: nil},
				{Repos: buildRepos(t, 100, "dittto"), Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
			},
			expectedLen:   200,
			expectedError: false,
		},
		{
			name: "with single failure",
			mockResponses: []searchResponse{
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
			},
			expectedLen:   0,
			expectedError: true,
		},
		{
			name: "with multiple failure",
			mockResponses: []searchResponse{
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: nil},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: nil},
			},
			expectedLen:   0,
			expectedError: true,
		},
		{
			name: "with all failure",
			mockResponses: []searchResponse{
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
				{Repos: nil, Err: fmt.Errorf("504")},
			},
			expectedLen:   0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader(&apiClientMock{
				responses: tt.mockResponses,
			}, slog.Default())

			repos, err := l.LoadRepos(context.Background(), 300)
			if !tt.expectedError && err != nil {
				t.Errorf("expected no error, got: %v", err)
				return
			}
			if tt.expectedError && err == nil {
				t.Errorf("expected error, got: nil")
				return
			}
			if len(repos) != tt.expectedLen {
				t.Errorf("expected a total of %d repos, got: %d", tt.expectedLen, len(repos))
				return
			}
		})
	}
}

func buildRepos(t testing.TB, cap int, prefix string) []api.GitHubRepo {
	t.Helper()
	repos := make([]api.GitHubRepo, 0, cap)
	for i := range cap {
		repos = append(repos, api.GitHubRepo{Id: fmt.Sprintf("%s%d", prefix, i)})
	}
	return repos
}
