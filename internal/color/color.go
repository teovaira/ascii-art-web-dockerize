// Package color parses color specifications into ANSI 24-bit terminal escape codes.
//
// Supported formats:
//   - Named colors: black, red, green, yellow, blue, magenta, cyan, white,
//     orange, purple, pink, brown, gray (case-insensitive)
//   - Hex: #RRGGBB (e.g. #ff0000)
//   - RGB: rgb(R, G, B) (e.g. rgb(255, 0, 0))
//
// Example:
//
//	rgb, _ := Parse("red")
//	fmt.Print(color.ANSI(rgb) + "Hello" + "\033[0m")
package color

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	rgbComponents = 3
	decimalBase   = 10
	hexBase       = 16
	uint8Bits     = 8
	ansi24BitFmt  = "\033[38;2;%d;%d;%dm"
)

// RGB represents a 24-bit color.
type RGB struct {
	R, G, B uint8
}

// namedColors maps lowercase color names to their RGB values.
var namedColors = map[string]RGB{
	"black":   {0, 0, 0},
	"red":     {255, 0, 0},
	"green":   {0, 255, 0},
	"yellow":  {255, 255, 0},
	"blue":    {0, 0, 255},
	"magenta": {255, 0, 255},
	"cyan":    {0, 255, 255},
	"white":   {255, 255, 255},
	"orange":  {255, 165, 0},
	"purple":  {128, 0, 128},
	"pink":    {255, 192, 203},
	"brown":   {165, 42, 42},
	"gray":    {128, 128, 128},
}

// ErrInvalidFormat is returned when color specification is malformed.
var ErrInvalidFormat = errors.New("invalid color format")

// Parse converts a color specification string to an RGB value.
//
// Supported formats:
//   - Named (case-insensitive): red, green, blue, orange, purple, etc.
//   - Hex: #RRGGBB (e.g. #ff0000)
//   - RGB: rgb(R,G,B) (e.g. rgb(255,0,0))
//
// Parameters:
//   - colorSpec: The color specification string to parse.
//
// Returns:
//   - An RGB value representing the parsed color.
//   - ErrInvalidFormat (wrapped) if the input is empty, unknown, or malformed.
func Parse(colorSpec string) (RGB, error) {
	colorSpec = strings.TrimSpace(colorSpec)

	if colorSpec == "" {
		return RGB{}, fmt.Errorf("empty color specification: %w", ErrInvalidFormat)
	}
	lower := strings.ToLower(colorSpec)

	if color, ok := namedColors[lower]; ok {
		return color, nil
	}
	if len(colorSpec) == 7 && colorSpec[0] == '#' {
		return parseHex(colorSpec)
	}
	if strings.HasPrefix(lower, "rgb(") {
		return parseRGB(lower)
	}
	return RGB{}, fmt.Errorf("unknown color format %q: %w", colorSpec, ErrInvalidFormat)
}

// parseHex parses a hex color string in format #RRGGBB to RGB.
//
// Parameters:
//   - hex: Color string in format #RRGGBB (e.g., "#ff0000").
//
// Returns:
//   - RGB value representing the parsed color.
//   - An error if the hex format is invalid or values are out of range.
func parseHex(hex string) (RGB, error) {
	r, err := strconv.ParseUint(hex[1:3], hexBase, uint8Bits)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid red hex: %w", err)
	}

	g, err := strconv.ParseUint(hex[3:5], hexBase, uint8Bits)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid green hex: %w", err)
	}
	b, err := strconv.ParseUint(hex[5:7], hexBase, uint8Bits)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid blue hex: %w", err)
	}

	return RGB{uint8(r), uint8(g), uint8(b)}, nil
}

// parseRGB parses an RGB color string in format rgb(R,G,B) to RGB.
//
// Parameters:
//   - rgbStr: Color string in format rgb(R,G,B) (e.g., "rgb(255,0,0)").
//
// Returns:
//   - RGB value representing the parsed color.
//   - An error if the format is invalid or component values are out of range (0-255).
func parseRGB(rgbStr string) (RGB, error) {
	if !strings.HasSuffix(rgbStr, ")") {
		return RGB{}, fmt.Errorf("missing closing parenthesis: %w", ErrInvalidFormat)
	}
	content := strings.TrimPrefix(rgbStr, "rgb(")
	content = strings.TrimSuffix(content, ")")

	content = strings.TrimSpace(content)
	if content == "" {
		return RGB{}, fmt.Errorf("rgb() components cannot be empty")
	}

	parts := strings.Split(content, ",")
	if len(parts) != rgbComponents {
		return RGB{}, fmt.Errorf("rgb() requires exactly 3 components, got %d", len(parts))
	}

	var r, g, b uint8
	for i, part := range parts {
		valueString := strings.TrimSpace(part)
		value, err := strconv.ParseUint(valueString, decimalBase, uint8Bits)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid rgb() component %q: %w", valueString, err)
		}

		switch i {
		case 0:
			r = uint8(value)
		case 1:
			g = uint8(value)
		case 2:
			b = uint8(value)
		}
	}

	return RGB{r, g, b}, nil
}

// ANSI returns the 24-bit ANSI escape sequence for the given color.
//
// The returned string has the format \033[38;2;<r>;<g>;<b>m and can be
// written directly to a terminal that supports 24-bit color.
//
// Parameters:
//   - rgb: The RGB color value to convert.
//
// Returns:
//   - The ANSI escape sequence string.
func ANSI(rgb RGB) string {
	return fmt.Sprintf(ansi24BitFmt, rgb.R, rgb.G, rgb.B)
}
