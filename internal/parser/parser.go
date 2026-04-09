package parser

import (
	"cherkasyoblenergo-api/internal/models"
	"cherkasyoblenergo-api/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	_ "time/tzdata"

	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var kievLocation *time.Location

var allowedDomains = []string{
	"cherkasyoblenergo.com",
	"gita.cherkasyoblenergo.com",
}

var privateIPBlocks = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
}

func init() {
	var err error
	kievLocation, err = time.LoadLocation("Europe/Kiev")
	if err != nil {
		log.Printf("Failed to load Europe/Kiev timezone, using UTC: %v", err)
		kievLocation = time.UTC
	}
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, cidr := range privateIPBlocks {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func ValidateNewsURL(rawURL string) error {
	if rawURL == "" {
		return errors.New("URL cannot be empty")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "https" {
		return errors.New("only HTTPS URLs are allowed")
	}

	host := parsedURL.Hostname()
	if host == "" {
		return errors.New("URL must have a hostname")
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("failed to resolve hostname: %w", err)
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return errors.New("private IP addresses are not allowed")
		}
	}

	allowed := false
	for _, domain := range allowedDomains {
		if host == domain || strings.HasSuffix(host, "."+domain) {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("domain is not in the allowlist")
	}

	return nil
}

func createSecureHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, fmt.Errorf("invalid address: %w", err)
			}

			ips, err := net.LookupIP(host)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve hostname during dial: %w", err)
			}

			for _, ip := range ips {
				if isPrivateIP(ip) {
					return nil, fmt.Errorf("private IP addresses are not allowed during connection")
				}
			}

			return dialer.DialContext(ctx, network, addr)
		},
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}

var scheduleKeywords = []string{
	"оновлені графіки",
	"графіки погодинних",
	"графіки відключень",
}

type NewsResponse struct {
	NewsList []struct {
		ID       int    `json:"id"`
		Date     string `json:"date"`
		Title    string `json:"title"`
		HtmlBody string `json:"htmlBody"`
	} `json:"newsList"`
}

type scheduleNews struct {
	ID       int
	Date     time.Time
	Title    string
	HtmlBody string
}

type syncResult string

const (
	syncResultCreated       syncResult = "created"
	syncResultUpdated       syncResult = "updated"
	syncResultUnchanged     syncResult = "unchanged"
	syncResultSkippedNoData syncResult = "skipped_no_schedule"
	syncResultDatabaseError syncResult = "database_error"
)

func StartCron(db *gorm.DB, newsURL string) *cron.Cron {
	if newsURL == "" {
		log.Fatal("NEWS_URL environment variable is required")
	}

	intervalStr := os.Getenv("PARSING_INTERVAL_MINUTES")
	interval := 5
	if intervalStr != "" {
		if parsed, err := strconv.Atoi(intervalStr); err == nil && parsed > 0 {
			interval = parsed
		}
	}
	cronSchedule := "@every " + strconv.Itoa(interval) + "m"

	c := cron.New()
	c.AddFunc(cronSchedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		FetchAndStoreNews(ctx, db, newsURL)
	})

	c.Start()
	return c
}

var isParsing int32

func FetchAndStoreNews(ctx context.Context, db *gorm.DB, newsURL string) {
	if !atomic.CompareAndSwapInt32(&isParsing, 0, 1) {
		log.Println("Parsing job already running, skipping")
		return
	}
	defer atomic.StoreInt32(&isParsing, 0)

	if err := ValidateNewsURL(newsURL); err != nil {
		log.Printf("URL validation failed: %v", err)
		return
	}

	log.Println("Starting news parsing")

	req, err := http.NewRequestWithContext(ctx, "GET", newsURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	client := createSecureHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return
	}

	var newsResp NewsResponse
	if err = json.Unmarshal(body, &newsResp); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		return
	}
	log.Printf("Fetched %d news items", len(newsResp.NewsList))

	var filteredNews []scheduleNews
	for _, news := range newsResp.NewsList {
		hasScheduleKeywords := containsScheduleKeywords(news.Title)
		hasSchedulePatterns := containsSchedulePatterns(news.HtmlBody)

		if !hasScheduleKeywords && !hasSchedulePatterns {
			continue
		}
		parsedDate, err := time.ParseInLocation("02.01.2006 15:04", news.Date, kievLocation)
		if err != nil {
			continue
		}
		filteredNews = append(filteredNews, scheduleNews{
			ID:       news.ID,
			Date:     parsedDate,
			Title:    news.Title,
			HtmlBody: news.HtmlBody,
		})
	}

	sort.Slice(filteredNews, func(i, j int) bool {
		if filteredNews[i].Date.Equal(filteredNews[j].Date) {
			return filteredNews[i].ID < filteredNews[j].ID
		}
		return filteredNews[i].Date.Before(filteredNews[j].Date)
	})
	log.Printf("Filtered to %d relevant news items", len(filteredNews))

	createdCount := 0
	updatedCount := 0
	unchangedCount := 0
	skippedCount := 0
	errorCount := 0
	for _, news := range filteredNews {
		result, err := syncScheduleRecord(db, news)
		if err != nil {
			log.Printf("Failed to sync news ID %d: %v", news.ID, err)
			errorCount++
			continue
		}

		switch result {
		case syncResultCreated:
			createdCount++
			log.Printf("Created schedule data from news ID: %d", news.ID)
		case syncResultUpdated:
			updatedCount++
			log.Printf("Updated schedule data from news ID: %d", news.ID)
		case syncResultUnchanged:
			unchangedCount++
			log.Printf("Schedule data unchanged for news ID: %d", news.ID)
		case syncResultSkippedNoData:
			skippedCount++
			log.Printf("Skipping news ID %d because no schedule data was parsed", news.ID)
		}
	}
	log.Printf(
		"Sync summary: created=%d updated=%d unchanged=%d skipped=%d errors=%d",
		createdCount,
		updatedCount,
		unchangedCount,
		skippedCount,
		errorCount,
	)
}

func buildScheduleFromNews(news scheduleNews) (models.Schedule, bool) {
	sch := parseScheduleData(news.HtmlBody)
	if !hasScheduleData(sch) {
		return models.Schedule{}, false
	}

	sch.NewsID = news.ID
	sch.Title = news.Title
	sch.Date = news.Date
	sch.ScheduleDate = utils.ExtractScheduleDateFromTitle(news.Title)
	return sch, true
}

func scheduleChanged(existing, incoming models.Schedule) bool {
	return existing.NewsID != incoming.NewsID ||
		existing.Title != incoming.Title ||
		!existing.Date.Equal(incoming.Date) ||
		existing.ScheduleDate != incoming.ScheduleDate ||
		existing.Col1_1 != incoming.Col1_1 ||
		existing.Col1_2 != incoming.Col1_2 ||
		existing.Col2_1 != incoming.Col2_1 ||
		existing.Col2_2 != incoming.Col2_2 ||
		existing.Col3_1 != incoming.Col3_1 ||
		existing.Col3_2 != incoming.Col3_2 ||
		existing.Col4_1 != incoming.Col4_1 ||
		existing.Col4_2 != incoming.Col4_2 ||
		existing.Col5_1 != incoming.Col5_1 ||
		existing.Col5_2 != incoming.Col5_2 ||
		existing.Col6_1 != incoming.Col6_1 ||
		existing.Col6_2 != incoming.Col6_2
}

func syncScheduleRecord(db *gorm.DB, news scheduleNews) (syncResult, error) {
	incoming, ok := buildScheduleFromNews(news)
	if !ok {
		return syncResultSkippedNoData, nil
	}

	var existing models.Schedule
	err := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
		Where("news_id = ?", news.ID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&incoming).Error; err != nil {
			return syncResultDatabaseError, err
		}
		return syncResultCreated, nil
	}
	if err != nil {
		return syncResultDatabaseError, err
	}

	if !scheduleChanged(existing, incoming) {
		return syncResultUnchanged, nil
	}

	updates := map[string]interface{}{
		"title":         incoming.Title,
		"date":          incoming.Date,
		"schedule_date": incoming.ScheduleDate,
		"1_1":           incoming.Col1_1,
		"1_2":           incoming.Col1_2,
		"2_1":           incoming.Col2_1,
		"2_2":           incoming.Col2_2,
		"3_1":           incoming.Col3_1,
		"3_2":           incoming.Col3_2,
		"4_1":           incoming.Col4_1,
		"4_2":           incoming.Col4_2,
		"5_1":           incoming.Col5_1,
		"5_2":           incoming.Col5_2,
		"6_1":           incoming.Col6_1,
		"6_2":           incoming.Col6_2,
	}

	if err := db.Model(&existing).Updates(updates).Error; err != nil {
		return syncResultDatabaseError, err
	}

	return syncResultUpdated, nil
}

func containsScheduleKeywords(title string) bool {
	titleLower := strings.ToLower(title)
	for _, kw := range scheduleKeywords {
		if strings.Contains(titleLower, kw) {
			return true
		}
	}
	return false
}

func containsSchedulePatterns(htmlBody string) bool {
	re := regexp.MustCompile(`\b[1-6]\.[1-2]:?\s+\d{1,2}:\d{2}`)
	return re.MatchString(htmlBody)
}

func normalizeSpaces(text string) string {
	text = strings.ReplaceAll(text, "\u00a0", " ")
	text = strings.ReplaceAll(text, "\u202f", " ")
	text = strings.TrimSpace(text)
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}
	return text
}

func normalizeTimeRanges(value string) string {
	value = normalizeSpaces(value)
	value = strings.TrimRight(value, ",.;")
	return strings.TrimSpace(value)
}

func parseScheduleFromParagraphs(htmlBody string) (models.Schedule, bool) {
	var data models.Schedule
	found := false

	re := regexp.MustCompile(`^(\d)\.(\d)\.?:?\s*(.+)`)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		log.Printf("Failed to parse HTML in parseScheduleFromParagraphs: %v", err)
		return data, false
	}

	paragraphCount := 0
	seen := make(map[string]struct{})

	doc.Find("p, div").Each(func(i int, s *goquery.Selection) {
		text := normalizeSpaces(s.Text())

		if text == "" {
			return
		}

		if _, exists := seen[text]; exists {
			return
		}
		seen[text] = struct{}{}
		paragraphCount++

		matches := re.FindStringSubmatch(text)
		if matches == nil {
			return
		}

		mainQueue, err1 := strconv.Atoi(matches[1])
		subQueue, err2 := strconv.Atoi(matches[2])
		timeRanges := normalizeTimeRanges(matches[3])

		if err1 != nil || err2 != nil || mainQueue < 1 || mainQueue > 6 || subQueue < 1 || subQueue > 2 {
			log.Printf("Invalid queue values: mainQueue=%d, subQueue=%d", mainQueue, subQueue)
			return
		}

		log.Printf("Found schedule: %d.%d = %s", mainQueue, subQueue, timeRanges)
		setQueueValue(&data, mainQueue, subQueue, timeRanges)
		found = true
	})

	log.Printf("parseScheduleFromParagraphs: checked %d paragraphs, found=%v", paragraphCount, found)
	return data, found
}

func hasScheduleData(s models.Schedule) bool {
	cols := []string{
		s.Col1_1, s.Col1_2, s.Col2_1, s.Col2_2, s.Col3_1, s.Col3_2,
		s.Col4_1, s.Col4_2, s.Col5_1, s.Col5_2, s.Col6_1, s.Col6_2,
	}
	for _, c := range cols {
		if strings.TrimSpace(c) != "" {
			return true
		}
	}
	return false
}

func parseScheduleData(htmlBody string) models.Schedule {
	var data models.Schedule

	if parsedData, found := parseScheduleFromParagraphs(htmlBody); found {
		return parsedData
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err == nil {
		foundTableData := false
		doc.Find("table tr").Each(func(i int, tr *goquery.Selection) {
			if tr.Find("th").Length() > 0 {
				return
			}
			queueStr := strings.TrimSpace(tr.Find("td").First().Text())
			timeStr := strings.TrimSpace(tr.Find("td").Last().Text())
			parts := strings.Split(queueStr, ".")
			if len(parts) != 2 {
				return
			}
			mainQueue, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return
			}
			var subQueue int
			switch strings.TrimSpace(parts[1]) {
			case "І", "l", "I":
				subQueue = 1
			case "ІІ", "ll", "II":
				subQueue = 2
			default:
				return
			}
			setQueueValue(&data, mainQueue, subQueue, timeStr)
			foundTableData = true
		})

		if foundTableData {
			return data
		}
	}

	if strings.Contains(strings.ToLower(htmlBody), "скасовано") {
		for i := 1; i <= 6; i++ {
			for j := 1; j <= 2; j++ {
				setQueueValue(&data, i, j, "скасовано")
			}
		}
		return data
	}

	return data
}

func setQueueValue(data *models.Schedule, queue, subQueue int, value string) {
	switch queue {
	case 1:
		if subQueue == 1 {
			data.Col1_1 = value
		} else {
			data.Col1_2 = value
		}
	case 2:
		if subQueue == 1 {
			data.Col2_1 = value
		} else {
			data.Col2_2 = value
		}
	case 3:
		if subQueue == 1 {
			data.Col3_1 = value
		} else {
			data.Col3_2 = value
		}
	case 4:
		if subQueue == 1 {
			data.Col4_1 = value
		} else {
			data.Col4_2 = value
		}
	case 5:
		if subQueue == 1 {
			data.Col5_1 = value
		} else {
			data.Col5_2 = value
		}
	case 6:
		if subQueue == 1 {
			data.Col6_1 = value
		} else {
			data.Col6_2 = value
		}
	}
}
