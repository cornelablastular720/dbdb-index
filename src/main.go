// Package main generates the dbdb-index README.md by fetching and indexing
// all database management systems catalogued at https://dbdb.io.
//
// It parses the sitemap for discovery, scrapes individual pages for metadata
// (name, description, data models), caches results locally, and renders an
// alphabetically-organised Markdown index.
package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"html"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

const baseURL = "https://dbdb.io"

// ── Sitemap types ──────────────────────────────────────────────────

type urlSet struct {
	URLs []siteURL `xml:"url"`
}

type siteURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

// ── Database entry ─────────────────────────────────────────────────

// Database holds the metadata for a single DBMS.
type Database struct {
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Desc      string   `json:"description,omitempty"`
	Models    []string `json:"data_models,omitempty"`
	Country   string   `json:"country,omitempty"`
	StartYear string   `json:"start_year,omitempty"`
	ProjTypes []string `json:"project_types,omitempty"`
	WrittenIn []string `json:"written_in,omitempty"`
	Licenses  []string `json:"licenses,omitempty"`
	LastMod   string   `json:"lastmod,omitempty"`
	Fetched   bool     `json:"fetched,omitempty"`
}

func (db *Database) pageURL() string { return baseURL + "/db/" + db.Slug }

// ── HTTP helper ────────────────────────────────────────────────────

var httpClient = &http.Client{Timeout: 30 * time.Second}

func httpGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "dbdb-index/1.0 (+https://github.com/tamnd/dbdb-index)")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// ── Sitemap parsing ────────────────────────────────────────────────

var reDBSlug = regexp.MustCompile(`/db/([^/]+)/?$`)

func parseSitemap(data []byte) []Database {
	var set urlSet
	if err := xml.Unmarshal(data, &set); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing sitemap: %v\n", err)
		os.Exit(1)
	}
	var dbs []Database
	for _, u := range set.URLs {
		m := reDBSlug.FindStringSubmatch(u.Loc)
		if m == nil {
			continue
		}
		dbs = append(dbs, Database{
			Name:    slugToName(m[1]),
			Slug:    m[1],
			LastMod: u.LastMod,
		})
	}
	return dbs
}

func slugToName(slug string) string {
	name := strings.ReplaceAll(slug, "-", " ")
	words := strings.Fields(name)
	for i, w := range words {
		runes := []rune(w)
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// ── Page scraping ──────────────────────────────────────────────────

var (
	reH1        = regexp.MustCompile(`(?i)<h1[^>]*>\s*([^<]+?)\s*</h1>`)
	reTitle     = regexp.MustCompile(`(?i)<title>\s*([^<]+?)\s*</title>`)
	reOGDesc    = regexp.MustCompile(`(?i)<meta\s+(?:property="og:description"|name="description")\s+content="([^"]*)"`)
	reOGDescAlt = regexp.MustCompile(`(?i)<meta\s+content="([^"]*)"\s+(?:property="og:description"|name="description")`)

	// Browse filter link patterns: /browse?key=value with display text
	reDataModel = regexp.MustCompile(`/browse\?data-model=([^"&]+)"[^>]*>([^<]+)`)
	reCountry   = regexp.MustCompile(`/browse\?country=([^"&]+)"\s+title="View other systems from ([^"]+)"`)
	reProjType  = regexp.MustCompile(`/browse\?type=([^"&]+)"[^>]*>([^<]+)`)
	reWrittenIn = regexp.MustCompile(`/browse\?programming=([^"&]+)"[^>]*>([^<]+)`)
	reLicense   = regexp.MustCompile(`/browse\?license=([^"&]+)"[^>]*>([^<]+)`)
	reStartYear = regexp.MustCompile(`(?s)Start Year.*?<p class="card-text">\s*(\d{4})\s*</p>`)
)

func scrapeDB(db *Database) error {
	data, err := httpGet(db.pageURL())
	if err != nil {
		return err
	}
	page := string(data)

	// Name — prefer <h1>, fall back to <title>
	if m := reH1.FindStringSubmatch(page); m != nil {
		if name := strings.TrimSpace(html.UnescapeString(m[1])); name != "" {
			db.Name = name
		}
	} else if m := reTitle.FindStringSubmatch(page); m != nil {
		parts := strings.SplitN(html.UnescapeString(m[1]), " - ", 2)
		if name := strings.TrimSpace(parts[0]); name != "" {
			db.Name = name
		}
	}

	// Description from meta tags
	desc := ""
	if m := reOGDesc.FindStringSubmatch(page); m != nil {
		desc = m[1]
	} else if m := reOGDescAlt.FindStringSubmatch(page); m != nil {
		desc = m[1]
	}
	if desc != "" {
		desc = html.UnescapeString(desc)
		desc = strings.Join(strings.Fields(desc), " ")
		db.Desc = desc
	}

	// Data models
	db.Models = extractBrowseValues(reDataModel, page)

	// Country — take the display text from the first match
	if m := reCountry.FindStringSubmatch(page); m != nil {
		db.Country = strings.TrimSpace(html.UnescapeString(m[2]))
	}

	// Start year
	if m := reStartYear.FindStringSubmatch(page); m != nil {
		db.StartYear = m[1]
	}

	// Project types
	db.ProjTypes = extractBrowseValues(reProjType, page)

	// Written in
	db.WrittenIn = extractBrowseValues(reWrittenIn, page)

	// Licenses
	db.Licenses = extractBrowseValues(reLicense, page)

	db.Fetched = true
	return nil
}

// extractBrowseValues returns deduplicated display texts from browse filter links.
func extractBrowseValues(re *regexp.Regexp, page string) []string {
	matches := re.FindAllStringSubmatch(page, -1)
	seen := map[string]bool{}
	var vals []string
	for _, m := range matches {
		text := strings.TrimSpace(html.UnescapeString(m[2]))
		if text != "" && !seen[text] {
			seen[text] = true
			vals = append(vals, text)
		}
	}
	return vals
}

func fetchDetails(dbs []Database, workers int) {
	var (
		mu   sync.Mutex
		wg   sync.WaitGroup
		done int
	)
	sem := make(chan struct{}, workers)

	for i := range dbs {
		if dbs[i].Fetched {
			continue
		}
		wg.Add(1)
		go func(db *Database) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if err := scrapeDB(db); err != nil {
				fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", db.Slug, err)
			}

			mu.Lock()
			done++
			if done%100 == 0 {
				fmt.Fprintf(os.Stderr, "  fetched %d databases...\n", done)
			}
			mu.Unlock()

			time.Sleep(100 * time.Millisecond) // polite rate limiting
		}(&dbs[i])
	}
	wg.Wait()
}

// ── Cache ──────────────────────────────────────────────────────────

func loadCache(path string) []Database {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var dbs []Database
	if err := json.Unmarshal(data, &dbs); err != nil {
		return nil
	}
	return dbs
}

func saveCache(path string, dbs []Database) error {
	data, err := json.MarshalIndent(dbs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// mergeSitemap updates cached entries with fresh sitemap data and adds new entries.
func mergeSitemap(cached, fresh []Database) []Database {
	bySlug := map[string]*Database{}
	for i := range cached {
		bySlug[cached[i].Slug] = &cached[i]
	}
	for _, f := range fresh {
		if existing, ok := bySlug[f.Slug]; ok {
			existing.LastMod = f.LastMod
		} else {
			cached = append(cached, f)
		}
	}
	// Remove cached entries no longer in sitemap
	freshSet := map[string]bool{}
	for _, f := range fresh {
		freshSet[f.Slug] = true
	}
	filtered := cached[:0]
	for _, db := range cached {
		if freshSet[db.Slug] {
			filtered = append(filtered, db)
		}
	}
	return filtered
}

// ── README generation ──────────────────────────────────────────────

var (
	reAnchorClean = regexp.MustCompile(`[^a-z0-9\s-]`)
)

func anchor(label string) string {
	a := strings.ToLower(label)
	a = reAnchorClean.ReplaceAllString(a, "")
	a = strings.TrimSpace(a)
	return strings.ReplaceAll(a, " ", "-")
}

// ── Stats helpers ──────────────────────────────────────────────────

type kv struct {
	k string
	v int
}

func countField(dbs []Database, extract func(Database) []string) []kv {
	m := map[string]int{}
	for _, db := range dbs {
		for _, v := range extract(db) {
			m[v]++
		}
	}
	pairs := make([]kv, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].v != pairs[j].v {
			return pairs[i].v > pairs[j].v
		}
		return pairs[i].k < pairs[j].k
	})
	return pairs
}

// barChart renders a horizontal ASCII bar chart in a fenced code block.
// maxBars limits how many rows to show; barWidth is the max bar length in chars.
func barChart(pairs []kv, maxBars, barWidth int) string {
	if len(pairs) == 0 {
		return ""
	}
	n := len(pairs)
	if n > maxBars {
		n = maxBars
	}
	top := pairs[:n]

	// Find max label width and max value for scaling
	maxLabel, maxVal := 0, 0
	for _, p := range top {
		if len(p.k) > maxLabel {
			maxLabel = len(p.k)
		}
		if p.v > maxVal {
			maxVal = p.v
		}
	}

	var b strings.Builder
	b.WriteString("```\n")
	for _, p := range top {
		w := int(math.Round(float64(p.v) / float64(maxVal) * float64(barWidth)))
		if w < 1 && p.v > 0 {
			w = 1
		}
		fmt.Fprintf(&b, "  %-*s %s %d\n", maxLabel, p.k, strings.Repeat("█", w), p.v)
	}
	b.WriteString("```\n")
	return b.String()
}

// decadeTimeline builds a decade-grouped histogram of start years.
func decadeTimeline(dbs []Database) ([]kv, int, int) {
	years := map[int]int{}
	minY, maxY := 9999, 0
	for _, db := range dbs {
		if db.StartYear == "" {
			continue
		}
		y, err := strconv.Atoi(db.StartYear)
		if err != nil {
			continue
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
		decade := (y / 10) * 10
		years[decade]++
	}
	var pairs []kv
	for d := (minY / 10) * 10; d <= maxY; d += 10 {
		label := fmt.Sprintf("%ds", d)
		pairs = append(pairs, kv{label, years[d]})
	}
	return pairs, minY, maxY
}

func generate(dbs []Database) string {
	sort.Slice(dbs, func(i, j int) bool {
		return strings.ToLower(dbs[i].Name) < strings.ToLower(dbs[j].Name)
	})

	var b strings.Builder

	// ── Header ─────────────────────────────────────────────────────
	fmt.Fprintf(&b, "# Database of Databases — Index\n\n")
	fmt.Fprintf(&b, "A comprehensive index of **%d** database management systems catalogued by "+
		"[Carnegie Mellon University's Database of Databases](https://dbdb.io).\n\n", len(dbs))
	fmt.Fprintf(&b, "> Auto-generated by [dbdb-index](src/). Data sourced from [dbdb.io](https://dbdb.io).\n\n")

	// ── Stats section ──────────────────────────────────────────────
	b.WriteString("## At a Glance\n\n")

	// Collect all stats
	models := countField(dbs, func(d Database) []string { return d.Models })
	countries := countField(dbs, func(d Database) []string {
		if d.Country != "" {
			return []string{d.Country}
		}
		return nil
	})
	projTypes := countField(dbs, func(d Database) []string { return d.ProjTypes })
	languages := countField(dbs, func(d Database) []string { return d.WrittenIn })
	licenses := countField(dbs, func(d Database) []string { return d.Licenses })
	decades, minY, maxY := decadeTimeline(dbs)

	// Count open source vs commercial
	openSource, commercial := 0, 0
	for _, p := range projTypes {
		if strings.Contains(strings.ToLower(p.k), "open source") {
			openSource = p.v
		}
		if strings.Contains(strings.ToLower(p.k), "commercial") {
			commercial = p.v
		}
	}

	// Prose summary
	fmt.Fprintf(&b, "The database landscape spans **%d years** of innovation (from %d to %d), "+
		"across **%d countries**. ", maxY-minY, minY, maxY, len(countries))
	if len(models) > 0 {
		fmt.Fprintf(&b, "The most popular data model is **%s** (%d systems), followed by "+
			"**%s** (%d) and **%s** (%d). ", models[0].k, models[0].v, models[1].k, models[1].v, models[2].k, models[2].v)
	}
	if openSource > 0 && commercial > 0 {
		fmt.Fprintf(&b, "Open source projects (%d) outnumber commercial ones (%d) by %.1fx. ",
			openSource, commercial, float64(openSource)/float64(commercial))
	}
	if len(languages) > 0 {
		fmt.Fprintf(&b, "**%s** is the most common implementation language (%d systems), with **%s** (%d) and **%s** (%d) close behind.",
			languages[0].k, languages[0].v, languages[1].k, languages[1].v, languages[2].k, languages[2].v)
	}
	b.WriteString("\n\n")

	// Data models chart
	b.WriteString("### By Data Model\n\n")
	b.WriteString(barChart(models, 15, 40))
	b.WriteByte('\n')

	// Timeline chart
	b.WriteString("### By Decade\n\n")
	fmt.Fprintf(&b, "Database systems through the decades — the explosion of new projects in the 2010s "+
		"reflects the rise of NoSQL, NewSQL, and cloud-native databases.\n\n")
	b.WriteString(barChart(decades, 20, 40))
	b.WriteByte('\n')

	// Top countries
	b.WriteString("### By Country of Origin\n\n")
	b.WriteString(barChart(countries, 15, 40))
	b.WriteByte('\n')

	// Implementation languages
	b.WriteString("### By Implementation Language\n\n")
	b.WriteString(barChart(languages, 15, 40))
	b.WriteByte('\n')

	// Licenses
	b.WriteString("### By License\n\n")
	b.WriteString(barChart(licenses, 12, 40))
	b.WriteByte('\n')

	// Project types
	b.WriteString("### By Project Type\n\n")
	b.WriteString(barChart(projTypes, 10, 40))
	b.WriteByte('\n')

	// ── Table of contents ──────────────────────────────────────────
	groups := groupByLetter(dbs)
	letters := sortedKeys(groups)

	// Table of contents
	b.WriteString("## Contents\n\n")
	for _, letter := range letters {
		n := len(groups[letter])
		fmt.Fprintf(&b, "[%s](#%s)(%d)", letter, anchor(letter), n)
		b.WriteString(" · ")
	}
	// trim trailing separator
	s := b.String()
	s = strings.TrimSuffix(s, " · ")
	b.Reset()
	b.WriteString(s)
	b.WriteString("\n\n")

	// Sections
	for _, letter := range letters {
		entries := groups[letter]
		fmt.Fprintf(&b, "## %s\n\n", letter)
		fmt.Fprintln(&b, "| Database | Description | Data Models | Country | Year | Type | Written In | License | Modified |")
		fmt.Fprintln(&b, "|----------|-------------|-------------|---------|------|------|------------|---------|----------|")
		for _, db := range entries {
			desc := db.Desc
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			desc = strings.ReplaceAll(desc, "|", "\\|")
			desc = strings.ReplaceAll(desc, "\n", " ")
			lastmod := ""
			if db.LastMod != "" {
				lastmod = db.LastMod[:min(10, len(db.LastMod))]
			}
			fmt.Fprintf(&b, "| [%s](%s) | %s | %s | %s | %s | %s | %s | %s | %s |\n",
				db.Name, db.pageURL(), desc,
				strings.Join(db.Models, ", "),
				db.Country, db.StartYear,
				strings.Join(db.ProjTypes, ", "),
				strings.Join(db.WrittenIn, ", "),
				strings.Join(db.Licenses, ", "),
				lastmod)
		}
		b.WriteByte('\n')
	}

	// Footer
	b.WriteString("---\n\n")
	fmt.Fprintf(&b, "*Last generated: %s · Source: [dbdb.io](https://dbdb.io) · "+
		"Tool: [dbdb-index](src/)*\n", time.Now().UTC().Format("2006-01-02"))

	return b.String()
}

func groupByLetter(dbs []Database) map[string][]Database {
	groups := map[string][]Database{}
	for _, db := range dbs {
		r := []rune(db.Name)[0]
		letter := strings.ToUpper(string(r))
		if r < 'A' || (r > 'Z' && r < 'a') || r > 'z' {
			letter = "#"
		}
		groups[letter] = append(groups[letter], db)
	}
	return groups
}

func sortedKeys(m map[string][]Database) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// Move "#" to the front if present
	for i, k := range keys {
		if k == "#" {
			keys = append([]string{k}, append(keys[:i], keys[i+1:]...)...)
			break
		}
	}
	return keys
}

// ── CLI ────────────────────────────────────────────────────────────

func main() {
	output := flag.String("o", "README.md", "output file path")
	cache := flag.String("cache", "databases.json", "cache file for fetched metadata")
	workers := flag.Int("workers", 20, "concurrent fetch workers")
	noFetch := flag.Bool("no-fetch", false, "skip page fetching; generate from cache/sitemap names only")
	flag.Parse()

	// 1. Fetch sitemap
	fmt.Fprintln(os.Stderr, "Fetching sitemap...")
	sitemapData, err := httpGet(baseURL + "/sitemap.xml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching sitemap: %v\n", err)
		os.Exit(1)
	}
	fresh := parseSitemap(sitemapData)
	fmt.Fprintf(os.Stderr, "Found %d databases in sitemap\n", len(fresh))

	// 2. Merge with cache
	dbs := loadCache(*cache)
	if len(dbs) > 0 {
		fmt.Fprintf(os.Stderr, "Loaded %d cached entries\n", len(dbs))
		dbs = mergeSitemap(dbs, fresh)
	} else {
		dbs = fresh
	}

	// 3. Fetch individual pages for metadata
	if !*noFetch {
		unfetched := 0
		for _, db := range dbs {
			if !db.Fetched {
				unfetched++
			}
		}
		if unfetched > 0 {
			fmt.Fprintf(os.Stderr, "Fetching details for %d databases (%d workers)...\n", unfetched, *workers)
			fetchDetails(dbs, *workers)
			if err := saveCache(*cache, dbs); err != nil {
				fmt.Fprintf(os.Stderr, "warning: cache save failed: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "Cache saved to %s\n", *cache)
			}
		} else {
			fmt.Fprintln(os.Stderr, "All databases already cached")
		}
	}

	// 4. Generate README
	readme := generate(dbs)
	if err := os.WriteFile(*output, []byte(readme), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Written %s: %d databases, %d lines\n",
		*output, len(dbs), strings.Count(readme, "\n")+1)
}
