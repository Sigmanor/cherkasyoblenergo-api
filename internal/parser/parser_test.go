package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cherkasyoblenergo-api/internal/models"
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
