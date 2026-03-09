# secret-scraper

A Go-based HTTP service that scrapes public GitHub repositories for exposed secrets and API keys using configurable regex patterns.

## How It Works

1. You send a POST request specifying a `secret_provider` and `secret_type`.
2. The service looks up the matching search keyword and regex pattern from `config.yaml`.
3. It queries the GitHub Code Search API across N pages (iterations), fetches the raw content of each matched file, and extracts secrets using the configured regex.
4. Deduplicated secrets are returned in the response.

## Project Structure

```
secret-scraper/
├── cmd/main.go                          # Entry point, wires dependencies, starts Fiber HTTP server
├── config/config.go                     # Config loading from YAML
├── files/config/
│   ├── config.yaml                      # Active config (gitignored)
│   └── config.yaml.sample               # Sample config
├── handler/scraper/rest/
│   └── scrape_secret.go                 # POST /v1/scraper/scrape_secret handler
├── usecase/scraper/
│   └── scrape_secret.go                 # Core scraping logic: search → fetch → regex match
├── repository/
│   ├── github/client/                   # GitHub API: code search + raw file fetch
│   └── config/static/                   # Config-based scraper rule lookup
├── model/
│   ├── secret.go                        # Secret, SecretProvider, SecretType
│   └── github.go                        # GithubCode model
└── util/rest/response.go                # JSON error helper
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

```bash
go run ./cmd/main.go
```

The server starts on port `3000`.

## API

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

## Dependencies

| Package | Purpose |
|---------|---------|
| `gofiber/fiber/v3` | HTTP framework |
| `google/go-github/v84` | GitHub API client |
| `gopkg.in/yaml.v3` | Config parsing |

## License

MIT