package main

import (
	"context"
	"log"

	"github.com/whosonfirst/spelunker/v2/app/httpd/server"
)

func main() {

	ctx := context.Background()
	err := server.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run server, %v", err)
	}
}
