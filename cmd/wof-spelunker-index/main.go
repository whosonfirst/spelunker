package main

import (
	"context"
	"log"

	"github.com/whosonfirst/spelunker/v2/app/index"
)

func main() {

	ctx := context.Background()
	err := index.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to index spelunker, %v", err)
	}
}
