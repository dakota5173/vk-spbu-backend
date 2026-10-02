package main

import (
	"context"
	"dakota5173/vk-spbu-backend/task-2/argparse"
	"dakota5173/vk-spbu-backend/task-2/downloader"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := argparse.ParseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	downloader.DownloadMovies(ctx, cfg)
}
