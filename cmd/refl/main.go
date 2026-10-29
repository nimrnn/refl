package main

import (
	"fmt"
	"context"
	"os"
	"os/signal"
	"io"
	"syscall"
	"sync"
	"time"

	"refl/internal/config"
	"refl/internal/load"
	"refl/internal/probe"
)

const burst = 10

func run(args []string, stderr io.Writer) error {
	o, err := config.ParseFlags(args, stderr)
	if err != nil {
		return err
	}

	if err := o.Validate(); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	p := probe.New(o.Rate, burst)

	input := io.Reader(os.Stdin)
	if o.Input != "" {
		file, err := os.Open(o.Input)
		if err != nil {
			return fmt.Errorf("Open file %q: %w", o.Input, err)
		}
		defer file.Close()

		input = file
	}

	urls, err := load.LoadSource(input)
	if err != nil {
		return err
	}

	jc := make(chan string)
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < o.Workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for u := range jc {
				r, err := p.ProbeParams(ctx, u)
				if err != nil {
					fmt.Printf("Failed for %s. Continue anyway.", u)
					continue
				}

				fmt.Printf("%s\t%v\n", r.Url, r.Params)
			}
		}()
	}

	go func() {
		defer close(jc)
		for _, u := range urls {
			select {
				case jc <- u:
				case <-ctx.Done():
					return
			}
		}
	}()

	wg.Wait()
	fmt.Println("\nDuration: ", time.Since(start))

	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
