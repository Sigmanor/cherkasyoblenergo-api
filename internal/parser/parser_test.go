package parser

import (
	"cherkasyoblenergo-api/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseScheduleFromParagraphs_ToleratesExtraDotAndSpaces(t *testing.T) {
	html := `
<p>1.1 00:00 - 02:00</p>
<p>2.1. 06:00 - 08:00, 11:00 - 13:00, 15:00 - 17:00, 19:00 - 21:00</p>
<p>2.2&nbsp;06:00 - 08:00, 12:00 - 14:00, 16:00 - 18:00, 20:00 - 22:00</p>
<p>4.1 02:00 - 04:00, 08:00 - 10:00, 12:00 - 14:00, 16:00 - 18:00, 20:00 - 22:00</p>
<p>4.2 02:00 - 04:00, 08:00 - 10:00, 14:00 - 16:00, 18:00 - 20:00, 22:00 - 00:00&nbsp;</p>
`

	result, found := parseScheduleFromParagraphs(html)

	require.True(t, found, "paragraph parser should find schedule rows")
	assert.Equal(t, "00:00 - 02:00", result.Col1_1)
	assert.Equal(t, "06:00 - 08:00, 11:00 - 13:00, 15:00 - 17:00, 19:00 - 21:00", result.Col2_1)
	assert.Equal(t, "06:00 - 08:00, 12:00 - 14:00, 16:00 - 18:00, 20:00 - 22:00", result.Col2_2)
	assert.Equal(t, "02:00 - 04:00, 08:00 - 10:00, 12:00 - 14:00, 16:00 - 18:00, 20:00 - 22:00", result.Col4_1)
	assert.Equal(t, "02:00 - 04:00, 08:00 - 10:00, 14:00 - 16:00, 18:00 - 20:00, 22:00 - 00:00", result.Col4_2)
}

func TestNormalizeTimeRanges_TrimsTrailingSeparators(t *testing.T) {
	assert.Equal(t, "02:00 - 04:00", normalizeTimeRanges(" 02:00 - 04:00, "))
}

func TestParseScheduleFromParagraphs_HandlesDivBlocks(t *testing.T) {
	html := `
<div data-block="true"><div>1.1 00:00 - 02:00, 06:00 - 08:00</div></div>
<div data-block="true"><div>2.1 08:00 - 10:00, 13:00 - 15:00</div></div>
<div data-block="true"><div>6.2 06:00 - 08:00, 11:00 - 13:00, 16:00 - 18:00, 21:00 - 23:00</div></div>
`

	result, found := parseScheduleFromParagraphs(html)

	require.True(t, found, "div parser should find schedule rows")
	assert.Equal(t, "00:00 - 02:00, 06:00 - 08:00", result.Col1_1)
	assert.Equal(t, "08:00 - 10:00, 13:00 - 15:00", result.Col2_1)
	assert.Equal(t, "06:00 - 08:00, 11:00 - 13:00, 16:00 - 18:00, 21:00 - 23:00", result.Col6_2)
}

func TestHasScheduleData(t *testing.T) {
	assert.False(t, hasScheduleData(models.Schedule{}), "empty schedule should be treated as missing")
	assert.True(t, hasScheduleData(models.Schedule{Col3_2: "01:00 - 02:00"}), "any non-empty column should count")
	assert.True(t, hasScheduleData(models.Schedule{Col1_1: " 01:00 - 02:00 "}), "trimming should still detect data")
}

func TestParseScheduleFromParagraphs_HandlesColonFormat(t *testing.T) {
	html := `
<p>1.1: 00:30 – 04:00, 06:00 – 10:00, 12:00 – 16:00, 18:00 – 22:00</p>
<p>1.2: 01:30 – 05:30, 07:30 – 11:30, 13:30 – 17:30, 19:30 – 22:30</p>
<p>2.1: 00:00 – 00:30, 03:00 – 06:30, 08:30 – 12:30, 14:30 – 18:30, 20:30 – 00:00</p>
<p>6.2: 00:00 – 03:30, 06:00 – 09:30, 11:30 – 15:30, 17:30 – 21:30, 23:30 – 00:00</p>
`

	result, found := parseScheduleFromParagraphs(html)

	require.True(t, found, "paragraph parser should find schedule rows with colon format")
	assert.Equal(t, "00:30 – 04:00, 06:00 – 10:00, 12:00 – 16:00, 18:00 – 22:00", result.Col1_1)
	assert.Equal(t, "01:30 – 05:30, 07:30 – 11:30, 13:30 – 17:30, 19:30 – 22:30", result.Col1_2)
	assert.Equal(t, "00:00 – 00:30, 03:00 – 06:30, 08:30 – 12:30, 14:30 – 18:30, 20:30 – 00:00", result.Col2_1)
	assert.Equal(t, "00:00 – 03:30, 06:00 – 09:30, 11:30 – 15:30, 17:30 – 21:30, 23:30 – 00:00", result.Col6_2)
}

func TestContainsSchedulePatterns_HandlesColonFormat(t *testing.T) {
	// Old format without colon
	assert.True(t, containsSchedulePatterns("1.1 00:30 - 04:00"), "should match old format without colon")
	// New format with colon
	assert.True(t, containsSchedulePatterns("1.1: 00:30 – 04:00"), "should match new format with colon")
	assert.True(t, containsSchedulePatterns("<p>2.2: 06:00 – 08:00</p>"), "should match new format in HTML")
}
