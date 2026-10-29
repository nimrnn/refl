package probe

import (
	"fmt"
	"net/url"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/time/rate"
	"context"
	"net/http"
	"strings"
	"io"
	"time"
)

const canaryLength = 8

type Reflected struct {
	Url		string
	Params	[]string
}

type Probe struct {
	client		*http.Client
	limiter		*rate.Limiter
}

func New(rps float64, burst int) *Probe {
	return &Probe{
		client: &http.Client{Timeout: 10 * time.Second},
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

func canary() (string, error) {
	bytes := make([]byte, canaryLength)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("Generate canary: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func (p *Probe) send(ctx context.Context, targetUrl string) (string, error) {
	if err := p.limiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("Wait rate limit: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetUrl, nil)
	if err != nil {
		return "", fmt.Errorf("Build request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read response body: %w", err)
	}

	return string(body), nil
	
}

func (p *Probe) ProbeParams(ctx context.Context, rawUrl string) (Reflected, error) {
	var params []string
	
	u, err := url.Parse(rawUrl)
	if err != nil {
		return Reflected{}, fmt.Errorf("Parse url: %w", err)
	}

	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return Reflected{}, fmt.Errorf("Parse query parameters from %q: %w", rawUrl, err)
	}

	for key := range query {
		q := url.Values{}

		for k, v := range query {
			q[k] = append([]string(nil), v...)
		}
		
		c, err := canary()
		if err != nil {
			return Reflected{}, err
		}
		
		q.Set(key, c)
		u.RawQuery = q.Encode()

		response, err := p.send(ctx, u.String())
		if err != nil {
			return Reflected{}, err
		}

		if strings.Contains(response, c) {
			params = append(params, key)
		}
	}

	return Reflected{
		Url: rawUrl,
		Params: params,
	}, nil
	
}
