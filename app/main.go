package main

import (
	"context"
	"image-app/internal/server"
)

func main() {
	ctx := context.Background()
	server.StartServer(ctx)
}
