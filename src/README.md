# How the index is generated

This is a small Go program that builds the top-level `README.md` by pulling metadata for every database system listed on [dbdb.io](https://dbdb.io), the Database of Databases maintained by Carnegie Mellon University.

## Prerequisites

- Go 1.22 or newer
- Internet access to reach dbdb.io

## Usage

```sh
# Full run: fetch the sitemap, scrape every page, and generate the README
go run ./src -o README.md

# Incremental: reuse cached data, only fetch entries that are new
go run ./src -o README.md -cache databases.json

# Quick: skip scraping entirely, just use sitemap slugs as names
go run ./src -o README.md --no-fetch

# Dial down concurrency if you want to be extra polite
go run ./src -o README.md -workers 10
```

## How it works

1. **Discover.** The program fetches `https://dbdb.io/sitemap.xml` and pulls out every `/db/{slug}` URL. This is the canonical list of all databases on the site.

2. **Merge with cache.** If a `databases.json` file already exists from a previous run, it gets merged with the fresh sitemap. New databases are added, removed ones are pruned, and everything else keeps its cached metadata.

3. **Scrape.** For each database that hasn't been scraped yet, the program fetches its page and extracts:
   - **Name** from the `<h1>` tag (falls back to `<title>`)
   - **Description** from the `og:description` meta tag
   - **Data models, country, start year, project type, implementation languages, and licenses** from the browse filter links on the page

   This runs with 20 concurrent workers by default, with a 100ms delay between requests to keep things polite.

4. **Cache.** All the scraped metadata gets saved to `databases.json`. On the next run, only new or missing entries need fetching, so it finishes in seconds.

5. **Generate.** Finally, the program renders `README.md` with:
   - A prose summary and ASCII bar charts covering data models, decades, countries, languages, licenses, and project types
   - A compact A-Z table of contents with per-letter counts
   - Alphabetical sections, each a Markdown table with name, description, data model, country, year, type, language, and license

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-o` | `README.md` | Where to write the output |
| `-cache` | `databases.json` | Path to the metadata cache file |
| `-workers` | `20` | How many pages to fetch at once |
| `--no-fetch` | `false` | Skip scraping, generate from cache or sitemap only |

## Files

| File | Purpose |
|------|---------|
| `main.go` | The whole program: sitemap parsing, scraping, caching, and README rendering |
| `databases.json` | Cached metadata, auto-generated on each run (not committed) |
