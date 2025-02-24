package parser

import (
	"strings"
	"testing"
	"time"
)

func TestContainsScheduleKeywords(t *testing.T) {
	tests := []struct {
		title    string
		expected bool
	}{
		{"оновлені графіки", true},
		{"random title", false},
	}
	for _, tt := range tests {
		got := containsScheduleKeywords(tt.title)
		if got != tt.expected {
			t.Errorf("For title %q, expected %v but got %v", tt.title, tt.expected, got)
		}
	}
}

func TestParseScheduleDataCanceled(t *testing.T) {
	htmlBody := "скасовано"
	s := parseScheduleData(htmlBody)
	if s.Col1_1 != "скасовано" {
		t.Errorf("Expected Col1_1 to be 'скасовано', got %q", s.Col1_1)
	}
}

func TestParseScheduleDataTable(t *testing.T) {
	htmlBody := `
		<table>
			<tr><td>1.І</td><td>10:00</td></tr>
			<tr><td>2.ІІ</td><td>11:00</td></tr>
		</table>
	`
	s := parseScheduleData(htmlBody)
	if strings.TrimSpace(s.Col1_1) != "10:00" {
		t.Errorf("Expected Col1_1 to be '10:00', got %q", s.Col1_1)
	}
	if strings.TrimSpace(s.Col2_2) != "11:00" {
		t.Errorf("Expected Col2_2 to be '11:00', got %q", s.Col2_2)
	}
}

func TestDateParsing(t *testing.T) {
	dateStr := "31.12.2022 23:59"
	parsed, err := time.Parse("02.01.2006 15:04", dateStr)
	if err != nil {
		t.Errorf("Expected valid date parsing but got error: %v", err)
	}
	expected := time.Date(2022, 12, 31, 23, 59, 0, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, parsed)
	}
}
