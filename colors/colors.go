// This package is inspired by the chalk's ansi-styles and supports-color
// (JavaScript).
//
// See: https://en.wikipedia.org/wiki/ANSI_escape_code .
// See: https://misc.flogisoft.com/bash/tip_colors_and_formatting .
// See: https://www.npmjs.com/package/ansi-styles .
// See: https://www.npmjs.com/package/supports-color .
// See: https://github.com/termstandard/colors .

package colors

import (
	"os"
	"strconv"

	"github.com/mattn/go-isatty"
)

var (
	supportsColor     bool
	supportsANSI256   bool
	supportsTrueColor bool
)

func init() {
	// Check is TTY.
	isTTY := isatty.IsTerminal(os.Stdout.Fd()) ||
		isatty.IsCygwinTerminal(os.Stdout.Fd())

	if !isTTY {
		return
	}

	// Terminals.
	term := os.Getenv("TERM")

	if term == "dumb" {
		return
	}

	if colorTerm, ok := os.LookupEnv("COLORTERM"); ok {
		supportsColor = true

		if colorTerm == "truecolor" {
			supportsANSI256 = true
			supportsTrueColor = true
		}
		return
	}

	if termProg, ok := os.LookupEnv("TERM_PROGRAM"); ok {
		if termProg == "Apple_Terminal" {
			supportsColor = true
			supportsANSI256 = true
			return
		}

		// TODO(SuperPaintman): add iTerm.app
	}

	// TODO(SuperPaintman): /-256(color)?$/i

	// TODO(SuperPaintman): /^screen|^xterm|^vt100|^vt220|^rxvt|color|ansi|cygwin|linux/i

	// TODO(SuperPaintman): add win32 checker.

	// CI.
	if _, ok := os.LookupEnv("CI"); ok {
		cis := [...]string{
			"TRAVIS",
			"CIRCLECI",
			"APPVEYOR",
			"GITLAB_CI",
			"GITHUB_ACTIONS",
			"BUILDKITE",
			"DRONE",
		}

		for _, name := range cis {
			if _, ok := os.LookupEnv(name); ok {
				supportsColor = true
				return
			}
		}

		if os.Getenv("CI_NAME") == "codeship" {
			supportsColor = true
			return
		}
	}

	// TODO(SuperPaintman): add TeamCity checker.
}

func SupportsColor() bool     { return supportsColor }
func SupportsANSI256() bool   { return supportsANSI256 }
func SupportsTrueColor() bool { return supportsTrueColor }

type Mode uint8

const (
	Auto Mode = iota
	Never
	Always
)

var Colors Mode = Auto

var (
	ForceANSI256   bool
	ForceTrueColor bool
)

func shouldUseColors(clrs Mode) bool {
	return clrs == Always ||
		(clrs == Auto && supportsColor)
}

func shouldUseANSI256(clrs Mode, force bool) bool {
	return (supportsANSI256 || force) &&
		shouldUseColors(clrs)
}

func shouldUseTrueColor(clrs Mode, force bool) bool {
	return (supportsTrueColor || force) &&
		shouldUseColors(clrs)
}

type Attribute uint8

func (a Attribute) String() string {
	if !shouldUseColors(Colors) {
		return ""
	}

	return attributeToString(uint8(a))
}

//go:generate python ./generate_reset_attributes.py

func (s Attribute) Reset() Attribute {
	return resetAttributes[s]
}

const (
	Reset              Attribute = 0
	ResetBold          Attribute = 22 // 21 isn't widely supported and 22 does the same thing.
	ResetDim           Attribute = 22
	ResetItalic        Attribute = 23
	ResetUnderline     Attribute = 24
	ResetInverse       Attribute = 27
	ResetHidden        Attribute = 28
	ResetStrikethrough Attribute = 29
	ResetOverline      Attribute = 55

	Bold          Attribute = 1
	Dim           Attribute = 2
	Italic        Attribute = 3
	Underline     Attribute = 4
	Inverse       Attribute = 7
	Hidden        Attribute = 8
	Strikethrough Attribute = 9
	Overline      Attribute = 53
)

const (
	ResetColor Attribute = 39

	Black   Attribute = 30
	Red     Attribute = 31
	Green   Attribute = 32
	Yellow  Attribute = 33
	Blue    Attribute = 34
	Magenta Attribute = 35
	Cyan    Attribute = 36
	White   Attribute = 37

	BlackBright   Attribute = 90
	RedBright     Attribute = 91
	GreenBright   Attribute = 92
	YellowBright  Attribute = 93
	BlueBright    Attribute = 94
	MagentaBright Attribute = 95
	CyanBright    Attribute = 96
	WhiteBright   Attribute = 97

	// Aliases.
	Gray Attribute = BlackBright
)

const (
	ResetBgColor Attribute = 49

	BgBlack   Attribute = 40
	BgRed     Attribute = 41
	BgGreen   Attribute = 42
	BgYellow  Attribute = 43
	BgBlue    Attribute = 44
	BgMagenta Attribute = 45
	BgCyan    Attribute = 46
	BgWhite   Attribute = 47

	BgBlackBright   Attribute = 100
	BgRedBright     Attribute = 101
	BgGreenBright   Attribute = 102
	BgYellowBright  Attribute = 103
	BgBlueBright    Attribute = 104
	BgMagentaBright Attribute = 105
	BgCyanBright    Attribute = 106
	BgWhiteBright   Attribute = 107

	// Aliases.
	BgGray Attribute = BgBlackBright
)

//go:generate python ./generate_ansi_attribute_string.py

// attributeToString converts Attribute to string.
//
// A hack for fast and inlinable Attribute to string converion (like what the
// stringer does).
//
// Unfortunately Go can't inline strconv.Itoa (at least Go 1.16) and we can't
// write a regular function. So we need some evil slice hacks :(.
//
// If you look in the git log you will find a faster version of this function
// but it has fewer ways to be inlined.
//
// NOTE(SuperPaintman): I would like to remove it in the future.
func attributeToString(i uint8) string {
	return ansiAttributeString[ansiAttributeIndex[uint16(i)]:ansiAttributeIndex[uint16(i)+1]]
}

// TODO(SuperPaintman): make it inlinable.
func ANSI256(color uint8) string {
	if !shouldUseANSI256(Colors, ForceANSI256) {
		return ""
	}

	// TODO(SuperPaintman): optimize it with a preassembled slice.
	return "\x1b[38;5;" + strconv.Itoa(int(color)) + "m"
}

// TODO(SuperPaintman): make it inlinable.
func BgANSI256(color uint8) string {
	if !shouldUseANSI256(Colors, ForceANSI256) {
		return ""
	}

	// TODO(SuperPaintman): optimize it with a preassembled slice.
	return "\x1b[48;5;" + strconv.Itoa(int(color)) + "m"
}

// 24-bit or truecolor or ANSI 16 millions.
func TrueColor(r, g, b uint8) string {
	if !shouldUseTrueColor(Colors, ForceTrueColor) {
		return ""
	}

	// TODO(SuperPaintman): optimize it with a preassembled slice.
	return "\x1b[38;2;" +
		strconv.Itoa(int(r)) + ";" +
		strconv.Itoa(int(g)) + ";" +
		strconv.Itoa(int(b)) + "m"
}

func BgTrueColor(r, g, b uint8) string {
	if !shouldUseTrueColor(Colors, ForceTrueColor) {
		return ""
	}

	// TODO(SuperPaintman): optimize it with a preassembled slice for uint8.
	return "\x1b[48;2;" +
		strconv.Itoa(int(r)) + ";" +
		strconv.Itoa(int(g)) + ";" +
		strconv.Itoa(int(b)) + "m"
}

type RGB struct {
	R, G, B uint8
}

func TrueColorRGB(color RGB) string {
	return TrueColor(color.R, color.G, color.B)
}

func BgTrueColorRGB(color RGB) string {
	return BgTrueColor(color.R, color.G, color.B)
}
