// Copyright 2026 mvanhorn. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"encoding/binary"
	"strings"
	"testing"
)

// buildSimpleBlob constructs a minimal typedstream-prefixed blob containing a
// single UTF-8 string with the given length-prefix encoding. Used by tests
// to exercise the decoder against known-shape inputs without depending on
// captured chat.db rows.
func buildSimpleBlob(text string, prefixVariant byte) []byte {
	var b []byte
	b = append(b, []byte("streamtyped")...)
	// Pad with a few non-tag bytes so the scanner has to skip past the header
	// before finding the string tag.
	b = append(b, 0x04, 0x0b, 0x06)
	b = append(b, typedStreamStringTag)

	textBytes := []byte(text)
	switch prefixVariant {
	case 0: // direct single-byte length (length < 0x81)
		b = append(b, byte(len(textBytes)))
	case lengthPrefix1Byte:
		b = append(b, lengthPrefix1Byte, byte(len(textBytes)))
	case lengthPrefix2Byte:
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, uint16(len(textBytes)))
		b = append(b, lengthPrefix2Byte)
		b = append(b, buf...)
	case lengthPrefix4Byte:
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(len(textBytes)))
		b = append(b, lengthPrefix4Byte)
		b = append(b, buf...)
	}
	b = append(b, textBytes...)
	return b
}

func TestDecodeAttributedBody_Empty(t *testing.T) {
	text, source := decodeAttributedBody(nil)
	if text != "" || source != textSourceUnrecoverable {
		t.Errorf("nil blob: got (%q, %q), want (\"\", %q)", text, source, textSourceUnrecoverable)
	}

	text, source = decodeAttributedBody([]byte{})
	if text != "" || source != textSourceUnrecoverable {
		t.Errorf("empty blob: got (%q, %q), want (\"\", %q)", text, source, textSourceUnrecoverable)
	}
}

func TestDecodeAttributedBody_NonTypedStream(t *testing.T) {
	blob := []byte("this is not a typedstream blob")
	text, source := decodeAttributedBody(blob)
	if text != "" || source != textSourceUnrecoverable {
		t.Errorf("non-typedstream prefix: got (%q, %q), want (\"\", %q)", text, source, textSourceUnrecoverable)
	}
}

func TestDecodeAttributedBody_ShortASCII(t *testing.T) {
	blob := buildSimpleBlob("hello world", 0)
	text, source := decodeAttributedBody(blob)
	if text != "hello world" || source != textSourceDecoded {
		t.Errorf("short ASCII: got (%q, %q), want (\"hello world\", %q)", text, source, textSourceDecoded)
	}
}

func TestDecodeAttributedBody_MultiByteUTF8(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"emoji", "🎉🎊👋"},
		{"accented", "café résumé"},
		{"cjk", "你好世界"},
		{"mixed", "hi 👋 from café 你好"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			blob := buildSimpleBlob(tc.in, 0)
			text, source := decodeAttributedBody(blob)
			if text != tc.in || source != textSourceDecoded {
				t.Errorf("got (%q, %q), want (%q, %q)", text, source, tc.in, textSourceDecoded)
			}
		})
	}
}

func TestDecodeAttributedBody_Length1Byte(t *testing.T) {
	// Build a string long enough to require the 1-byte length prefix
	// (length 0x90 = 144, which is >= 0x81 so direct encoding can't be used).
	msg := strings.Repeat("a", 144)
	blob := buildSimpleBlob(msg, lengthPrefix1Byte)
	text, source := decodeAttributedBody(blob)
	if text != msg || source != textSourceDecoded {
		t.Errorf("1-byte length: got source=%q text-len=%d, want source=%q text-len=%d",
			source, len(text), textSourceDecoded, len(msg))
	}
}

func TestDecodeAttributedBody_Length2Byte(t *testing.T) {
	// Build a longer string that requires the 2-byte length prefix.
	msg := strings.Repeat("hello ", 100) // 600 bytes
	blob := buildSimpleBlob(msg, lengthPrefix2Byte)
	text, source := decodeAttributedBody(blob)
	if text != msg || source != textSourceDecoded {
		t.Errorf("2-byte length: got source=%q text-len=%d, want source=%q text-len=%d",
			source, len(text), textSourceDecoded, len(msg))
	}
}

func TestDecodeAttributedBody_Length4Byte(t *testing.T) {
	// Very large message requiring the 4-byte length prefix.
	msg := strings.Repeat("xyz ", 20000) // 80000 bytes
	blob := buildSimpleBlob(msg, lengthPrefix4Byte)
	text, source := decodeAttributedBody(blob)
	if text != msg || source != textSourceDecoded {
		t.Errorf("4-byte length: got source=%q text-len=%d, want source=%q text-len=%d",
			source, len(text), textSourceDecoded, len(msg))
	}
}

func TestDecodeAttributedBody_TruncatedMidLength(t *testing.T) {
	// Build a 2-byte-prefix blob then chop off in the middle of the length.
	blob := buildSimpleBlob("some message body", lengthPrefix2Byte)
	// Find the prefix marker and truncate to immediately after it.
	idx := -1
	for i := 0; i < len(blob)-1; i++ {
		if blob[i] == typedStreamStringTag && blob[i+1] == lengthPrefix2Byte {
			idx = i + 2 // keep tag + marker, drop length bytes
			break
		}
	}
	if idx < 0 {
		t.Fatal("test setup error: prefix marker not found")
	}
	truncated := blob[:idx]
	text, source := decodeAttributedBody(truncated)
	// We don't require unrecoverable specifically — the scanner may find an
	// earlier 0x2B byte and try to decode there too. The contract is just
	// "don't crash and don't return obviously wrong content".
	if source == textSourceDecoded && text == "some message body" {
		t.Errorf("truncated blob should not have decoded full text, got (%q, %q)", text, source)
	}
}

func TestDecodeAttributedBody_LengthExceedsBuffer(t *testing.T) {
	// Construct a blob whose claimed length is larger than remaining bytes.
	var blob []byte
	blob = append(blob, []byte("streamtyped")...)
	blob = append(blob, 0x04, 0x0b)
	blob = append(blob, typedStreamStringTag)
	blob = append(blob, lengthPrefix2Byte)
	// Claim 10000 bytes but only provide 5.
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, 10000)
	blob = append(blob, buf...)
	blob = append(blob, []byte("short")...)

	text, source := decodeAttributedBody(blob)
	if source == textSourceDecoded {
		t.Errorf("over-length claim should not decode, got (%q, %q)", text, source)
	}
}

func TestDecodeAttributedBody_InvalidUTF8(t *testing.T) {
	// Build a blob with a valid prefix but invalid UTF-8 bytes.
	var blob []byte
	blob = append(blob, []byte("streamtyped")...)
	blob = append(blob, 0x04, 0x0b)
	blob = append(blob, typedStreamStringTag)
	blob = append(blob, 0x05)                                  // direct length 5
	blob = append(blob, 0xff, 0xfe, 0xfd, 0xfc, 0xfb)          // invalid UTF-8
	text, source := decodeAttributedBody(blob)
	if source == textSourceDecoded {
		t.Errorf("invalid UTF-8 should not decode, got (%q, %q)", text, source)
	}
}

func TestDecodeAttributedBody_FiltersClassName(t *testing.T) {
	// "NSMutableAttributedString" is a known class name that should be
	// filtered out by looksLikeMessageText. Build a blob that has it
	// as the first decodable candidate and ensure the decoder rejects it
	// (or skips past it to find no further valid string).
	blob := buildSimpleBlob("NSMutableAttributedString", 0)
	text, source := decodeAttributedBody(blob)
	if source == textSourceDecoded {
		t.Errorf("class name should be filtered, got (%q, %q)", text, source)
	}
}

func TestDecodeAttributedBody_LooksLikeMessageText(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"normal text", "hey are you free for lunch", true},
		{"with emoji", "yes! 🎉", true},
		{"empty", "", false},
		{"class name NSString", "NSString", false},
		{"class name NSMutableAttributedString", "NSMutableAttributedString", false},
		{"control chars dominant", "\x01\x02\x03\x04", false},
		{"low control byte ratio", "hello\nworld\tfine", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := looksLikeMessageText(tc.in)
			if got != tc.want {
				t.Errorf("looksLikeMessageText(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
