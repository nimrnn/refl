package load

import (
	"fmt"
	"bufio"
	"io"
	"strings"
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
