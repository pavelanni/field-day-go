package morse

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 44100
)

// Morse code table
var morseCodeMap = map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".", 'F': "..-.",
	'G': "--.", 'H': "....", 'I': "..", 'J': ".---", 'K': "-.-", 'L': ".-..",
	'M': "--", 'N': "-.", 'O': "---", 'P': ".--.", 'Q': "--.-", 'R': ".-.",
	'S': "...", 'T': "-", 'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-",
	'Y': "-.--", 'Z': "--..",
	'0': "-----", '1': ".----", '2': "..---", '3': "...--", '4': "....-",
	'5': ".....", '6': "-....", '7': "--...", '8': "---..", '9': "----.",
	'/': "-..-.", '?': "..--..", '.': ".-.-.-", ',': "--..--",
}

type Player struct {
	audioContext *audio.Context
	freq         int
	wpm          int
	samples      morseAudio
}

func NewPlayer(freq int, wpm int) *Player {
	acontext := audio.NewContext(sampleRate)
	samples := newMorseAudio(wpm, freq)
	return &Player{audioContext: acontext, freq: freq, wpm: wpm, samples: *samples}
}

func (p *Player) Play(text string) error {

	// Generate Morse code audio
	samples, totalSamples := p.generateMorseAudio(text)

	// Create a buffer and write WAV data
	buf := &bytes.Buffer{}
	writeWavHeader(buf, totalSamples*2, sampleRate)
	// Write PCM data
	for _, sample := range samples {
		err := binary.Write(buf, binary.LittleEndian, sample)
		if err != nil {
			return err
		}
	}

	// Create a reader from the buffer
	reader := bytes.NewReader(buf.Bytes())

	// Play the sound
	audioPlayer, err := wav.DecodeWithSampleRate(sampleRate, reader)
	if err != nil {
		return err
	}

	player, err := p.audioContext.NewPlayer(audioPlayer)
	if err != nil {
		return err
	}

	fmt.Printf("Playing '%s' in Morse code at %d WPM, %d Hz\n", text, p.wpm, p.freq)
	player.Play()

	// Calculate total duration and wait for playback to complete
	totalDuration := time.Duration(float64(totalSamples) / float64(sampleRate) * float64(time.Second))
	time.Sleep(totalDuration)

	return nil
}

// Pre-generated audio samples for each character
type morseAudio struct {
	dotSamples  []int16
	dashSamples []int16
	elementGap  []int16
	charGap     []int16
	wordGap     []int16
	charSamples map[rune][]int16
}

// calculateMorseTiming calculates timing from WPM using PARIS standard
func calculateMorseTiming(wpm int) (dotDuration, dashDuration, elementGap, charGap, wordGap int) {
	if wpm <= 0 {
		wpm = 20 // Default to 20 WPM
	}

	// 1 time unit duration in milliseconds
	timeUnit := 60000 / (wpm * 50) // 60 seconds * 1000 ms / (wpm * 50 units per PARIS)

	dotDuration = timeUnit
	dashDuration = timeUnit * 3
	elementGap = timeUnit
	charGap = timeUnit * 3
	wordGap = timeUnit * 7

	return
}

// newMorseAudio creates a new morseAudio instance with pre-generated samples
func newMorseAudio(wpm int, freq int) *morseAudio {
	dotDuration, dashDuration, elementGap, charGap, wordGap := calculateMorseTiming(wpm)

	// Convert durations from milliseconds to samples
	dotSamples := int(float64(dotDuration) * sampleRate / 1000)
	dashSamples := int(float64(dashDuration) * sampleRate / 1000)
	elementGapSamples := int(float64(elementGap) * sampleRate / 1000)
	charGapSamples := int(float64(charGap) * sampleRate / 1000)
	wordGapSamples := int(float64(wordGap) * sampleRate / 1000)

	// Generate basic elements
	dot := make([]int16, dotSamples)
	dash := make([]int16, dashSamples)
	elementGapAudio := make([]int16, elementGapSamples)
	charGapAudio := make([]int16, charGapSamples)
	wordGapAudio := make([]int16, wordGapSamples)

	// Generate tone for dot and dash
	for i := 0; i < dotSamples; i++ {
		dot[i] = int16(math.Sin(2*math.Pi*float64(freq)*float64(i)/float64(sampleRate)) * 32767)
	}
	for i := 0; i < dashSamples; i++ {
		dash[i] = int16(math.Sin(2*math.Pi*float64(freq)*float64(i)/float64(sampleRate)) * 32767)
	}

	// Pre-generate samples for each character
	charSamples := make(map[rune][]int16)
	for char, morse := range morseCodeMap {
		var samples []int16
		for i, element := range morse {
			if element == '.' {
				samples = append(samples, dot...)
			} else if element == '-' {
				samples = append(samples, dash...)
			}
			if i < len(morse)-1 {
				samples = append(samples, elementGapAudio...)
			}
		}
		charSamples[char] = samples
	}

	return &morseAudio{
		dotSamples:  dot,
		dashSamples: dash,
		elementGap:  elementGapAudio,
		charGap:     charGapAudio,
		wordGap:     wordGapAudio,
		charSamples: charSamples,
	}
}

// generateMorseAudio generates audio for a given text in Morse code
func (p *Player) generateMorseAudio(text string) ([]int16, int) {
	morse := p.samples

	// Calculate total duration needed
	totalSamples := 0
	for i, char := range strings.ToUpper(text) {
		if char == ' ' {
			totalSamples += len(morse.wordGap)
			continue
		}

		if samples, ok := morse.charSamples[char]; ok {
			totalSamples += len(samples)
			if i < len(text)-1 && text[i+1] != ' ' {
				totalSamples += len(morse.charGap)
			}
		}
	}

	// Generate the audio samples
	samples := make([]int16, totalSamples)
	currentSample := 0

	for i, char := range strings.ToUpper(text) {
		if char == ' ' {
			copy(samples[currentSample:], morse.wordGap)
			currentSample += len(morse.wordGap)
			continue
		}

		if charSamples, ok := morse.charSamples[char]; ok {
			copy(samples[currentSample:], charSamples)
			currentSample += len(charSamples)

			// Add character gap if not the last character
			if i < len(text)-1 && text[i+1] != ' ' {
				copy(samples[currentSample:], morse.charGap)
				currentSample += len(morse.charGap)
			}
		}
	}

	return samples, totalSamples
}

func writeWavHeader(w *bytes.Buffer, dataSize int, sampleRate int) {
	// RIFF header
	w.Write([]byte("RIFF"))
	binary.Write(w, binary.LittleEndian, uint32(36+dataSize))
	w.Write([]byte("WAVE"))

	// fmt chunk
	w.Write([]byte("fmt "))
	binary.Write(w, binary.LittleEndian, uint32(16)) // fmt chunk size
	binary.Write(w, binary.LittleEndian, uint16(1))  // audio format (1 for PCM)
	binary.Write(w, binary.LittleEndian, uint16(1))  // number of channels
	binary.Write(w, binary.LittleEndian, uint32(sampleRate))
	binary.Write(w, binary.LittleEndian, uint32(sampleRate*2)) // byte rate
	binary.Write(w, binary.LittleEndian, uint16(2))            // block align
	binary.Write(w, binary.LittleEndian, uint16(16))           // bits per sample

	// data chunk
	w.Write([]byte("data"))
	binary.Write(w, binary.LittleEndian, uint32(dataSize))
}

// GenerateWav generates a WAV file as a byte slice for the given text using default WPM and frequency.
func GenerateWav(text string) ([]byte, error) {
	const defaultWPM = 15
	const defaultFreq = 600
	// Create a Player (audioContext is not used here)
	samples := newMorseAudio(defaultWPM, defaultFreq)
	p := &Player{freq: defaultFreq, wpm: defaultWPM, samples: *samples}

	audioSamples, totalSamples := p.generateMorseAudio(text)
	buf := &bytes.Buffer{}
	writeWavHeader(buf, totalSamples*2, sampleRate)
	for _, sample := range audioSamples {
		err := binary.Write(buf, binary.LittleEndian, sample)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
