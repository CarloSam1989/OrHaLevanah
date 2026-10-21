package sighting

import (
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const devorahSourceName = "Devorah's Date Tree"

var (
	devorahLinkPattern = regexp.MustCompile(`(?i)href=["'](https://mailchi\.mp/devorahsdatetree/[^"'#]+)`)
	monthNumberPattern = regexp.MustCompile(`(?i)\b(1[0-3]|[1-9])(?:st|nd|rd|th)\s+Biblical\s+Month\b`)
	englishDatePattern = regexp.MustCompile(`(?i)\b(?:Monday,?\s+|Tuesday,?\s+|Wednesday,?\s+|Thursday,?\s+|Friday,?\s+|Saturday,?\s+|Sunday,?\s+)?((?:January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2},?\s+\d{4}|\d{1,2}\s+(?:January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{4})\b`)
)

func refreshFromDevorah(client *http.Client, req *http.Request) error {
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Devorah Telegram source returned HTTP %d", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 3<<20))
	if err != nil {
		return err
	}
	links := devorahReportLinks(string(body))
	for _, item := range telegramEvidence(string(body)) { rememberEvidence(item) }
	setEvidenceState("complete")
	if len(links) == 0 {
		return errors.New("no Devorah crescent report links found")
	}

	// The public preview contains several months. Backfill the latest reports so
	// the whole biblical year and its feasts share the same observed anchors.
	if len(links) > 8 {
		links = links[len(links)-8:]
	}
	var saved int
	var lastErr error
	for _, link := range links {
		report, err := fetchDevorahReport(client, link)
		if err != nil {
			lastErr = err
			continue
		}
		if err := saveReport(report); err != nil {
			lastErr = err
			continue
		}
		saved++
	}
	if saved == 0 {
		if lastErr != nil {
			return lastErr
		}
		return errors.New("no authoritative Devorah month declaration found")
	}
	return nil
}

func devorahReportLinks(body string) []string {
	matches := devorahLinkPattern.FindAllStringSubmatch(body, -1)
	seen := make(map[string]bool)
	links := make([]string, 0, len(matches))
	for _, match := range matches {
		link := html.UnescapeString(match[1])
		parsed, err := url.Parse(link)
		if err != nil || !strings.EqualFold(parsed.Host, "mailchi.mp") || !strings.HasPrefix(parsed.Path, "/devorahsdatetree/") {
			continue
		}
		lower := strings.ToLower(parsed.Path)
		if strings.Contains(lower, "upcoming") || !strings.Contains(lower, "report") {
			continue
		}
		parsed.RawQuery = ""
		link = parsed.String()
		if !seen[link] {
			seen[link] = true
			links = append(links, link)
		}
	}
	return links
}

func fetchDevorahReport(client *http.Client, link string) (Report, error) {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return Report{}, err
	}
	req.Header.Set("User-Agent", "Or-HaLevanah/1.0 Devorah-report-checker")
	response, err := client.Do(req)
	if err != nil {
		setEvidenceState("partial")
		return Report{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		setEvidenceState("partial")
		return Report{}, fmt.Errorf("Devorah report returned HTTP %d", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 3<<20))
	if err != nil {
		setEvidenceState("partial")
		return Report{}, err
	}
	recordPublicationEvidence(string(body), link)
	return parseDevorahMonthReport(string(body), link, time.Now().UTC())
}

func parseDevorahMonthReport(body, link string, checkedAt time.Time) (Report, error) {
	plain := html.UnescapeString(tagPattern.ReplaceAllString(body, " "))
	plain = strings.Join(strings.Fields(plain), " ")
	monthMatch := monthNumberPattern.FindStringSubmatch(plain)
	if len(monthMatch) != 2 {
		return Report{}, errors.New("Devorah report has no biblical month number")
	}
	month, _ := strconv.Atoi(monthMatch[1])

	lower := strings.ToLower(plain)
	hasMonthDeclaration := strings.Contains(lower, "biblical month has begun") || strings.Contains(lower, "biblical month began")
	hasIsraelContext := strings.Contains(lower, "from israel") || strings.Contains(lower, "land of israel")
	positiveAt := firstPhraseIndex(
		lower,
		"sighted from israel on",
		"spotted from israel on",
		"sighted from the land of israel on",
		"spotted from the land of israel on",
		"crescent moon was sighted on",
		"crescent moon was spotted on",
	)
	method := "israel_eyewitness"
	dateSearch := ""
	if positiveAt >= 0 && hasMonthDeclaration && hasIsraelContext {
		dateSearch = plain[positiveAt:]
	} else {
		defaultAt := firstPhraseIndex(lower, "will begin on", "began on")
		isDayThirty := strings.Contains(lower, "by default") || strings.Contains(lower, "30 days") || strings.Contains(lower, "30th day") || strings.Contains(lower, "thirty days")
		if defaultAt < 0 || !isDayThirty {
			return Report{}, errors.New("Devorah report does not confirm an Israel sighting or day-30 month start")
		}
		method = "day_30_declaration"
		dateSearch = plain[defaultAt:]
	}

	dateMatch := englishDatePattern.FindStringSubmatch(dateSearch)
	if len(dateMatch) != 2 {
		return Report{}, errors.New("Devorah report has no confirmed start date")
	}
	date, err := parseEnglishDate(dateMatch[1])
	if err != nil {
		return Report{}, err
	}
	return Report{
		Date:               date.Format("2006-01-02"),
		Country:            "Israel",
		Location:           "Israel",
		Status:             "confirmed",
		SourceName:         devorahSourceName,
		SourceURL:          link,
		BiblicalMonth:      month,
		ConfirmationMethod: method,
		CheckedAt:          checkedAt,
	}, nil
}

func firstPhraseIndex(value string, phrases ...string) int {
	best := -1
	for _, phrase := range phrases {
		if index := strings.Index(value, phrase); index >= 0 && (best < 0 || index < best) {
			best = index
		}
	}
	return best
}

func parseEnglishDate(value string) (time.Time, error) {
	value = strings.ReplaceAll(strings.TrimSpace(value), ",", "")
	for _, layout := range []string{"January 2 2006", "2 January 2006"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid English report date %q", value)
}
