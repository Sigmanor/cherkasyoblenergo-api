package parser

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"cherkasyoblenergo-api/internal/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

func StartCron(db *gorm.DB, newsURL string) {
	if newsURL == "" {
		log.Fatal("NEWS_URL environment variable is required")
	}

	intervalStr := os.Getenv("PARSING_INTERVAL_MINUTES")
	interval := 10
	if intervalStr != "" {
		if parsed, err := strconv.Atoi(intervalStr); err == nil && parsed > 0 {
			interval = parsed
		}
	}
	cronSchedule := "@every " + strconv.Itoa(interval) + "m"

	c := cron.New()
	c.AddFunc(cronSchedule, func() { FetchAndStoreNews(db, newsURL) })

	c.Start()
}

func FetchAndStoreNews(db *gorm.DB, newsURL string) {
	log.Println("Starting news parsing")
	resp, err := http.Get(newsURL)
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
		// Check if title contains schedule keywords OR if content contains schedule patterns
		hasScheduleKeywords := containsScheduleKeywords(news.Title)
		hasSchedulePatterns := containsSchedulePatterns(news.HtmlBody)
		
		if !hasScheduleKeywords && !hasSchedulePatterns {
			continue
		}
		parsedDate, err := time.Parse("02.01.2006 15:04", news.Date)
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
		return filteredNews[i].Date.Before(filteredNews[j].Date)
	})
	log.Printf("Filtered to %d relevant news items", len(filteredNews))

	savedCount := 0
	for _, news := range filteredNews {
		var existing models.Schedule
		err = db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Where("news_id = ?", news.ID).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sch := parseScheduleData(news.HtmlBody)
			sch.NewsID = news.ID
			sch.Title = news.Title
			sch.Date = news.Date
			if err = db.Create(&sch).Error; err != nil {
				log.Printf("Failed to save data to DB: %v", err)
			} else {
				log.Printf("Successfully saved schedule data from news ID: %d", news.ID)
				savedCount++
			}
		} else if err != nil {
			log.Printf("Database error when checking news ID %d: %v", news.ID, err)
		}
	}
	log.Printf("Saved %d new schedule data items", savedCount)
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
	// Check if the HTML body contains schedule patterns like "1.1 <time>" or "2.2 <time>"
	// This regex looks for patterns like "1.1 10:00" or "6.2 15:00-17:00"
	re := regexp.MustCompile(`\b[1-6]\.[1-2]\s+\d{1,2}:\d{2}`)
	return re.MatchString(htmlBody)
}

func parseScheduleFromParagraphs(htmlBody string) (models.Schedule, bool) {
	var data models.Schedule
	found := false
	
	re := regexp.MustCompile(`^(\d)\.(\d)\s+(.+)`)
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return data, false
	}
	
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text == "" {
			return
		}
		
		matches := re.FindStringSubmatch(text)
		if matches == nil {
			return
		}
		
		mainQueue, err1 := strconv.Atoi(matches[1])
		subQueue, err2 := strconv.Atoi(matches[2])
		timeRanges := strings.TrimSpace(matches[3])
		
		if err1 != nil || err2 != nil || mainQueue < 1 || mainQueue > 6 || subQueue < 1 || subQueue > 2 {
			return
		}
		
		setQueueValue(&data, mainQueue, subQueue, timeRanges)
		found = true
	})
	
	return data, found
}

func parseScheduleData(htmlBody string) models.Schedule {
	var data models.Schedule
	if strings.Contains(strings.ToLower(htmlBody), "скасовано") {
		for i := 1; i <= 6; i++ {
			for j := 1; j <= 2; j++ {
				setQueueValue(&data, i, j, "скасовано")
			}
		}
		return data
	}
	
	if parsedData, found := parseScheduleFromParagraphs(htmlBody); found {
		return parsedData
	}
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return data
	}
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
	})
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
