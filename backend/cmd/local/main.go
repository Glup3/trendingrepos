package main

import (
	"context"
	"fmt"
	"os"

	"github.com/glup3/trendingrepos/api"
	"github.com/glup3/trendingrepos/internal/loader"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("PAT_TOKEN")

	apiClient := api.NewAPIClient(apiKey)
	l := loader.NewLoader(apiClient)

	starCounts, err := l.CollectStarsUpperBounds(ctx, []string{"Go"}, []string{})
	if err != nil {
		fmt.Println("error happened", err)
	}

	for _, stars := range starCounts {
		fmt.Println(stars)
	}
}
