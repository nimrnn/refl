# refl

A small, concurrent probe for finding query parameters whose values are reflected in the HTTP response body. It is intended as a quick first step when testing for reflected XSS.

## How it works

For each URL, `refl`:

1. Parses the query string.
2. Sends a baseline request for every parameter, replacing that parameter's value with a unique random canary while leaving all other parameters untouched.
3. Checks whether the canary appears verbatim in the response body.
4. Reports the URL and the parameters that reflected.

Requests are rate-limited and distributed across a pool of workers, so a large input list can be processed without hammering the target.

## Requirements

- Go 1.27 or newer (to build from source)
- A target you are authorized to test

## Build

```sh
git clone https://github.com/nimrnn/refl.git
cd refl
go build -o refl ./cmd/refl
```

## Usage

Provide URLs via stdin:

```sh
cat urls.txt | ./refl
```

Or from a file:

```sh
./refl -input urls.txt
```

Example:

```sh
$ printf '%s\n' 'https://example.com/search?q=hello&page=1' | ./refl
https://example.com/search?q=hello&page=1    [q]

Duration:  312ms
```

In this example, `q` was reflected in the response body while `page` was not.

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-input` | stdin | File containing one URL per line. |
| `-workers` | `5` | Number of concurrent workers. |
| `-rate` | `10` | Maximum requests per second across all workers. |
| `-output` | none | Append results to a file. URLs already present in the file are skipped. |

Blank lines in the input are ignored. Press `Ctrl+C` to stop early.

## Notes and limitations

- Only query parameters are tested; path segments, headers, and POST bodies are not.
- Detection is a simple canary substring match, so it can produce false positives if the response reflects the canary in an unexpected context (for example, an error page or log echo) or false negatives if the application transforms the value (for example, truncation or encoding).
- A reflected parameter is not automatically exploitable; confirm the context before reporting XSS.

## Disclaimer

Use this tool only against systems you own or have explicit permission to test. You are responsible for complying with applicable laws and program rules.
