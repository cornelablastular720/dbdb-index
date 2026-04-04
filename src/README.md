# How the index is generated

`main.go` generates the `README.md` index by fetching and caching metadata for every database management system listed on [dbdb.io](https://dbdb.io) (Carnegie Mellon University's Database of Databases).

## Prerequisites

- Go 1.22+
- Internet access (to reach dbdb.io)

## Usage

```sh
# Full run — fetch sitemap, scrape all pages, generate README:
go run ./src -o README.md

# Use cached data, only fetch new/unfetched entries:
go run ./src -o README.md -cache databases.json

# Quick run — sitemap only, no page scraping (names derived from slugs):
go run ./src -o README.md --no-fetch

# Custom concurrency:
go run ./src -o README.md -workers 10
```

## How it works

1. **Discover** — Fetches `https://dbdb.io/sitemap.xml` and extracts all `/db/{slug}` URLs. The sitemap provides the canonical list of every database along with its last-modified date.

2. **Merge cache** — Loads `databases.json` (if present) and merges with the fresh sitemap. New databases are added; removed ones are pruned; existing entries keep their cached metadata.

3. **Scrape** — For each database not yet in cache, fetches its page at `https://dbdb.io/db/{slug}` and extracts:
   - **Name** — from `<h1>` tag (or `<title>` fallback)
   - **Description** — from `og:description` or `description` meta tag
   - **Data models** — from `/browse?data-model=` links on the page

   Scraping runs concurrently (default 20 workers) with 100ms per-request rate limiting to be polite.

4. **Cache** — Saves all fetched metadata to `databases.json` so subsequent runs are instant for already-scraped databases.

5. **Generate** — Renders the final `README.md` with:
   - Total database count and data model statistics
   - Compact table-of-contents with per-letter counts
   - Alphabetical sections (A–Z, `#` for numbers/symbols), each a Markdown table with database name (linked) and description

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-o` | `README.md` | Output file path |
| `-cache` | `databases.json` | Cache file for fetched metadata |
| `-workers` | `20` | Concurrent fetch workers |
| `--no-fetch` | `false` | Skip page scraping; generate from cache/sitemap only |

## Files

| File | Purpose |
|------|---------|
| `main.go` | Generator — sitemap parsing, page scraping, caching, README rendering |
| `databases.json` | Cached metadata (auto-generated, not committed) |
