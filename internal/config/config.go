package config

import (
	"flag"
	"fmt"
	"io"
)

type Options struct {
	Input		string
	Output		string
	Workers		int
	Rate		float64
}

func ParseFlags(args []string, stderr io.Writer) (Options, error) {
	var o Options

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.StringVar(&o.Input, "input", "", "File with list of urls to probe. (Default: stdin)")
	fs.StringVar(&o.Output, "output", "", "File name to save output. If not supplied, only print to stdout.")
	fs.IntVar(&o.Workers, "workers", 5, "Number of concurrent workers. (Default: 5)")
	fs.Float64Var(&o.Rate, "rate", 10.0, "Number of requests to send per second. (Default: 10/sec)")

	fs.Usage = func() {
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return Options{}, fmt.Errorf("Parsing arguments: %w", err)
	}

	return o, nil
}

func (o *Options) Validate() error {
	if o.Workers < 1 {
		return fmt.Errorf("Workers need to be at least one. Got %d", o.Workers)
	}

	return nil
}
