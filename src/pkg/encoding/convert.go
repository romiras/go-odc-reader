// Package encoding provides character encoding conversion utilities.
package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"odcread/pkg/oberon"
)

// ConvertLatin1 converts ISO-8859-1 (Latin-1) encoded bytes to UTF-8 string.
// It also normalizes line endings by converting \r to \n.
func ConvertLatin1(input []oberon.ShortChar) (string, error) {
	// Find the null terminator
	length := len(input)
	for i, ch := range input {
		if ch == 0 {
			length = i
			break
		}
	}

	// Convert to byte slice (excluding null terminator)
	data := input[:length]

	// ISO-8859-1 decoder
	decoder := charmap.ISO8859_1.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(data), decoder)

	result, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to convert Latin-1: %w", err)
	}

	// Normalize line endings
	output := strings.ReplaceAll(string(result), "\r", "\n")
	return output, nil
}

// ConvertUCS2 converts UCS-2 (16-bit Unicode) encoded data to UTF-8 string.
// It also normalizes line endings by converting \r to \n.
func ConvertUCS2(input []oberon.Char) (string, error) {
	// Find the null terminator
	length := len(input)
	for i, ch := range input {
		if ch == 0 {
			length = i
			break
		}
	}

	// Convert uint16 slice to bytes (little-endian)
	buf := new(bytes.Buffer)
	for _, ch := range input[:length] {
		if err := binary.Write(buf, binary.LittleEndian, ch); err != nil {
			return "", fmt.Errorf("failed to write UCS-2 char: %w", err)
		}
	}

	// UTF-16 decoder (UCS-2 is a subset of UTF-16)
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	reader := transform.NewReader(buf, decoder)

	result, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to convert UCS-2: %w", err)
	}

	// Normalize line endings
	output := strings.ReplaceAll(string(result), "\r", "\n")
	return output, nil
}
