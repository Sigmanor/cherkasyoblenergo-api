package handlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cherkasyoblenergo-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Schedule struct {
	ID           int64     `json:"id"`
	NewsID       int       `json:"news_id"`
	Title        string    `json:"title"`
	Date         time.Time `json:"date"`
	ScheduleDate string    `json:"schedule_date"`
	OneOne       string    `gorm:"column:1_1" json:"1_1"`
	OneTwo       string    `gorm:"column:1_2" json:"1_2"`
	TwoOne       string    `gorm:"column:2_1" json:"2_1"`
	TwoTwo       string    `gorm:"column:2_2" json:"2_2"`
	ThreeOne     string    `gorm:"column:3_1" json:"3_1"`
	ThreeTwo     string    `gorm:"column:3_2" json:"3_2"`
	FourOne      string    `gorm:"column:4_1" json:"4_1"`
	FourTwo      string    `gorm:"column:4_2" json:"4_2"`
	FiveOne      string    `gorm:"column:5_1" json:"5_1"`
	FiveTwo      string    `gorm:"column:5_2" json:"5_2"`
	SixOne       string    `gorm:"column:6_1" json:"6_1"`
	SixTwo       string    `gorm:"column:6_2" json:"6_2"`
}

type ScheduleFilter struct {
	Option string `json:"option"`
	Date   string `json:"date"`
	Limit  int    `json:"limit"`
	Queue  string `json:"queue"`
}

func parseAndValidateQueues(queueStr string) ([]string, error) {
	if strings.TrimSpace(queueStr) == "" {
		return []string{}, nil
	}

	tokens := strings.Split(queueStr, ",")

	queuePattern := regexp.MustCompile(`^[1-6]_[1-2]$`)

	seen := make(map[string]bool)
	result := []string{}

	for _, token := range tokens {
		queue := strings.TrimSpace(token)

		if !queuePattern.MatchString(queue) {
			return nil, fmt.Errorf("invalid queue value: '%s'. Each queue must match format X_Y where X is 1-6 and Y is 1-2", queue)
		}

		if !seen[queue] {
			seen[queue] = true
			result = append(result, queue)
		}
	}

	return result, nil
}

func getQueueValue(schedule *Schedule, queueName string) string {
	switch queueName {
	case "1_1":
		return schedule.OneOne
	case "1_2":
		return schedule.OneTwo
	case "2_1":
		return schedule.TwoOne
	case "2_2":
		return schedule.TwoTwo
	case "3_1":
		return schedule.ThreeOne
	case "3_2":
		return schedule.ThreeTwo
	case "4_1":
		return schedule.FourOne
	case "4_2":
		return schedule.FourTwo
	case "5_1":
		return schedule.FiveOne
	case "5_2":
		return schedule.FiveTwo
	case "6_1":
		return schedule.SixOne
	case "6_2":
		return schedule.SixTwo
	default:
		return ""
	}
}

func buildFilteredResponse(schedules []Schedule, queueNames []string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(schedules))
	for i, schedule := range schedules {
		resultMap := map[string]interface{}{
			"id":            schedule.ID,
			"news_id":       schedule.NewsID,
			"title":         schedule.Title,
			"date":          schedule.Date,
			"schedule_date": schedule.ScheduleDate,
		}

		for _, queueName := range queueNames {
			resultMap[queueName] = getQueueValue(&schedule, queueName)
		}

		result[i] = resultMap
	}
	return result
}

func handleScheduleRequest(c *fiber.Ctx, db *gorm.DB, filter ScheduleFilter) error {
	var schedules []Schedule
	query := db.Table("schedules")
	switch filter.Option {
	case "all":
		if err := query.Find(&schedules).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve records"})
		}
	case "latest_n":
		if filter.Limit <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit value, it must be greater than zero"})
		}
		if err := query.Order("date desc").Limit(filter.Limit).Find(&schedules).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve latest records"})
		}
	case "by_date":
		if filter.Date == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Date must be specified in YYYY-MM-DD format"})
		}
		date, err := time.Parse("2006-01-02", filter.Date)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
		}
		if err := query.Where("DATE(date) = ?", date.Format("2006-01-02")).Find(&schedules).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve records by date"})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid option parameter value"})
	}

	for i := range schedules {
		schedules[i].ScheduleDate = utils.ExtractScheduleDateFromTitle(schedules[i].Title)
	}

	validatedQueues, err := parseAndValidateQueues(filter.Queue)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if len(validatedQueues) == 0 {
		return c.JSON(schedules)
	}

	return c.JSON(buildFilteredResponse(schedules, validatedQueues))
}

func GetSchedule(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filter := ScheduleFilter{
			Option: c.Query("option"),
			Date:   c.Query("date"),
			Queue:  c.Query("queue"),
		}

		if limitStr := c.Query("limit"); limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit parameter, expected integer value"})
			}
			filter.Limit = limit
		}

		return handleScheduleRequest(c, db, filter)
	}
}
