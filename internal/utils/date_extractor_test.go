package utils

import (
	"testing"
)

func TestExtractScheduleDateFromTitle_ValidTitles(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"Графік погодинних відключень на 14 листопада", "14.11"},
		{"Графік на 5 травня", "05.05"},
		{"Графік на 28 лютого", "28.02"},
		{"Графік відключень на 1 січня", "01.01"},
		{"Графік на 15 грудня", "15.12"},
		{"Оновлений графік на 3 березня", "03.03"},
		{"Графік на 10 квітня", "10.04"},
		{"Графік на 22 червня", "22.06"},
		{"Графік на 7 липня", "07.07"},
		{"Графік на 31 серпня", "31.08"},
		{"Графік на 9 вересня", "09.09"},
		{"Графік на 18 жовтня", "18.10"},
		
		{"Графік на 23 жовтня 2025 року", "23.10"},
		{"Графік 22 жовтня 2025 року", "22.10"},
		{"Оновлений графік 15 листопада 2024 року", "15.11"},
		
		{"Графік з 15:30, 22 жовтня 2025 року", "22.10"},
		{"Оновлення з 09:00, 5 травня 2025 року", "05.05"},
		{"Графік з 18:45, 31 грудня 2024 року", "31.12"},
	}

	for _, tt := range tests {
		got := ExtractScheduleDateFromTitle(tt.title)
		if got != tt.expected {
			t.Errorf("ExtractScheduleDateFromTitle(%q) = %q, want %q", tt.title, got, tt.expected)
		}
	}
}

func TestExtractScheduleDateFromTitle_Typos(t *testing.T) {
	tests := []struct {
		title    string
		expected string
		desc     string
	}{
		{"Графік на 14 листопаде", "14.11", "typo in month ending"},
		{"Графік на 5 листопда", "05.11", "missing letter in month"},
		{"Графік на 28 листопадаа", "28.11", "extra letter in month"},
		{"Графік на 15 груня", "15.12", "missing letter д"},
		{"Графік на 7 травна", "07.05", "typo in month ending"},
		{"Графік на 31 серпна", "31.08", "typo in august ending"},
	}

	for _, tt := range tests {
		got := ExtractScheduleDateFromTitle(tt.title)
		if got != tt.expected {
			t.Errorf("%s: ExtractScheduleDateFromTitle(%q) = %q, want %q", tt.desc, tt.title, got, tt.expected)
		}
	}
}

func TestExtractScheduleDateFromTitle_EdgeCases(t *testing.T) {
	tests := []struct {
		title    string
		expected string
		desc     string
	}{
		{"Графік без дати", "", "title without date"},
		{"Графік на 0 листопада", "", "invalid day 0"},
		{"Графік на 32 листопада", "", "invalid day 32"},
		{"Графік на 15 unknownmonth", "", "unrecognizable month name"},
		{"", "", "empty string"},
		{"На 14 number", "", "month is not Ukrainian"},
		{"Графік відключень", "", "no date pattern"},
		{"На листопада 14", "", "wrong order - month before day"},
		{"Графік - 5 травня", "05.05", "dash treated as punctuation, not negative"},
		{"Графік на 40 грудня", "", "day out of range"},
		{"Some random text", "", "completely different format"},
		{"Графік на 15", "", "missing month"},
		{"Графік листопада", "", "missing day"},
	}

	for _, tt := range tests {
		got := ExtractScheduleDateFromTitle(tt.title)
		if got != tt.expected {
			t.Errorf("%s: ExtractScheduleDateFromTitle(%q) = %q, want %q", tt.desc, tt.title, got, tt.expected)
		}
	}
}

func TestParseUkrainianMonth(t *testing.T) {
	tests := []struct {
		monthName string
		expected  int
		shouldOk  bool
		desc      string
	}{
		{"січня", 1, true, "January"},
		{"лютого", 2, true, "February"},
		{"березня", 3, true, "March"},
		{"квітня", 4, true, "April"},
		{"травня", 5, true, "May"},
		{"червня", 6, true, "June"},
		{"липня", 7, true, "July"},
		{"серпня", 8, true, "August"},
		{"вересня", 9, true, "September"},
		{"жовтня", 10, true, "October"},
		{"листопада", 11, true, "November"},
		{"грудня", 12, true, "December"},

		{"СІЧНЯ", 1, true, "uppercase January"},
		{"ЛюТоГо", 2, true, "mixed case February"},

		{"листопаде", 11, true, "typo in November ending"},
		{"листопда", 11, true, "missing letter in November"},
		{"листопадаа", 11, true, "extra letter in November"},
		{"груня", 12, true, "missing letter in December"},
		{"травна", 5, true, "typo in May ending"},
		{"серпна", 8, true, "typo in August ending"},

		{"invalid", 0, false, "completely invalid month"},
		{"january", 0, false, "English month name"},
		{"месяц", 0, false, "Russian month name"},
		{"", 0, false, "empty string"},
		{"а", 0, false, "single character"},
		{"xyz", 0, false, "random characters"},
		{"листо", 0, false, "too short - below similarity threshold"},
	}

	for _, tt := range tests {
		got, ok := parseUkrainianMonth(tt.monthName)
		if ok != tt.shouldOk {
			t.Errorf("%s: parseUkrainianMonth(%q) ok = %v, want %v", tt.desc, tt.monthName, ok, tt.shouldOk)
		}
		if ok && got != tt.expected {
			t.Errorf("%s: parseUkrainianMonth(%q) = %d, want %d", tt.desc, tt.monthName, got, tt.expected)
		}
	}
}

func TestExtractDayAndMonth(t *testing.T) {
	tests := []struct {
		title        string
		expectedDay  int
		expectedMonth int
		shouldOk     bool
		desc         string
	}{
		{"Графік на 14 листопада", 14, 11, true, "standard format"},
		{"на 5 травня", 5, 5, true, "minimal format"},
		{"Оновлений графік на 28 лютого 2024", 28, 2, true, "with year"},
		{"На 1 січня", 1, 1, true, "first day of month"},
		{"на 31 грудня", 31, 12, true, "last day of month"},
	
		{"на 1 березня", 1, 3, true, "minimum valid day"},
		{"на 31 серпня", 31, 8, true, "maximum valid day"},
	
		{"на 23 жовтня 2025 року", 23, 10, true, "format with year after month"},
		{"22 жовтня 2025 року", 22, 10, true, "format without 'на' prefix with year"},
		{"15 листопада 2024 року", 15, 11, true, "another format with year"},
	
		{"з 15:30, 22 жовтня 2025 року", 22, 10, true, "format with time and year"},
		{"з 09:00, 5 травня 2025 року", 5, 5, true, "format with time in morning"},
		{"з 18:45, 31 грудня 2024 року", 31, 12, true, "format with time in evening"},

		{"на 0 травня", 0, 0, false, "day is 0"},
		{"на 32 липня", 0, 0, false, "day is 32"},
		{"на - 5 червня", 5, 6, true, "dash before number (treated as punctuation)"},
		{"на 100 вересня", 0, 0, false, "day is 100"},

		{"на 15", 0, 0, false, "missing month"},
		{"на листопада", 0, 0, false, "missing day"},
		{"графік відключень", 0, 0, false, "no date pattern"},
		{"", 0, 0, false, "empty string"},

		{"15 листопада", 15, 11, true, "universal format without 'на' keyword"},
		{"5 травня", 5, 5, true, "universal format single digit day"},
		{"28 лютого", 28, 2, true, "universal format february"},
		{"Графік 22 червня", 22, 6, true, "universal format with prefix text"},
		{"на листопада 15", 0, 0, false, "month before day"},
		{"нa 14 листопада", 0, 0, false, "Latin 'a' in 'на'"},

		{"на 14 invalid", 0, 0, false, "invalid month name"},
		{"на 14 january", 0, 0, false, "English month name"},
		{"на 14 месяца", 0, 0, false, "Russian-like month name"},

		{"на  14  листопада", 14, 11, true, "extra spaces"},
		{"на\t14\tлистопада", 14, 11, true, "tabs instead of spaces"},
	}

	for _, tt := range tests {
		day, month, ok := extractDayAndMonth(tt.title)
		if ok != tt.shouldOk {
			t.Errorf("%s: extractDayAndMonth(%q) ok = %v, want %v", tt.desc, tt.title, ok, tt.shouldOk)
		}
		if ok {
			if day != tt.expectedDay {
				t.Errorf("%s: extractDayAndMonth(%q) day = %d, want %d", tt.desc, tt.title, day, tt.expectedDay)
			}
			if month != tt.expectedMonth {
				t.Errorf("%s: extractDayAndMonth(%q) month = %d, want %d", tt.desc, tt.title, month, tt.expectedMonth)
			}
		}
	}
}

func TestCalculateSimilarity(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		minScore float64
		desc     string
	}{
		{"листопада", "листопада", 1.0, "identical strings"},
		{"листопада", "листопаде", 0.8, "one character different"},
		{"листопада", "листопда", 0.8, "one character missing"},
		{"січня", "сiчня", 0.8, "one character substitution"},
		{"", "", 1.0, "both empty"},
		{"abc", "xyz", 0.0, "completely different"},
	}

	for _, tt := range tests {
		got := calculateSimilarity(tt.s1, tt.s2)
		if got < tt.minScore {
			t.Errorf("%s: calculateSimilarity(%q, %q) = %f, want >= %f", tt.desc, tt.s1, tt.s2, got, tt.minScore)
		}
	}
}

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
		desc     string
	}{
		{"", "", 0, "both empty"},
		{"a", "", 1, "one empty"},
		{"", "a", 1, "other empty"},
		{"abc", "abc", 0, "identical"},
		{"abc", "ab", 1, "one deletion"},
		{"ab", "abc", 1, "one insertion"},
		{"abc", "adc", 1, "one substitution"},
		{"листопада", "листопаде", 1, "one character different"},
		{"листопада", "листопда", 2, "one character missing"},
	}

	for _, tt := range tests {
		got := levenshteinDistance(tt.s1, tt.s2)
		if got != tt.expected {
			t.Errorf("%s: levenshteinDistance(%q, %q) = %d, want %d", tt.desc, tt.s1, tt.s2, got, tt.expected)
		}
	}
}

func TestMin3(t *testing.T) {
	tests := []struct {
		a, b, c  int
		expected int
	}{
		{1, 2, 3, 1},
		{3, 2, 1, 1},
		{2, 1, 3, 1},
		{1, 1, 1, 1},
		{-1, 0, 1, -1},
		{0, 0, 0, 0},
	}

	for _, tt := range tests {
		got := min3(tt.a, tt.b, tt.c)
		if got != tt.expected {
			t.Errorf("min3(%d, %d, %d) = %d, want %d", tt.a, tt.b, tt.c, got, tt.expected)
		}
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{1, 2, 2},
		{2, 1, 2},
		{1, 1, 1},
		{-1, 0, 0},
		{0, 0, 0},
		{-5, -3, -3},
	}

	for _, tt := range tests {
		got := max(tt.a, tt.b)
		if got != tt.expected {
			t.Errorf("max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
		}
	}
}