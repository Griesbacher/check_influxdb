package Units

import (
	"testing"
)

var byteToString = []struct {
	b        float64
	expected string
}{
	{1, "1B"},
	{1024, "1KiB"},
	{1024 * 2, "2KiB"},
	{1024 * 2.5, "2.5KiB"},
	{1024 * 1024, "1MiB"},
	{1024 * 1024 * 1024, "1GiB"},
	{1024 * 1024 * 1024 * 1024, "1TiB"},
	{1024 * 1024 * 1024 * 1024 * 1024, "1PiB"},
	{1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1EiB"},
	{1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1ZiB"},
	{1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1YiB"},
}

func TestByteSize_String(t *testing.T) {
	for i, data := range byteToString {
		result := ByteSize(data.b).String()
		if result != data.expected {
			t.Errorf("%d - Expected: %s, got: %s", i, data.expected, result)
		}
	}
}
