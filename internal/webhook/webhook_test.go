package webhook

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cherkasyoblenergo-api/internal/models"
)

func TestSendWebhook_PayloadIsArray(t *testing.T) {
	var receivedPayload []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		receivedPayload, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	testSchedules := []models.Schedule{
		{
			ID:           1,
			NewsID:       100,
			Title:        "Test Schedule 1",
			Date:         time.Date(2025, 12, 15, 10, 0, 0, 0, time.UTC),
			ScheduleDate: "15.12.2025",
			Col1_1:       "08:00-12:00",
			Col1_2:       "16:00-20:00",
		},
		{
			ID:           2,
			NewsID:       101,
			Title:        "Test Schedule 2",
			Date:         time.Date(2025, 12, 16, 10, 0, 0, 0, time.UTC),
			ScheduleDate: "16.12.2025",
			Col1_1:       "09:00-13:00",
			Col1_2:       "17:00-21:00",
		},
	}

	apiKey := models.APIKey{
		Key:        "test-api-key",
		WebhookURL: server.URL,
	}

	err := SendWebhook(apiKey, testSchedules)
	if err != nil {
		t.Fatalf("SendWebhook failed: %v", err)
	}

	var receivedPayloadStruct struct {
		Schedules []models.Schedule `json:"schedules"`
	}
	err = json.Unmarshal(receivedPayload, &receivedPayloadStruct)
	if err != nil {
		t.Fatalf("Failed to unmarshal payload as object with schedules key: %v", err)
	}

	receivedSchedules := receivedPayloadStruct.Schedules

	if len(receivedSchedules) != len(testSchedules) {
		t.Errorf("Expected %d schedules, got %d", len(testSchedules), len(receivedSchedules))
	}

	for i, expected := range testSchedules {
		if i >= len(receivedSchedules) {
			break
		}
		received := receivedSchedules[i]

		if received.ID != expected.ID {
			t.Errorf("Schedule %d: expected ID %d, got %d", i, expected.ID, received.ID)
		}
		if received.NewsID != expected.NewsID {
			t.Errorf("Schedule %d: expected NewsID %d, got %d", i, expected.NewsID, received.NewsID)
		}
		if received.Title != expected.Title {
			t.Errorf("Schedule %d: expected Title %s, got %s", i, expected.Title, received.Title)
		}
		if received.ScheduleDate != expected.ScheduleDate {
			t.Errorf("Schedule %d: expected ScheduleDate %s, got %s", i, expected.ScheduleDate, received.ScheduleDate)
		}
	}

	if len(receivedPayload) == 0 || receivedPayload[0] != '{' {
		t.Error("Payload should start with '{' to be a JSON object")
	}
}
