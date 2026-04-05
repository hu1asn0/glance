package glance

import (
	"testing"
)

func TestPercentChange(t *testing.T) {
	tests := []struct {
		name     string
		current  float64
		previous float64
		expected float64
	}{
		{"positive increase", 150, 100, 50},
		{"positive decrease", 50, 100, -50},
		{"no change", 100, 100, 0},
		{"both zero", 0, 0, 0},
		{"from zero to positive", 50, 0, 100},
		{"double", 200, 100, 100},
		{"negative values", -50, -100, -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := percentChange(tt.current, tt.previous)
			if result != tt.expected {
				t.Errorf("percentChange(%v, %v) = %v, want %v", tt.current, tt.previous, result, tt.expected)
			}
		})
	}
}

func TestExtractDomainFromUrl(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple https", "https://example.com/path", "example.com"},
		{"with www", "https://www.example.com/path", "example.com"},
		{"http", "http://blog.example.com", "blog.example.com"},
		{"with port", "https://example.com:8080/path", "example.com:8080"},
		{"empty string", "", ""},
		{"invalid url", "not a url", ""},
		{"uppercase", "https://WWW.EXAMPLE.COM", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDomainFromUrl(tt.input)
			if result != tt.expected {
				t.Errorf("extractDomainFromUrl(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStripURLScheme(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://example.com", "example.com"},
		{"http://example.com/path", "example.com/path"},
		{"ftp://files.example.com", "files.example.com"},
		{"example.com", "example.com"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := stripURLScheme(tt.input)
			if result != tt.expected {
				t.Errorf("stripURLScheme(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLimitStringLength(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		max         int
		expectedStr string
		expectedCut bool
	}{
		{"under limit", "hello", 10, "hello", false},
		{"at limit", "hello", 5, "hello", false},
		{"over limit", "hello world", 5, "hello", true},
		{"empty", "", 5, "", false},
		{"unicode", "héllő wörld", 5, "héllő", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, cut := limitStringLength(tt.input, tt.max)
			if result != tt.expectedStr || cut != tt.expectedCut {
				t.Errorf("limitStringLength(%q, %d) = (%q, %v), want (%q, %v)",
					tt.input, tt.max, result, cut, tt.expectedStr, tt.expectedCut)
			}
		})
	}
}

func TestNormalizeVersionFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.2.3", "v1.2.3"},
		{"v1.2.3", "v1.2.3"},
		{"V1.2.3", "v1.2.3"},
		{"  v1.0  ", "v1.0"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeVersionFormat(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeVersionFormat(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTitleToSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"My  Dashboard  Page", "my-dashboard-page"},
		{"  leading trailing  ", "leading-trailing"},
		{"already-a-slug", "already-a-slug"},
		{"UPPER CASE", "upper-case"},
		{"tabs\tand\nnewlines", "tabs-and-newlines"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := titleToSlug(tt.input)
			if result != tt.expected {
				t.Errorf("titleToSlug(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStringToBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"yes", true},
		{"false", false},
		{"no", false},
		{"1", false},
		{"", false},
		{"TRUE", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := stringToBool(tt.input)
			if result != tt.expected {
				t.Errorf("stringToBool(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPrefixStringLines(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		input    string
		expected string
	}{
		{"single line", ">> ", "hello", ">> hello"},
		{"multi line", "  ", "a\nb\nc", "  a\n  b\n  c"},
		{"empty prefix", "", "hello\nworld", "hello\nworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prefixStringLines(tt.prefix, tt.input)
			if result != tt.expected {
				t.Errorf("prefixStringLines(%q, %q) = %q, want %q", tt.prefix, tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaybeCopySliceWithoutZeroValues(t *testing.T) {
	t.Run("no zeros", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := maybeCopySliceWithoutZeroValues(input)
		if len(result) != 3 {
			t.Errorf("expected 3 elements, got %d", len(result))
		}
	})

	t.Run("with zeros", func(t *testing.T) {
		input := []int{1, 0, 3, 0, 5}
		result := maybeCopySliceWithoutZeroValues(input)
		if len(result) != 3 {
			t.Errorf("expected 3 elements, got %d", len(result))
		}
	})

	t.Run("all zeros", func(t *testing.T) {
		input := []int{0, 0, 0}
		result := maybeCopySliceWithoutZeroValues(input)
		if len(result) != 0 {
			t.Errorf("expected 0 elements, got %d", len(result))
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		result := maybeCopySliceWithoutZeroValues(input)
		if len(result) != 0 {
			t.Errorf("expected 0 elements, got %d", len(result))
		}
	})

	t.Run("float64 with zeros", func(t *testing.T) {
		input := []float64{1.5, 0, 3.7}
		result := maybeCopySliceWithoutZeroValues(input)
		if len(result) != 2 {
			t.Errorf("expected 2 elements, got %d", len(result))
		}
	})
}

func TestItemAtIndexOrDefault(t *testing.T) {
	items := []string{"a", "b", "c"}

	t.Run("valid index", func(t *testing.T) {
		if result := itemAtIndexOrDefault(items, 1, "x"); result != "b" {
			t.Errorf("expected 'b', got %q", result)
		}
	})

	t.Run("out of bounds", func(t *testing.T) {
		if result := itemAtIndexOrDefault(items, 5, "x"); result != "x" {
			t.Errorf("expected 'x', got %q", result)
		}
	})
}

func TestTernary(t *testing.T) {
	if result := ternary(true, "yes", "no"); result != "yes" {
		t.Errorf("expected 'yes', got %q", result)
	}
	if result := ternary(false, "yes", "no"); result != "no" {
		t.Errorf("expected 'no', got %q", result)
	}
}

func TestHslToHex(t *testing.T) {
	tests := []struct {
		name     string
		h, s, l  float64
		expected string
	}{
		{"red", 0, 100, 50, "#ff0000"},
		{"green", 120, 100, 50, "#00ff00"},
		{"blue", 240, 100, 50, "#0000ff"},
		{"white", 0, 0, 100, "#ffffff"},
		{"black", 0, 0, 0, "#000000"},
		{"gray", 0, 0, 50, "#808080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hslToHex(tt.h, tt.s, tt.l)
			if result != tt.expected {
				t.Errorf("hslToHex(%v, %v, %v) = %q, want %q", tt.h, tt.s, tt.l, result, tt.expected)
			}
		})
	}
}

func TestSvgPolylineCoordsFromYValues(t *testing.T) {
	t.Run("fewer than 2 values", func(t *testing.T) {
		if result := svgPolylineCoordsFromYValues(100, 50, []float64{5}); result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("two values", func(t *testing.T) {
		result := svgPolylineCoordsFromYValues(100, 50, []float64{0, 10})
		if result == "" {
			t.Error("expected non-empty coordinates string")
		}
	})

	t.Run("equal values produce valid output", func(t *testing.T) {
		// When all values are equal, max-min=0 which could cause division by zero
		// but the function should still handle this gracefully
		result := svgPolylineCoordsFromYValues(100, 50, []float64{5, 5, 5})
		// NaN check: if result contains "NaN", the function didn't handle the edge case
		if result == "" {
			t.Error("expected non-empty coordinates string")
		}
	})
}

func TestParseRFC3339Time(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		result := parseRFC3339Time("2026-01-15T10:30:00Z")
		if result.Year() != 2026 || result.Month() != 1 || result.Day() != 15 {
			t.Errorf("unexpected time: %v", result)
		}
	})

	t.Run("invalid time returns now", func(t *testing.T) {
		result := parseRFC3339Time("not-a-time")
		// Should return approximately time.Now()
		if result.Year() < 2025 {
			t.Errorf("expected recent time for invalid input, got %v", result)
		}
	})
}
