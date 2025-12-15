package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"cherkasyoblenergo-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	err = db.AutoMigrate(&models.Schedule{})
	if err != nil {
		panic("failed to migrate test database")
	}

	testSchedules := []models.Schedule{
		{
			ID:     1,
			NewsID: 101,
			Title:  "Графік погодинних відключень на 14 листопада",
			Date:   time.Now().Add(-48 * time.Hour),
			Col1_1: "08:00-10:00",
			Col1_2: "10:00-12:00",
			Col2_1: "12:00-14:00",
			Col2_2: "14:00-16:00",
			Col3_1: "09:00-11:00",
			Col3_2: "11:00-13:00",
			Col4_1: "13:00-15:00",
			Col4_2: "15:00-17:00",
			Col5_1: "07:00-09:00",
			Col5_2: "09:00-11:00",
			Col6_1: "11:00-13:00",
			Col6_2: "13:00-15:00",
		},
		{
			ID:     2,
			NewsID: 102,
			Title:  "Графік погодинних відключень на 25 грудня",
			Date:   time.Now().Add(-24 * time.Hour),
			Col1_1: "09:00-11:00",
			Col1_2: "11:00-13:00",
			Col2_1: "13:00-15:00",
			Col2_2: "15:00-17:00",
			Col3_1: "10:00-12:00",
			Col3_2: "12:00-14:00",
			Col4_1: "14:00-16:00",
			Col4_2: "16:00-18:00",
			Col5_1: "08:00-10:00",
			Col5_2: "10:00-12:00",
			Col6_1: "12:00-14:00",
			Col6_2: "14:00-16:00",
		},
		{
			ID:     3,
			NewsID: 103,
			Title:  "Schedule without date",
			Date:   time.Now(),
			Col1_1: "07:00-09:00",
			Col1_2: "09:00-11:00",
			Col2_1: "11:00-13:00",
			Col2_2: "13:00-15:00",
			Col3_1: "08:00-10:00",
			Col3_2: "10:00-12:00",
			Col4_1: "12:00-14:00",
			Col4_2: "14:00-16:00",
			Col5_1: "06:00-08:00",
			Col5_2: "08:00-10:00",
			Col6_1: "10:00-12:00",
			Col6_2: "12:00-14:00",
		},
	}

	for _, schedule := range testSchedules {
		db.Create(&schedule)
	}

	return db
}

func newGetScheduleRequest(params map[string]string) *http.Request {
	path := "/schedule"
	query := url.Values{}

	for key, value := range params {
		if value != "" {
			query.Set(key, value)
		}
	}

	if encoded := query.Encode(); encoded != "" {
		path = path + "?" + encoded
	}

	return httptest.NewRequest("GET", path, nil)
}

func TestGetSchedule_AllOption(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "all",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGetSchedule_QueueFilter(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "all",
		"queue":  "3_2",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody)

	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "id")
		assert.Contains(t, schedule, "news_id")
		assert.Contains(t, schedule, "title")
		assert.Contains(t, schedule, "date")
		assert.Contains(t, schedule, "schedule_date")

		assert.Contains(t, schedule, "3_2")
		assert.NotEmpty(t, schedule["3_2"])

		assert.NotContains(t, schedule, "3_1")
		assert.NotContains(t, schedule, "1_1")
		assert.NotContains(t, schedule, "1_2")
		assert.NotContains(t, schedule, "2_1")
		assert.NotContains(t, schedule, "2_2")
		assert.NotContains(t, schedule, "4_1")
		assert.NotContains(t, schedule, "4_2")
		assert.NotContains(t, schedule, "5_1")
		assert.NotContains(t, schedule, "5_2")
		assert.NotContains(t, schedule, "6_1")
		assert.NotContains(t, schedule, "6_2")
	}
}

func TestGetSchedule_QueueFilter_InvalidFormat(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	invalidQueues := []string{"3-2", "7_1", "1_3", "abc"}

	for _, invalidQueue := range invalidQueues {
		t.Run("Queue_"+invalidQueue, func(t *testing.T) {
			req := newGetScheduleRequest(map[string]string{
				"option": "all",
				"queue":  invalidQueue,
			})

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, 400, resp.StatusCode)

			bodyBytes, _ := io.ReadAll(resp.Body)
			var errorResponse map[string]interface{}
			err = json.Unmarshal(bodyBytes, &errorResponse)
			assert.NoError(t, err)
			assert.Contains(t, errorResponse, "error")
		})
	}
}

func TestGetSchedule_QueueFilter_WithLatestN(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "latest_n",
		"limit":  "1",
		"queue":  "3_2",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)

	assert.Len(t, responseBody, 1)

	schedule := responseBody[0]
	assert.Contains(t, schedule, "3_2")
	assert.NotEmpty(t, schedule["3_2"])
	assert.NotContains(t, schedule, "3_1")
	assert.NotContains(t, schedule, "4_1")
}

func TestGetSchedule_QueueFilter_WithByDate(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "by_date",
		"date":   time.Now().Format("2006-01-02"),
		"queue":  "3_2",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)

	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "3_2")
		assert.NotEmpty(t, schedule["3_2"])
		assert.NotContains(t, schedule, "3_1")
		assert.NotContains(t, schedule, "5_1")
	}
}

func TestGetSchedule_ScheduleDateField(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "all",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody)

	var scheduleWithDate map[string]interface{}
	var scheduleWithoutDate map[string]interface{}

	for _, schedule := range responseBody {
		title, ok := schedule["title"].(string)
		if ok {
			switch title {
			case "Графік погодинних відключень на 14 листопада":
				scheduleWithDate = schedule
			case "Schedule without date":
				scheduleWithoutDate = schedule
			}
		}
	}

	assert.NotNil(t, scheduleWithDate)
	assert.Contains(t, scheduleWithDate, "schedule_date")
	scheduleDateStr, ok := scheduleWithDate["schedule_date"].(string)
	assert.True(t, ok)
	assert.Regexp(t, `^\d{4}-11-14$`, scheduleDateStr) 
	assert.NotNil(t, scheduleWithoutDate)
	assert.Contains(t, scheduleWithoutDate, "schedule_date")
	assert.Equal(t, "", scheduleWithoutDate["schedule_date"])
}

func TestGetSchedule_ScheduleDateField_AllOptions(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	testCases := []struct {
		name   string
		params map[string]string
	}{
		{
			name: "All_NoQueue",
			params: map[string]string{
				"option": "all",
			},
		},
		{
			name: "All_WithQueue",
			params: map[string]string{
				"option": "all",
				"queue":  "3_2",
			},
		},
		{
			name: "LatestN_NoQueue",
			params: map[string]string{
				"option": "latest_n",
				"limit":  "1",
			},
		},
		{
			name: "LatestN_WithQueue",
			params: map[string]string{
				"option": "latest_n",
				"limit":  "1",
				"queue":  "3_2",
			},
		},
		{
			name: "ByDate_NoQueue",
			params: map[string]string{
				"option": "by_date",
				"date":   time.Now().Format("2006-01-02"),
			},
		},
		{
			name: "ByDate_WithQueue",
			params: map[string]string{
				"option": "by_date",
				"date":   time.Now().Format("2006-01-02"),
				"queue":  "3_2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := newGetScheduleRequest(tc.params)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)

			bodyBytes, _ := io.ReadAll(resp.Body)
			var responseBody []map[string]interface{}
			err = json.Unmarshal(bodyBytes, &responseBody)
			assert.NoError(t, err)
			assert.NotEmpty(t, responseBody)

			for _, schedule := range responseBody {
				assert.Contains(t, schedule, "schedule_date")
				_, ok := schedule["schedule_date"].(string)
				assert.True(t, ok, "schedule_date should be a string")
			}
		})
	}
}

func TestGetSchedule_MultipleQueues_Success(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "all",
		"queue":  "4_1, 3_1",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody)

	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "id")
		assert.Contains(t, schedule, "news_id")
		assert.Contains(t, schedule, "title")
		assert.Contains(t, schedule, "date")
		assert.Contains(t, schedule, "schedule_date")

		assert.Contains(t, schedule, "4_1")
		assert.Contains(t, schedule, "3_1")

		assert.NotEmpty(t, schedule["4_1"])
		assert.NotEmpty(t, schedule["3_1"])
		_, ok1 := schedule["4_1"].(string)
		_, ok2 := schedule["3_1"].(string)
		assert.True(t, ok1, "4_1 should be a string")
		assert.True(t, ok2, "3_1 should be a string")

		assert.NotContains(t, schedule, "1_1")
		assert.NotContains(t, schedule, "1_2")
		assert.NotContains(t, schedule, "2_1")
		assert.NotContains(t, schedule, "2_2")
		assert.NotContains(t, schedule, "3_2")
		assert.NotContains(t, schedule, "4_2")
		assert.NotContains(t, schedule, "5_1")
		assert.NotContains(t, schedule, "5_2")
		assert.NotContains(t, schedule, "6_1")
		assert.NotContains(t, schedule, "6_2")
	}
}

func TestGetSchedule_MultipleQueues_NoSpaces(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "all",
		"queue":  "4_1,3_1,2_2",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody)

	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "4_1")
		assert.Contains(t, schedule, "3_1")
		assert.Contains(t, schedule, "2_2")

		assert.NotEmpty(t, schedule["4_1"])
		assert.NotEmpty(t, schedule["3_1"])
		assert.NotEmpty(t, schedule["2_2"])

		keys := make([]string, 0)
		for key := range schedule {
			if key == "4_1" || key == "3_1" || key == "2_2" {
				keys = append(keys, key)
			}
		}
		assert.Len(t, keys, 3)
	}
}

func TestGetSchedule_MultipleQueues_WithDuplicates(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "all",
		"queue":  "3_1, 4_1, 3_1",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody)

	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "3_1")
		assert.Contains(t, schedule, "4_1")

		assert.NotEmpty(t, schedule["3_1"])
		assert.NotEmpty(t, schedule["4_1"])

		queueCount := 0
		for key := range schedule {
			if key == "3_1" || key == "4_1" {
				queueCount++
			}
		}
		assert.Equal(t, 2, queueCount, "Should have exactly 2 unique queue fields")

		assert.NotContains(t, schedule, "1_1")
		assert.NotContains(t, schedule, "2_2")
		assert.NotContains(t, schedule, "5_1")
	}
}

func TestGetSchedule_MultipleQueues_InvalidQueue(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	testCases := []struct {
		name  string
		queue string
	}{
		{
			name:  "InvalidQueue_7_1",
			queue: "3_1, 7_1",
		},
		{
			name:  "InvalidQueue_abc",
			queue: "3_1, abc",
		},
		{
			name:  "InvalidQueue_1_3",
			queue: "1_3, 2_1",
		},
		{
			name:  "InvalidQueue_hyphen",
			queue: "3-1, 4_1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := newGetScheduleRequest(map[string]string{
				"option": "all",
				"queue":  tc.queue,
			})

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, 400, resp.StatusCode)

			bodyBytes, _ := io.ReadAll(resp.Body)
			var errorResponse map[string]interface{}
			err = json.Unmarshal(bodyBytes, &errorResponse)
			assert.NoError(t, err)
			assert.Contains(t, errorResponse, "error")

			errorMsg, ok := errorResponse["error"].(string)
			assert.True(t, ok)
			assert.NotEmpty(t, errorMsg)
		})
	}
}

func TestGetSchedule_MultipleQueues_WithLatestN(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "latest_n",
		"limit":  "1",
		"queue":  "3_2, 4_1",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)

	assert.Len(t, responseBody, 1)

	schedule := responseBody[0]

	assert.Contains(t, schedule, "3_2")
	assert.Contains(t, schedule, "4_1")

	assert.NotEmpty(t, schedule["3_2"])
	assert.NotEmpty(t, schedule["4_1"])

	assert.NotContains(t, schedule, "3_1")
	assert.NotContains(t, schedule, "5_1")
	assert.NotContains(t, schedule, "1_1")
}

func TestGetSchedule_MultipleQueues_WithByDate(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "by_date",
		"date":   time.Now().Format("2006-01-02"),
		"queue":  "3_2, 5_1",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)

	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "3_2")
		assert.Contains(t, schedule, "5_1")

		assert.NotEmpty(t, schedule["3_2"])
		assert.NotEmpty(t, schedule["5_1"])

		assert.NotContains(t, schedule, "3_1")
		assert.NotContains(t, schedule, "4_1")
		assert.NotContains(t, schedule, "1_1")
		assert.NotContains(t, schedule, "2_2")
	}
}

func TestGetSchedule_SingleQueue_BackwardCompatibility(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "all",
		"queue":  "3_2",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody)

	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "id")
		assert.Contains(t, schedule, "news_id")
		assert.Contains(t, schedule, "title")
		assert.Contains(t, schedule, "date")
		assert.Contains(t, schedule, "schedule_date")

		assert.Contains(t, schedule, "3_2")
		assert.NotEmpty(t, schedule["3_2"])

		assert.NotContains(t, schedule, "3_1")
		assert.NotContains(t, schedule, "1_1")
		assert.NotContains(t, schedule, "1_2")
		assert.NotContains(t, schedule, "2_1")
		assert.NotContains(t, schedule, "2_2")
		assert.NotContains(t, schedule, "4_1")
		assert.NotContains(t, schedule, "4_2")
		assert.NotContains(t, schedule, "5_1")
		assert.NotContains(t, schedule, "5_2")
		assert.NotContains(t, schedule, "6_1")
		assert.NotContains(t, schedule, "6_2")
	}
}

func TestGetSchedule_EmptyQueue_ReturnsAllQueues(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	testCases := []struct {
		name  string
		queue string
	}{
		{
			name:  "EmptyString",
			queue: "",
		},
		{
			name:  "WhitespaceOnly",
			queue: "  ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := newGetScheduleRequest(map[string]string{
				"option": "all",
				"queue":  tc.queue,
			})

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)

			bodyBytes, _ := io.ReadAll(resp.Body)
			var responseBody []map[string]interface{}
			err = json.Unmarshal(bodyBytes, &responseBody)
			assert.NoError(t, err)
			assert.NotEmpty(t, responseBody)

			for _, schedule := range responseBody {
				assert.Contains(t, schedule, "1_1")
				assert.Contains(t, schedule, "1_2")
				assert.Contains(t, schedule, "2_1")
				assert.Contains(t, schedule, "2_2")
				assert.Contains(t, schedule, "3_1")
				assert.Contains(t, schedule, "3_2")
				assert.Contains(t, schedule, "4_1")
				assert.Contains(t, schedule, "4_2")
				assert.Contains(t, schedule, "5_1")
				assert.Contains(t, schedule, "5_2")
				assert.Contains(t, schedule, "6_1")
				assert.Contains(t, schedule, "6_2")

				assert.NotEmpty(t, schedule["1_1"])
				assert.NotEmpty(t, schedule["1_2"])
				assert.NotEmpty(t, schedule["2_1"])
				assert.NotEmpty(t, schedule["2_2"])
				assert.NotEmpty(t, schedule["3_1"])
				assert.NotEmpty(t, schedule["3_2"])
				assert.NotEmpty(t, schedule["4_1"])
				assert.NotEmpty(t, schedule["4_2"])
				assert.NotEmpty(t, schedule["5_1"])
				assert.NotEmpty(t, schedule["5_2"])
				assert.NotEmpty(t, schedule["6_1"])
				assert.NotEmpty(t, schedule["6_2"])
			}
		})
	}
}

func TestGetSchedule_LatestNWithQueues(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := httptest.NewRequest("GET", "/schedule?option=latest_n&limit=2&queue=3_2,4_1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody)
	assert.LessOrEqual(t, len(responseBody), 2)

	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "id")
		assert.Contains(t, schedule, "news_id")
		assert.Contains(t, schedule, "title")
		assert.Contains(t, schedule, "date")
		assert.Contains(t, schedule, "schedule_date")

		assert.Contains(t, schedule, "3_2")
		assert.Contains(t, schedule, "4_1")
		assert.NotEmpty(t, schedule["3_2"])
		assert.NotEmpty(t, schedule["4_1"])

		assert.NotContains(t, schedule, "1_1")
		assert.NotContains(t, schedule, "1_2")
		assert.NotContains(t, schedule, "2_1")
		assert.NotContains(t, schedule, "2_2")
		assert.NotContains(t, schedule, "3_1")
		assert.NotContains(t, schedule, "4_2")
		assert.NotContains(t, schedule, "5_1")
		assert.NotContains(t, schedule, "5_2")
		assert.NotContains(t, schedule, "6_1")
		assert.NotContains(t, schedule, "6_2")
	}
}

func TestGetSchedule_InvalidLimit(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := httptest.NewRequest("GET", "/schedule?option=latest_n&limit=abc", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var errorResponse map[string]interface{}
	err = json.Unmarshal(bodyBytes, &errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, errorResponse, "error")
}

func TestGetSchedule_ByScheduleDate_Success(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	now := time.Now()
	year := now.Year()
	if 11 > int(now.Month()) {
		year--
	}
	expectedDate := fmt.Sprintf("%d-11-14", year)

	req := newGetScheduleRequest(map[string]string{
		"option": "by_schedule_date",
		"date":   expectedDate,
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)

	for _, schedule := range responseBody {
		scheduleDate, ok := schedule["schedule_date"].(string)
		assert.True(t, ok)
		assert.Equal(t, expectedDate, scheduleDate)
	}
}

func TestGetSchedule_ByScheduleDate_WithLimit(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	now := time.Now()
	year := now.Year()
	if 11 > int(now.Month()) {
		year--
	}
	expectedDate := fmt.Sprintf("%d-11-14", year)

	req := newGetScheduleRequest(map[string]string{
		"option": "by_schedule_date",
		"date":   expectedDate,
		"limit":  "1",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(responseBody), 1)
}

func TestGetSchedule_ByScheduleDate_WithQueue(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	now := time.Now()
	year := now.Year()
	if 11 > int(now.Month()) {
		year--
	}
	expectedDate := fmt.Sprintf("%d-11-14", year)

	req := newGetScheduleRequest(map[string]string{
		"option": "by_schedule_date",
		"date":   expectedDate,
		"queue":  "4_1",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)

	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "4_1")
		assert.NotContains(t, schedule, "3_1")
		assert.NotContains(t, schedule, "5_1")
	}
}

func TestGetSchedule_ByScheduleDate_NoDate(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "by_schedule_date",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var errorResponse map[string]interface{}
	err = json.Unmarshal(bodyBytes, &errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, errorResponse, "error")
}

func TestGetSchedule_ByScheduleDate_InvalidDateFormat(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	invalidDates := []string{"14.11", "2024/11/14", "14-11-2024", "abc"}

	for _, invalidDate := range invalidDates {
		t.Run("Date_"+invalidDate, func(t *testing.T) {
			req := newGetScheduleRequest(map[string]string{
				"option": "by_schedule_date",
				"date":   invalidDate,
			})

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

			bodyBytes, _ := io.ReadAll(resp.Body)
			var errorResponse map[string]interface{}
			err = json.Unmarshal(bodyBytes, &errorResponse)
			assert.NoError(t, err)
			assert.Contains(t, errorResponse, "error")
		})
	}
}

func TestGetSchedule_ByScheduleDate_Today(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "by_schedule_date",
		"date":   "today",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
}

func TestGetSchedule_ByScheduleDate_Tomorrow(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "by_schedule_date",
		"date":   "tomorrow",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
}

func TestGetSchedule_ByDate_Today(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "by_date",
		"date":   "today",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
}

func TestGetSchedule_ByDate_Tomorrow(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	req := newGetScheduleRequest(map[string]string{
		"option": "by_date",
		"date":   "tomorrow",
	})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var responseBody []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)
}

func TestGetSchedule_DateValues_CaseInsensitive(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	app.Get("/schedule", GetSchedule(db))

	testCases := []string{"TODAY", "Today", "TodAY", "TOMORROW", "Tomorrow", "toMORROW"}

	for _, dateValue := range testCases {
		t.Run("Date_"+dateValue, func(t *testing.T) {
			req := newGetScheduleRequest(map[string]string{
				"option": "by_schedule_date",
				"date":   dateValue,
			})

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})
	}
}
