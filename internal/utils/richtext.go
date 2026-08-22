package utils

import "strings"

type FontStyle string

const (
	FontBoldItalicSerif FontStyle = "bold_italic_serif"
	FontBoldSansSerif   FontStyle = "bold_sans"
	FontMonospace       FontStyle = "monospace"
	FontScript          FontStyle = "script"
	FontDoubleStruck    FontStyle = "double_struck"
	FontSmallCaps       FontStyle = "small_caps"
)

var fontMaps = map[FontStyle]map[rune]rune{}

func init() {
	build := func(style FontStyle, upperStart, lowerStart, digitStart rune, hasDigits bool) {
		m := make(map[rune]rune)
		for i := rune(0); i < 26; i++ {
			m['A'+i] = upperStart + i
			m['a'+i] = lowerStart + i
		}
		if hasDigits {
			for i := rune(0); i < 10; i++ {
				m['0'+i] = digitStart + i
			}
		}
		fontMaps[style] = m
	}

	build(FontBoldSansSerif, 0x1D5D4, 0x1D5EE, 0x1D7EC, true)
	build(FontMonospace, 0x1D670, 0x1D68A, 0x1D7F6, true)
	build(FontScript, 0x1D4D0, 0x1D4EA, 0, false)
	build(FontDoubleStruck, 0x1D538, 0x1D552, 0x1D7D8, true)

	biMap := make(map[rune]rune)
	for i := rune(0); i < 26; i++ {
		biMap['A'+i] = 0x1D468 + i
		biMap['a'+i] = 0x1D482 + i
	}
	fontMaps[FontBoldItalicSerif] = biMap

	scSrc := "abcdefghijklmnopqrstuvwxyz"
	scDst := []rune{'ᴀ', 'ʙ', 'ᴄ', 'ᴅ', 'ᴇ', 'ꜰ', 'ɢ', 'ʜ', 'ɪ', 'ᴊ', 'ᴋ', 'ʟ', 'ᴍ', 'ɴ', 'ᴏ', 'ᴘ', 'ǫ', 'ʀ', 'ꜱ', 'ᴛ', 'ᴜ', 'ᴠ', 'ᴡ', 'x', 'ʏ', 'ᴢ'}
	scMap := make(map[rune]rune)
	for i, c := range scSrc {
		scMap[c] = scDst[i]
		scMap[rune(strings.ToUpper(string(c))[0])] = scDst[i]
	}
	fontMaps[FontSmallCaps] = scMap
}

func StyleText(input string, style FontStyle) string {
	m, ok := fontMaps[style]
	if !ok {
		return input
	}
	var b strings.Builder
	for _, r := range input {
		if styled, found := m[r]; found {
			b.WriteRune(styled)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func Bold(s string) string      { return "<b>" + s + "</b>" }
func Italic(s string) string    { return "<i>" + s + "</i>" }
func Underline(s string) string { return "<u>" + s + "</u>" }
func Strike(s string) string    { return "<s>" + s + "</s>" }
func Spoiler(s string) string   { return "<tg-spoiler>" + s + "</tg-spoiler>" }
func Code(s string) string      { return "<code>" + s + "</code>" }
func Mention(name string, userID int64) string {
	return `<a href="tg://user?id=` + itoa(userID) + `">` + name + "</a>"
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}
