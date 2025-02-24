package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Schedule struct {
	ID       int64     `json:"id"`
	NewsID   int       `json:"news_id"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	OneOne   string    `gorm:"column:1_1" json:"1_1"`
	OneTwo   string    `gorm:"column:1_2" json:"1_2"`
	TwoOne   string    `gorm:"column:2_1" json:"2_1"`
	TwoTwo   string    `gorm:"column:2_2" json:"2_2"`
	ThreeOne string    `gorm:"column:3_1" json:"3_1"`
	ThreeTwo string    `gorm:"column:3_2" json:"3_2"`
	FourOne  string    `gorm:"column:4_1" json:"4_1"`
	FourTwo  string    `gorm:"column:4_2" json:"4_2"`
	FiveOne  string    `gorm:"column:5_1" json:"5_1"`
	FiveTwo  string    `gorm:"column:5_2" json:"5_2"`
	SixOne   string    `gorm:"column:6_1" json:"6_1"`
	SixTwo   string    `gorm:"column:6_2" json:"6_2"`
}

type ScheduleFilter struct {
	Option string `json:"option"`
	Date   string `json:"date"`
	Limit  int    `json:"limit"`
}

func PostSchedule(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var filter ScheduleFilter
		if err := c.BodyParser(&filter); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		}

		var schedules []Schedule
		query := db.Table("schedules")

		switch filter.Option {
		case "all":
			if err := query.Find(&schedules).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to retrieve records",
				})
			}
		case "latest_n":
			if filter.Limit <= 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid limit value, it must be greater than zero",
				})
			}
			if err := query.Order("date desc").Limit(filter.Limit).Find(&schedules).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to retrieve latest records",
				})
			}
		case "by_date":
			if filter.Date == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Date must be specified in YYYY-MM-DD format",
				})
			}
			date, err := time.Parse("2006-01-02", filter.Date)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid date format, expected YYYY-MM-DD",
				})
			}
			if err := query.Where("DATE(date) = ?", date.Format("2006-01-02")).Find(&schedules).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to retrieve records by date",
				})
			}
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid option parameter value",
			})
		}

		return c.JSON(schedules)
	}
}
