package sighting

// Supplementary publication evidence is deliberately separate from the month
// anchors. It never writes or changes the historical reports JSON.
import (
	"html"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type Evidence struct {
	Date string
	Title string
	URL string
	Status string
	Method string
}

var evidenceStore = struct {
	sync.RWMutex
	items map[string]Evidence
	state string
}{items: make(map[string]Evidence), state: "pending"}

var (
	postPattern = regexp.MustCompile(`data-post="([a-zA-Z0-9_]+/[0-9]+)"`)
	postDatePattern = regexp.MustCompile(`datetime="([^"]+)"`)
	previewTitlePattern = regexp.MustCompile(`(?s)class="tgme_widget_message_link_preview_title[^"]*"[^>]*>(.*?)</div>`)
	messageTextPattern = regexp.MustCompile(`(?s)class="tgme_widget_message_text[^"]*"[^>]*>(.*?)</div>`)
	pageTitlePattern = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	positiveUpdatePattern = regexp.MustCompile(`(?i)updated\s+list\s+of\s+positive\s+observations\s+from\s+Israel\s+on\s+`)
)

func plainText(body string) string {
	return strings.Join(strings.Fields(html.UnescapeString(tagPattern.ReplaceAllString(body, " "))), " ")
}

func lunarPublication(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "crescent") || strings.Contains(lower, "moon observation") || strings.Contains(lower, "biblical month")
}

func telegramEvidence(body string) []Evidence {
	positions := postPattern.FindAllStringSubmatchIndex(body, -1)
	items := make([]Evidence, 0)
	for i, position := range positions {
		end := len(body)
		if i+1 < len(positions) { end = positions[i+1][0] }
		block := body[position[0]:end]
		datetime := postDatePattern.FindStringSubmatch(block)
		if len(datetime) != 2 { continue }
		published, err := time.Parse(time.RFC3339, datetime[1])
		if err != nil { continue }
		title := ""
		for _, pattern := range []*regexp.Regexp{previewTitlePattern, messageTextPattern} {
			if match := pattern.FindStringSubmatch(block); len(match) == 2 { title = plainText(match[1]); if title != "" { break } }
		}
		if !lunarPublication(title) { continue }
		link := "https://t.me/"+body[position[2]:position[3]]
		// A preview and its linked article are one publication, not two witnesses.
		if match := devorahLinkPattern.FindStringSubmatch(block); len(match) == 2 {
			link = strings.Split(html.UnescapeString(match[1]), "?")[0]
		}
		status := "related"
		if strings.Contains(strings.ToLower(title), "upcoming") { status = "pending" }
		items = append(items, Evidence{Date: published.Format("2006-01-02"), Title: title, URL: link, Status: status, Method: "publication_title"})
	}
	return items
}

func rememberEvidence(item Evidence) {
	if len([]rune(item.Title)) > 220 { item.Title = string([]rune(item.Title)[:220])+"…" }
	evidenceStore.Lock()
	defer evidenceStore.Unlock()
	previous, exists := evidenceStore.items[item.URL]
	if !exists || previous.Method == "publication_title" || item.Method != "publication_title" { evidenceStore.items[item.URL] = item }
	// Bounded, recent in-memory index; no historical rewrite or per-click download.
	if len(evidenceStore.items) > 64 {
		var oldest string
		for key, value := range evidenceStore.items { if oldest == "" || value.Date < evidenceStore.items[oldest].Date { oldest = key } }
		delete(evidenceStore.items, oldest)
	}
}

func recordPublicationEvidence(body, link string) {
	titleMatch := pageTitlePattern.FindStringSubmatch(body)
	if len(titleMatch) != 2 { return }
	title := plainText(titleMatch[1])
	if !lunarPublication(title) { return }
	report, err := parseDevorahMonthReport(body, link, time.Now().UTC())
	if err == nil {
		status := "confirmed"
		if report.ConfirmationMethod == "day_30_declaration" { status = "declared" }
		rememberEvidence(Evidence{Date: report.Date, Title: title, URL: link, Status: status, Method: report.ConfirmationMethod})
		return
	}
	// Updates can list dated positive observations without repeating a declaration.
	plain := plainText(body)
	if anchor := positiveUpdatePattern.FindStringIndex(plain); anchor != nil {
		if match := englishDatePattern.FindStringSubmatch(plain[anchor[1]:]); len(match) == 2 {
			if date, err := parseEnglishDate(match[1]); err == nil {
				rememberEvidence(Evidence{Date: date.Format("2006-01-02"), Title: title, URL: link, Status: "confirmed", Method: "israel_eyewitness_update"})
			}
		}
	}
}

func setEvidenceState(state string) { evidenceStore.Lock(); evidenceStore.state = state; evidenceStore.Unlock() }

// EvidenceFor uses a narrow observation window around an already validated
// boundary. Nearby publications support the period, never confirm another date.
func EvidenceFor(date time.Time) ([]Evidence, string) {
	evidenceStore.RLock()
	defer evidenceStore.RUnlock()
	items := make([]Evidence, 0)
	start, end := date.AddDate(0, 0, -7).Format("2006-01-02"), date.AddDate(0, 0, 7).Format("2006-01-02")
	for _, item := range evidenceStore.items {
		if item.Date >= start && item.Date <= end { items = append(items, item) }
	}
	sort.Slice(items, func(i, j int) bool { return items[i].URL < items[j].URL })
	return items, evidenceStore.state
}
