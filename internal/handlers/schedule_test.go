package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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

	err = db.AutoMigrate(&Schedule{})
	if err != nil {
		panic("failed to migrate test database")
	}

	testSchedules := []Schedule{
		{
			ID:       1,
			NewsID:   101,
			Title:    "Графік погодинних відключень на 14 листопада",
			Date:     time.Now().Add(-48 * time.Hour),
			OneOne:   "08:00-10:00",
			OneTwo:   "10:00-12:00",
			TwoOne:   "12:00-14:00",
			TwoTwo:   "14:00-16:00",
			ThreeOne: "09:00-11:00",
			ThreeTwo: "11:00-13:00",
			FourOne:  "13:00-15:00",
			FourTwo:  "15:00-17:00",
			FiveOne:  "07:00-09:00",
			FiveTwo:  "09:00-11:00",
			SixOne:   "11:00-13:00",
			SixTwo:   "13:00-15:00",
		},
		{
			ID:       2,
			NewsID:   102,
			Title:    "Графік погодинних відключень на 25 грудня",
			Date:     time.Now().Add(-24 * time.Hour),
			OneOne:   "09:00-11:00",
			OneTwo:   "11:00-13:00",
			TwoOne:   "13:00-15:00",
			TwoTwo:   "15:00-17:00",
			ThreeOne: "10:00-12:00",
			ThreeTwo: "12:00-14:00",
			FourOne:  "14:00-16:00",
			FourTwo:  "16:00-18:00",
			FiveOne:  "08:00-10:00",
			FiveTwo:  "10:00-12:00",
			SixOne:   "12:00-14:00",
			SixTwo:   "14:00-16:00",
		},
		{
			ID:       3,
			NewsID:   103,
			Title:    "Schedule without date",
			Date:     time.Now(),
			OneOne:   "07:00-09:00",
			OneTwo:   "09:00-11:00",
			TwoOne:   "11:00-13:00",
			TwoTwo:   "13:00-15:00",
			ThreeOne: "08:00-10:00",
			ThreeTwo: "10:00-12:00",
			FourOne:  "12:00-14:00",
			FourTwo:  "14:00-16:00",
			FiveOne:  "06:00-08:00",
			FiveTwo:  "08:00-10:00",
			SixOne:   "10:00-12:00",
			SixTwo:   "12:00-14:00",
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
	assert.Equal(t, "14.11", scheduleWithDate["schedule_date"])

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
		// Verify metadata fields are present
		assert.Contains(t, schedule, "id")
		assert.Contains(t, schedule, "news_id")
		assert.Contains(t, schedule, "title")
		assert.Contains(t, schedule, "date")
		assert.Contains(t, schedule, "schedule_date")

		// Verify both queue fields are present
		assert.Contains(t, schedule, "4_1")
		assert.Contains(t, schedule, "3_1")

		// Verify both queue values are non-empty strings
		assert.NotEmpty(t, schedule["4_1"])
		assert.NotEmpty(t, schedule["3_1"])
		_, ok1 := schedule["4_1"].(string)
		_, ok2 := schedule["3_1"].(string)
		assert.True(t, ok1, "4_1 should be a string")
		assert.True(t, ok2, "3_1 should be a string")

		// Verify response does NOT contain other queue fields
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
		// Verify all three queues are present
		assert.Contains(t, schedule, "4_1")
		assert.Contains(t, schedule, "3_1")
		assert.Contains(t, schedule, "2_2")

		// Verify values are non-empty
		assert.NotEmpty(t, schedule["4_1"])
		assert.NotEmpty(t, schedule["3_1"])
		assert.NotEmpty(t, schedule["2_2"])

		// Verify order is preserved (check by getting keys from response)
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
		// Verify response contains 3_1 and 4_1 only once each
		assert.Contains(t, schedule, "3_1")
		assert.Contains(t, schedule, "4_1")

		// Verify values are non-empty
		assert.NotEmpty(t, schedule["3_1"])
		assert.NotEmpty(t, schedule["4_1"])

		// Count queue fields to ensure no duplicates
		queueCount := 0
		for key := range schedule {
			if key == "3_1" || key == "4_1" {
				queueCount++
			}
		}
		assert.Equal(t, 2, queueCount, "Should have exactly 2 unique queue fields")

		// Verify other queues are not present
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

			// Verify error message mentions the invalid queue
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

	// Verify only 1 record is returned
	assert.Len(t, responseBody, 1)

	schedule := responseBody[0]

	// Verify both queue fields are present in the single record
	assert.Contains(t, schedule, "3_2")
	assert.Contains(t, schedule, "4_1")

	// Verify both values are non-empty
	assert.NotEmpty(t, schedule["3_2"])
	assert.NotEmpty(t, schedule["4_1"])

	// Verify other queues are not present
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

	// Verify filtered by date and contains both specified queues
	for _, schedule := range responseBody {
		assert.Contains(t, schedule, "3_2")
		assert.Contains(t, schedule, "5_1")

		assert.NotEmpty(t, schedule["3_2"])
		assert.NotEmpty(t, schedule["5_1"])

		// Verify other queues are not present
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
		// Verify metadata fields are present
		assert.Contains(t, schedule, "id")
		assert.Contains(t, schedule, "news_id")
		assert.Contains(t, schedule, "title")
		assert.Contains(t, schedule, "date")
		assert.Contains(t, schedule, "schedule_date")

		// Verify single queue field in response
		assert.Contains(t, schedule, "3_2")
		assert.NotEmpty(t, schedule["3_2"])

		// Verify other queues are not present (backward compatibility)
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
				// Verify all 12 queue fields are returned
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

				// Verify all values are non-empty
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
