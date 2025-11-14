package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ExtractScheduleDateFromTitle(title string) string {
	day, month, ok := extractDayAndMonth(title)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%02d.%02d", day, month)
}

func extractDayAndMonth(title string) (day int, month int, ok bool) {
	patterns := []string{
		`на\s+(\d{1,2})\s+([а-яіїєґ]+)`,
		`(\d{1,2})\s+([а-яіїєґ]+)\s+\d{4}\s+року`,
		`з\s+\d{1,2}:\d{2},\s*(\d{1,2})\s+([а-яіїєґ]+)`,
	}
	
	lowerTitle := strings.ToLower(title)
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(lowerTitle)
		
		if len(matches) >= 3 {
			dayNum, err := strconv.Atoi(matches[1])
			if err != nil || dayNum < 1 || dayNum > 31 {
				continue
			}
			
			monthName := matches[2]
			monthNum, monthOk := parseUkrainianMonth(monthName)
			if !monthOk {
				continue
			}
			
			return dayNum, monthNum, true
		}
	}
	
	return 0, 0, false
}

func parseUkrainianMonth(monthName string) (int, bool) {
	monthName = strings.ToLower(strings.TrimSpace(monthName))
	
	monthMap := map[string]int{
		"січня":     1,
		"лютого":    2,
		"березня":   3,
		"квітня":    4,
		"травня":    5,
		"червня":    6,
		"липня":     7,
		"серпня":    8,
		"вересня":   9,
		"жовтня":    10,
		"листопада": 11,
		"грудня":    12,
	}
	
	if monthNum, exists := monthMap[monthName]; exists {
		return monthNum, true
	}
	
	const similarityThreshold = 0.8
	
	bestMatch := -1
	bestSimilarity := 0.0
	
	for knownMonth, monthNum := range monthMap {
		similarity := calculateSimilarity(monthName, knownMonth)
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			bestMatch = monthNum
		}
	}
	
	if bestSimilarity >= similarityThreshold {
		return bestMatch, true
	}
	
	return 0, false
}

func calculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	
	distance := levenshteinDistance(s1, s2)
	maxLen := max(len(s1), len(s2))
	
	if maxLen == 0 {
		return 1.0
	}
	
	return 1.0 - float64(distance)/float64(maxLen)
}

func levenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	
	if len(s1) == 0 {
		return len(s2)
	}
	
	if len(s2) == 0 {
		return len(s1)
	}
	
	rows := len(s1) + 1
	cols := len(s2) + 1
	matrix := make([][]int, rows)
	
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}
	
	for i := 0; i < rows; i++ {
		matrix[i][0] = i
	}
	for j := 0; j < cols; j++ {
		matrix[0][j] = j
	}
	
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			
			matrix[i][j] = min3(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}
	
	return matrix[rows-1][cols-1]
}

func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}