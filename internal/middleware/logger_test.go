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

func TestMaskAPIKey_ExactlyTen(t *testing.T) {
	apiKey := "1234567890"
	masked := maskAPIKey(apiKey)
	if masked != apiKey {
		t.Errorf("Expected %q unchanged, got %q", apiKey, masked)
	}
}

func TestMaskAPIKey_Eleven(t *testing.T) {
	apiKey := "12345678901"
	masked := maskAPIKey(apiKey)
	expected := "12345" + "*" + "78901"
	if masked != expected {
		t.Errorf("Expected %q, got %q", expected, masked)
	}
}
