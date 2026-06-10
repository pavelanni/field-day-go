package morse

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	freq := 600
	wpm := 15
	player := NewPlayer(freq, wpm)

	if player == nil {
		t.Error("NewPlayer() returned nil")
	}
	if player.freq != freq {
		t.Errorf("NewPlayer().freq = %d, want %d", player.freq, freq)
	}
	if player.wpm != wpm {
		t.Errorf("NewPlayer().wpm = %d, want %d", player.wpm, wpm)
	}
}

func TestCalculateMorseTiming(t *testing.T) {
	testCases := []struct {
		wpm         int
		expectedDot int
	}{
		{20, 60},  // 60000 / (20 * 50) = 60ms
		{15, 80},  // 60000 / (15 * 50) = 80ms
		{10, 120}, // 60000 / (10 * 50) = 120ms
		{0, 60},   // Default to 20 WPM
	}

	for _, tc := range testCases {
		dot, dash, elementGap, charGap, wordGap := calculateMorseTiming(tc.wpm)

		if dot != tc.expectedDot {
			t.Errorf("calculateMorseTiming(%d) dot = %d, want %d", tc.wpm, dot, tc.expectedDot)
		}
		if dash != dot*3 {
			t.Errorf("calculateMorseTiming(%d) dash = %d, want %d (3x dot)", tc.wpm, dash, dot*3)
		}
		if elementGap != dot {
			t.Errorf("calculateMorseTiming(%d) elementGap = %d, want %d (same as dot)", tc.wpm, elementGap, dot)
		}
		if charGap != dot*3 {
			t.Errorf("calculateMorseTiming(%d) charGap = %d, want %d (3x dot)", tc.wpm, charGap, dot*3)
		}
		if wordGap != dot*7 {
			t.Errorf("calculateMorseTiming(%d) wordGap = %d, want %d (7x dot)", tc.wpm, wordGap, dot*7)
		}
	}
}

func TestMorseCodeMap(t *testing.T) {
	// Test some basic morse code mappings
	testCases := map[rune]string{
		'A': ".-",
		'B': "-...",
		'S': "...",
		'O': "---",
		'1': ".----",
		'0': "-----",
		'?': "..--..",
	}

	for char, expectedMorse := range testCases {
		if morse, exists := morseCodeMap[char]; !exists {
			t.Errorf("morseCodeMap missing entry for '%c'", char)
		} else if morse != expectedMorse {
			t.Errorf("morseCodeMap['%c'] = %s, want %s", char, morse, expectedMorse)
		}
	}
}

func TestGenerateWav(t *testing.T) {
	testCases := []string{
		"SOS",
		"W1AW",
		"73",
		"A",
		"",
	}

	for _, text := range testCases {
		wavData, err := GenerateWav(text)
		if err != nil {
			t.Errorf("GenerateWav(%q) error = %v, want nil", text, err)
			continue
		}

		if len(wavData) == 0 {
			t.Errorf("GenerateWav(%q) returned empty data", text)
			continue
		}

		// Check WAV header
		if len(wavData) < 44 {
			t.Errorf("GenerateWav(%q) returned data too short for WAV header", text)
			continue
		}

		// Check RIFF header
		if string(wavData[0:4]) != "RIFF" {
			t.Errorf("GenerateWav(%q) missing RIFF header", text)
		}

		// Check WAVE format
		if string(wavData[8:12]) != "WAVE" {
			t.Errorf("GenerateWav(%q) missing WAVE format", text)
		}
	}
}

func TestPlayerGenerateMorseAudio(t *testing.T) {
	player := NewPlayer(600, 20)

	testCases := []struct {
		text    string
		minSize int
		maxSize int
	}{
		{"A", 1000, 20000},     // Single letter
		{"SOS", 10000, 100000}, // Classic distress call
		{"", 0, 100},           // Empty string
		{" ", 10000, 30000},    // Single space (word gap)
	}

	for _, tc := range testCases {
		samples, totalSamples := player.generateMorseAudio(tc.text)

		if len(samples) != totalSamples {
			t.Errorf("generateMorseAudio(%q) len(samples) = %d, want %d", tc.text, len(samples), totalSamples)
		}

		if totalSamples < tc.minSize || totalSamples > tc.maxSize {
			t.Errorf("generateMorseAudio(%q) totalSamples = %d, want between %d and %d", tc.text, totalSamples, tc.minSize, tc.maxSize)
		}
	}
}

func TestPlayerPlay(t *testing.T) {
	player := NewPlayer(600, 15)

	// Test various inputs
	testCases := []string{
		"W1AW",
		"73",
		"SOS",
		"HELLO WORLD",
		"",
	}

	for _, text := range testCases {
		err := player.Play(text)
		if err != nil {
			t.Errorf("Player.Play(%q) error = %v, want nil", text, err)
		}
	}
}

func TestMorseCodeMapCompleteness(t *testing.T) {
	// Test that we have morse codes for all letters and numbers
	for char := 'A'; char <= 'Z'; char++ {
		if _, exists := morseCodeMap[char]; !exists {
			t.Errorf("morseCodeMap missing letter '%c'", char)
		}
	}

	for char := '0'; char <= '9'; char++ {
		if _, exists := morseCodeMap[char]; !exists {
			t.Errorf("morseCodeMap missing digit '%c'", char)
		}
	}
}

func TestCaseInsensitivity(t *testing.T) {
	player := NewPlayer(600, 20)

	// Test that uppercase and lowercase generate the same audio
	upperSamples, upperTotal := player.generateMorseAudio("SOS")
	lowerSamples, lowerTotal := player.generateMorseAudio("sos")

	if upperTotal != lowerTotal {
		t.Errorf("Case sensitivity issue: uppercase total = %d, lowercase total = %d", upperTotal, lowerTotal)
	}

	if len(upperSamples) != len(lowerSamples) {
		t.Errorf("Case sensitivity issue: different sample lengths")
	}

	// Check that samples are identical
	for i := 0; i < len(upperSamples) && i < len(lowerSamples); i++ {
		if upperSamples[i] != lowerSamples[i] {
			t.Errorf("Case sensitivity issue: samples differ at position %d", i)
			break
		}
	}
}
