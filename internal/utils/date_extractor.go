package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type tokenType int

const (
	tokenNumber tokenType = iota
	tokenWord
	tokenPunctuation
	tokenWhitespace
)

type token struct {
	typ   tokenType
	value string
}

func ExtractScheduleDateFromTitle(title string) string {
	day, month, ok := extractDayAndMonth(title)
	if !ok {
		return ""
	}
	
	now := time.Now()
	year := now.Year()
	
	// If extracted month is after current month, assume previous year
	if month > int(now.Month()) {
		year--
	}
	
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func tokenizeText(text string) []token {
	var tokens []token
	runes := []rune(strings.ToLower(text))
	i := 0

	for i < len(runes) {
		r := runes[i]

		if unicode.IsSpace(r) {
			j := i
			for j < len(runes) && unicode.IsSpace(runes[j]) {
				j++
			}
			i = j
			continue
		}

		if unicode.IsDigit(r) {
			j := i
			for j < len(runes) && unicode.IsDigit(runes[j]) {
				j++
			}
			tokens = append(tokens, token{
				typ:   tokenNumber,
				value: string(runes[i:j]),
			})
			i = j
			continue
		}

		if unicode.IsLetter(r) {
			j := i
			for j < len(runes) && unicode.IsLetter(runes[j]) {
				j++
			}
			tokens = append(tokens, token{
				typ:   tokenWord,
				value: string(runes[i:j]),
			})
			i = j
			continue
		}

		tokens = append(tokens, token{
			typ:   tokenPunctuation,
			value: string(r),
		})
		i++
	}

	return tokens
}

func isCyrillic(s string) bool {
	for _, r := range s {
		if !unicode.Is(unicode.Cyrillic, r) {
			return false
		}
	}
	return true
}

func findDateInTokens(tokens []token) (day int, month int, ok bool) {
	keywordLike := map[string]bool{"на": true, "нa": true, "за": true, "зa": true, "з": true}
	for i := 0; i < len(tokens)-1; i++ {
		if tokens[i].typ == tokenNumber {
			prevWordIdx := -1
			for j := i - 1; j >= 0; j-- {
				if tokens[j].typ == tokenWord {
					prevWordIdx = j
					break
				}
				if tokens[j].typ != tokenPunctuation {
					break
				}
			}

			if prevWordIdx >= 0 && keywordLike[tokens[prevWordIdx].value] {
				if !isCyrillic(tokens[prevWordIdx].value) {
					continue
				}
			}

			dayNum, err := strconv.Atoi(tokens[i].value)
			if err != nil || dayNum < 1 || dayNum > 31 {
				continue
			}

			if i+1 < len(tokens) && tokens[i+1].typ == tokenWord {
				if !isCyrillic(tokens[i+1].value) {
					continue
				}

				monthNum, monthOk := parseUkrainianMonth(tokens[i+1].value)
				if monthOk {
					return dayNum, monthNum, true
				}
			}
		}
	}

	keywords := map[string]bool{"на": true, "за": true}
	for i := 0; i < len(tokens)-2; i++ {
		if tokens[i].typ == tokenWord && keywords[tokens[i].value] {
			if !isCyrillic(tokens[i].value) {
				continue
			}

			numIdx := i + 1
			for numIdx < len(tokens) && tokens[numIdx].typ == tokenPunctuation {
				numIdx++
			}

			if numIdx >= len(tokens) || tokens[numIdx].typ != tokenNumber {
				continue
			}

			monthIdx := numIdx + 1
			if monthIdx >= len(tokens) || tokens[monthIdx].typ != tokenWord {
				continue
			}

			dayNum, err := strconv.Atoi(tokens[numIdx].value)
			if err != nil || dayNum < 1 || dayNum > 31 {
				continue
			}

			if !isCyrillic(tokens[monthIdx].value) {
				continue
			}

			monthNum, monthOk := parseUkrainianMonth(tokens[monthIdx].value)
			if monthOk {
				return dayNum, monthNum, true
			}
		}
	}

	for i := 0; i < len(tokens)-3; i++ {
		if tokens[i].typ == tokenNumber {
			dayNum, err := strconv.Atoi(tokens[i].value)
			if err != nil || dayNum < 1 || dayNum > 31 {
				continue
			}

			if tokens[i+1].typ == tokenWord {
				if !isCyrillic(tokens[i+1].value) {
					continue
				}

				monthNum, monthOk := parseUkrainianMonth(tokens[i+1].value)
				if !monthOk {
					continue
				}

				if i+2 < len(tokens) && tokens[i+2].typ == tokenNumber {
					if i+3 < len(tokens) && tokens[i+3].typ == tokenWord && tokens[i+3].value == "року" {
						if !isCyrillic(tokens[i+3].value) {
							continue
						}
						return dayNum, monthNum, true
					}
				}
			}
		}
	}

	for i := 0; i < len(tokens); i++ {
		if tokens[i].typ == tokenWord && tokens[i].value == "з" {
			if !isCyrillic(tokens[i].value) {
				continue
			}

			for j := i + 1; j < len(tokens); j++ {
				if tokens[j].typ == tokenPunctuation && tokens[j].value == "," {
					numIdx := j + 1
					for numIdx < len(tokens) && tokens[numIdx].typ == tokenPunctuation {
						numIdx++
					}

					if numIdx >= len(tokens) || tokens[numIdx].typ != tokenNumber {
						break
					}

					dayNum, err := strconv.Atoi(tokens[numIdx].value)
					if err != nil || dayNum < 1 || dayNum > 31 {
						break
					}

					if numIdx+1 < len(tokens) && tokens[numIdx+1].typ == tokenWord {
						if !isCyrillic(tokens[numIdx+1].value) {
							break
						}

						monthNum, monthOk := parseUkrainianMonth(tokens[numIdx+1].value)
						if monthOk {
							return dayNum, monthNum, true
						}
					}
					break
				}
			}
		}
	}

	return 0, 0, false
}

func extractDayAndMonth(title string) (day int, month int, ok bool) {
	tokens := tokenizeText(title)
	return findDateInTokens(tokens)
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
