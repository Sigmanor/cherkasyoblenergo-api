package handlers

import (
	"cherkasyoblenergo-api/internal/models"
	"cherkasyoblenergo-api/internal/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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

func getQueueValue(schedule *models.Schedule, queueName string) string {
	switch queueName {
	case "1_1":
		return schedule.Col1_1
	case "1_2":
		return schedule.Col1_2
	case "2_1":
		return schedule.Col2_1
	case "2_2":
		return schedule.Col2_2
	case "3_1":
		return schedule.Col3_1
	case "3_2":
		return schedule.Col3_2
	case "4_1":
		return schedule.Col4_1
	case "4_2":
		return schedule.Col4_2
	case "5_1":
		return schedule.Col5_1
	case "5_2":
		return schedule.Col5_2
	case "6_1":
		return schedule.Col6_1
	case "6_2":
		return schedule.Col6_2
	default:
		return ""
	}
}

func buildFilteredResponse(schedules []models.Schedule, queueNames []string) []map[string]interface{} {
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

func resolveDateValue(dateStr string) string {
	now := time.Now()
	switch strings.ToLower(dateStr) {
	case "today":
		return now.Format("2006-01-02")
	case "tomorrow":
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	default:
		return dateStr
	}
}

func handleScheduleRequest(c *fiber.Ctx, db *gorm.DB, filter ScheduleFilter) error {
	filter.Date = resolveDateValue(filter.Date)
	var schedules []models.Schedule
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
		if err := query.Order("id desc").Limit(filter.Limit).Find(&schedules).Error; err != nil {
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
	case "by_schedule_date":
		if filter.Date == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Date must be specified in YYYY-MM-DD format"})
		}
		if _, err := time.Parse("2006-01-02", filter.Date); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
		}
		if filter.Limit < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit value, it must be greater than or equal to zero"})
		}

		scheduleQuery := db.Table("schedules").Where("schedule_date = ?", filter.Date).Order("date desc")
		if filter.Limit > 0 {
			scheduleQuery = scheduleQuery.Limit(filter.Limit)
		}
		if err := scheduleQuery.Find(&schedules).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve records"})
		}

		neededMore := filter.Limit == 0 || len(schedules) < filter.Limit
		if neededMore {
			var legacySchedules []models.Schedule
			legacyQuery := db.Table("schedules").Where("schedule_date = '' OR schedule_date IS NULL").Order("date desc")
			if err := legacyQuery.Find(&legacySchedules).Error; err == nil {
				for i := range legacySchedules {
					legacySchedules[i].ScheduleDate = utils.ExtractScheduleDateFromTitle(legacySchedules[i].Title)
					if legacySchedules[i].ScheduleDate == filter.Date {
						schedules = append(schedules, legacySchedules[i])
						if filter.Limit > 0 && len(schedules) >= filter.Limit {
							break
						}
					}
				}
			}
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid option parameter value"})
	}

	for i := range schedules {
		if schedules[i].ScheduleDate == "" {
			schedules[i].ScheduleDate = utils.ExtractScheduleDateFromTitle(schedules[i].Title)
		}
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
