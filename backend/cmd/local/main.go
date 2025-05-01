package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/glup3/trendingrepos/api"
	"github.com/glup3/trendingrepos/internal/loader"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("PAT_TOKEN")

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	apiClient := api.NewAPIClient(apiKey)
	l := loader.NewLoader(apiClient, logger)

	starCounts, err := l.CollectStarsUpperBounds(ctx, []string{"Go"}, []string{})
	if err != nil {
		fmt.Println("error happened", err)
	}

	for _, stars := range starCounts {
		fmt.Println(stars)
	}
}
