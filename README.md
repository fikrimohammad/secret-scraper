# secret-scraper

A Go-based tool that scrapes public GitHub repositories for exposed secrets and API keys using configurable regex patterns. Supports both HTTP server and CLI modes.

## How It Works

1. You specify a `secret_provider` and `secret_type` (via API or CLI).
2. The tool looks up the matching search keyword and regex pattern from `config.yaml`.
3. It queries the GitHub Code Search API across N pages (iterations), fetches the raw content of each matched file in parallel, and extracts secrets using the configured regex.
4. Deduplicated secrets are returned in the response.

## Project Structure

```
secret-scraper/
├── cmd/
│   ├── main.go                            # Entry point, routes to serve/scrape subcommand
│   ├── serve.go                           # HTTP server mode
│   └── scrape.go                          # CLI one-shot mode
├── config/config.go                       # Config loading from YAML
├── files/config/
│   ├── config.yaml                        # Active config (gitignored)
│   └── config.yaml.sample                 # Sample config
├── handler/scraper/rest/
│   └── scrape_secret.go                   # POST /v1/scraper/scrape_secret handler
├── usecase/scraper/
│   └── scrape_secret.go                   # Core scraping logic: search → parallel fetch → regex match
├── repository/
│   ├── github/client/                     # GitHub API: code search + raw file fetch (rate-limited)
│   └── config/static/                     # Config-based scraper rule lookup
├── model/
│   ├── secret.go                          # Secret, SecretProvider, SecretType
│   └── github.go                          # GithubCode model
└── util/rest/response.go                  # JSON error helper
```

## Prerequisites

- Go 1.22+
- A GitHub personal access token with `repo` (or `public_repo`) scope

## Setup

```bash
git clone https://github.com/fikrimohammad/secret-scraper.git
cd secret-scraper

cp files/config/config.yaml.sample files/config/config.yaml
# Edit config.yaml and set your GitHub access token
```

### `files/config/config.yaml`

```yaml
github:
  access_token: YOUR_GITHUB_TOKEN

secret_scraper:
  - secret_provider: anthropic
    secret_type: anthropic_api_key
    secret_query_keyword: sk-ant-api03
    secret_regex_pattern: sk-ant-api03-[a-zA-Z0-9\-_]+
  - secret_provider: anthropic
    secret_type: anthropic_admin_key
    secret_query_keyword: sk-ant-admin01
    secret_regex_pattern: sk-ant-admin01-[a-zA-Z0-9\-_]+
```

Add more entries to scan for other secret types (AWS, OpenAI, etc.).

## Running

### HTTP Server Mode

```bash
# Using go run
go run ./cmd/main.go serve

# Using make
make run

# Build and run
make dev
```

The server starts on port `3000`.

### CLI Mode

Run a one-shot scrape directly from the command line:

```bash
# Using go run
go run ./cmd/main.go scrape -provider anthropic -type anthropic_api_key

# With custom pagination
go run ./cmd/main.go scrape -provider anthropic -type anthropic_api_key -iterations 5 -limit 20

# Using the built binary
./bin/secret-scraper scrape -provider openai -type openai_api_key
```

#### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-provider` | required | Secret provider (e.g. `anthropic`, `openai`) |
| `-type` | required | Secret type (e.g. `anthropic_api_key`) |
| `-iterations` | `10` | Max number of search pages to scan |
| `-limit` | `10` | Max results per search page |

Output is JSON printed to stdout:

```json
{
  "Data": [
    {
      "provider": "anthropic",
      "type": "anthropic_api_key",
      "value": "sk-ant-api03-..."
    }
  ]
}
```

## API (Server Mode)

### `POST /v1/scraper/scrape_secret`

**Request body:**

```json
{
  "secret_provider": "anthropic",
  "secret_type": "anthropic_api_key",
  "max_limit_per_iterations": 10,
  "max_iterations": 5
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `secret_provider` | string | required | Provider name (e.g. `anthropic`) |
| `secret_type` | string | required | Secret type (e.g. `anthropic_api_key`) |
| `max_limit_per_iterations` | int | 10 | Results per GitHub search page |
| `max_iterations` | int | 10 | Number of pages to scan |

**Response:**

```json
{
  "data": [
    {
      "provider": "anthropic",
      "type": "anthropic_api_key",
      "value": "sk-ant-api03-..."
    }
  ]
}
```

## Rate Limiting

GitHub API requests are automatically rate-limited (1 request/sec, burst of 5) to stay within GitHub's authenticated rate limit of 5,000 requests/hour. File fetches that fail (e.g., file too large, network error) are skipped and logged rather than failing the entire scrape.

## File Size Limits

Files larger than 5MB are automatically skipped to prevent excessive memory usage. This is enforced both via the `Content-Length` header and by capping reads with `io.LimitReader`.

## Dependencies

| Package | Purpose |
|---------|---------|
| `gofiber/fiber/v3` | HTTP framework |
| `google/go-github/v84` | GitHub API client |
| `gopkg.in/yaml.v3` | Config parsing |
| `golang.org/x/time/rate` | API rate limiting |

## License

MIT
