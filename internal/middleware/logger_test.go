package middleware

import "testing"

func TestMaskAPIKey_Long(t *testing.T) {
	apiKey := "1234567890ABCDEF"
	masked := maskAPIKey(apiKey)
	expectedPrefix := "12345"
	expectedSuffix := "BCDEF"
	if masked[:5] != expectedPrefix {
		t.Errorf("Expected prefix %q, got %q", expectedPrefix, masked[:5])
	}
	if masked[len(masked)-5:] != expectedSuffix {
		t.Errorf("Expected suffix %q, got %q", expectedSuffix, masked[len(masked)-5:])
	}
}

func TestMaskAPIKey_Short(t *testing.T) {
	apiKey := "short"
	masked := maskAPIKey(apiKey)
	if masked != apiKey {
		t.Errorf("Expected %q, got %q", apiKey, masked)
	}
}
