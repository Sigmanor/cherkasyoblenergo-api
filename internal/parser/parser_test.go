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

// TestParseScheduleDataTable checks the old format (for backward compatibility)
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

func TestParseScheduleDataNewFormat(t *testing.T) {
	htmlBody := `
		<p>1.1 11:00 - 13:00, 19:00 - 21:00</p>
		<p>2.2 13:00 - 15:00, 21:00 - 23:00</p>
		<p>4.1 07:00 - 09:00, 15:00 - 17:00</p>
	`
	s := parseScheduleData(htmlBody)
	if strings.TrimSpace(s.Col1_1) != "11:00 - 13:00, 19:00 - 21:00" {
		t.Errorf("Expected Col1_1 to be '11:00 - 13:00, 19:00 - 21:00', got %q", s.Col1_1)
	}
	if strings.TrimSpace(s.Col2_2) != "13:00 - 15:00, 21:00 - 23:00" {
		t.Errorf("Expected Col2_2 to be '13:00 - 15:00, 21:00 - 23:00', got %q", s.Col2_2)
	}
	if strings.TrimSpace(s.Col4_1) != "07:00 - 09:00, 15:00 - 17:00" {
		t.Errorf("Expected Col4_1 to be '07:00 - 09:00, 15:00 - 17:00', got %q", s.Col4_1)
	}
}

func TestParseScheduleDataNewFormatSingleTimeRange(t *testing.T) {
	htmlBody := `
		<p>2.1 15:30 - 17:30</p>
	`
	s := parseScheduleData(htmlBody)
	if strings.TrimSpace(s.Col2_1) != "15:30 - 17:30" {
		t.Errorf("Expected Col2_1 to be '15:30 - 17:30', got %q", s.Col2_1)
	}
}

func TestParseScheduleDataMixedContent(t *testing.T) {
	htmlBody := `
		<p>Some text about Telegram</p>
		<p><a href="https://t.me/channel">Telegram Channel Link</a></p>
		<p>1.1 11:00 - 13:00, 19:00 - 21:00</p>
		<p>Some other text</p>
		<p>2.2 13:00 - 15:00, 21:00 - 23:00</p>
		<p>More text content</p>
		<p>4.1 07:00 - 09:00, 15:00 - 17:00</p>
	`
	s := parseScheduleData(htmlBody)
	if strings.TrimSpace(s.Col1_1) != "11:00 - 13:00, 19:00 - 21:00" {
		t.Errorf("Expected Col1_1 to be '11:00 - 13:00, 19:00 - 21:00', got %q", s.Col1_1)
	}
	if strings.TrimSpace(s.Col2_2) != "13:00 - 15:00, 21:00 - 23:00" {
		t.Errorf("Expected Col2_2 to be '13:00 - 15:00, 21:00 - 23:00', got %q", s.Col2_2)
	}
	if strings.TrimSpace(s.Col4_1) != "07:00 - 09:00, 15:00 - 17:00" {
		t.Errorf("Expected Col4_1 to be '07:00 - 09:00, 15:00 - 17:00', got %q", s.Col4_1)
	}
}

func TestParseScheduleFromParagraphs(t *testing.T) {
	htmlBody := `
		<p>1.1 10:00-12:00, 14:00-16:00</p>
		<p>2.2 11:00-13:00</p>
		<p>Some other text</p>
		<p>3.1 09:00-11:00, 15:00-17:00</p>
	`
	s, found := parseScheduleFromParagraphs(htmlBody)
	if !found {
		t.Errorf("Expected to find schedule data but got false")
	}
	if strings.TrimSpace(s.Col1_1) != "10:00-12:00, 14:00-16:00" {
		t.Errorf("Expected Col1_1 to be '10:00-12:00, 14:00-16:00', got %q", s.Col1_1)
	}
	if strings.TrimSpace(s.Col2_2) != "11:00-13:00" {
		t.Errorf("Expected Col2_2 to be '11:00-13:00', got %q", s.Col2_2)
	}
	if strings.TrimSpace(s.Col3_1) != "09:00-11:00, 15:00-17:00" {
		t.Errorf("Expected Col3_1 to be '09:00-11:00, 15:00-17:00', got %q", s.Col3_1)
	}
}

func TestParseScheduleDataWithParagraphs(t *testing.T) {
	htmlBody := `
		<p>1.1 10:00-12:00</p>
		<p>2.2 11:00-13:00</p>
	`
	s := parseScheduleData(htmlBody)
	if strings.TrimSpace(s.Col1_1) != "10:00-12:00" {
		t.Errorf("Expected Col1_1 to be '10:00-12:00', got %q", s.Col1_1)
	}
	if strings.TrimSpace(s.Col2_2) != "11:00-13:00" {
		t.Errorf("Expected Col2_2 to be '11:00-13:00', got %q", s.Col2_2)
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

func TestContainsSchedulePatterns(t *testing.T) {
	tests := []struct {
		htmlBody string
		expected bool
	}{
		{`<p>1.1 11:00 - 13:00, 19:00 - 21:00</p>`, true},
		{`<p>6.2 09:00-11:00</p>`, true},
		{`<p>Some random text without schedule patterns</p>`, false},
		{`<p>3.1 15:30 - 17:30</p><p>Other text</p>`, true},
		{`<p>No schedule here</p><p>Just regular text</p>`, false},
	}
	
	for _, tt := range tests {
		got := containsSchedulePatterns(tt.htmlBody)
		if got != tt.expected {
			t.Errorf("For htmlBody %q, expected %v but got %v", tt.htmlBody, tt.expected, got)
		}
	}
}
