package load

import (
	"fmt"
	"bufio"
	"io"
	"strings"
	"os"
	"errors"
)

func LoadSource(r io.Reader) ([]string, error) {
	var read []string

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		read = append(read, line)
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("Reading input: %w", err)
	}

	return read, nil
}

func LoadExisting(path string) (map[string]bool, error) {
	seen := make(map[string]bool)

	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return seen, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Open file: %w", err)
	}

	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		// lines are "url\tparams"; the dedupe key is only the url
		url, _, _ := strings.Cut(line, "\t")
		seen[url] = true
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("Read file: %w", err)
	}

	return seen, nil
}
