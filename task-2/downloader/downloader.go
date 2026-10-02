package downloader

import (
	"context"
	"dakota5173/vk-spbu-backend/task-2/argparse"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
)

const URLTemplate = "https://homeworksite.site/%d/info.0.json"

type Movie struct {
	ID       uint64 `json:"id"`
	Title    string `json:"title"`
	Year     uint64 `json:"year"`
	Director string `json:"director"`
}

func DownloadOne(ctx context.Context, client *http.Client, id uint64) {
	link := fmt.Sprintf(URLTemplate, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%d] request error: %v\n", id, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		if !errors.Is(ctx.Err(), context.Canceled) {
			fmt.Fprintf(os.Stderr, "[%d] request failed: %v\n", id, err)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "[%d] HTTP status %d\n", id, resp.StatusCode)
		return
	}

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		if !errors.Is(ctx.Err(), context.Canceled) {
			fmt.Fprintf(os.Stderr, "[%d] json decode error: %v\n", id, err)
		}
		return
	}

	fmt.Printf("%d — %s — %d — %s\n", movie.ID, movie.Title, movie.Year, movie.Director)
}

func DownloadMovies(ctx context.Context, cfg *argparse.CommandFlags) {
	var wg sync.WaitGroup
	jobs := make(chan uint64)
	client := &http.Client{Timeout: cfg.Timeout}

	for range cfg.Workers {
		wg.Go(
			func() {
				for j := range jobs {
					DownloadOne(ctx, client, j)
				}
			})
	}

loop:
	for j := cfg.From; j <= cfg.To; j++ {
		select {
		case <-ctx.Done():
			break loop
		case jobs <- j:
		}
	}
	close(jobs)
	wg.Wait()
}
