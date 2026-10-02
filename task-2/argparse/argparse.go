package argparse

import (
	"errors"
	"flag"
	"fmt"
	"time"
)

type CommandFlags struct {
	From    uint64
	To      uint64
	Workers uint64
	Timeout time.Duration
}

func ParseFlags() (*CommandFlags, error) {
	from := flag.Uint64("from", 0, "ID of the first movie (required)")
	to := flag.Uint64("to", 0, "ID of the last movie (required)")
	workers := flag.Uint64("workers", 10, "number of workers")
	timeoutSec := flag.Int64("timeout", 5, "request timeout in seconds")

	flag.Parse()

	var hasFrom, hasTo bool
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "from" {
			hasFrom = true
		}
		if f.Name == "to" {
			hasTo = true
		}
	})

	if !hasFrom {
		return nil, errors.New("missing required flag: --from")
	}
	if !hasTo {
		return nil, errors.New("missing required flag: --to")
	}
	if *from == 0 {
		return nil, errors.New("--from must be greater than 0")
	}
	if *to == 0 {
		return nil, errors.New("--to must be greater than 0")
	}
	if *from > *to {
		return nil, fmt.Errorf("invalid range: --from > --to (%d > %d)", *from, *to)
	}
	if *workers == 0 {
		return nil, errors.New("--workers must be greater than 0")
	}
	if *timeoutSec <= 0 {
		return nil, errors.New("--timeout must be greater than 0")
	}

	return &CommandFlags{
		From:    *from,
		To:      *to,
		Workers: *workers,
		Timeout: time.Duration(*timeoutSec) * time.Second,
	}, nil
}
