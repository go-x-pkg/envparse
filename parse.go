package envparse

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	reWhitespace = regexp.MustCompile(`[\t\v\f\r ]+`)

	// ErrMustHaveTwoArguments indicates key value
	// are two arguments
	ErrMustHaveTwoArguments = errors.New("must have two arguments")

	// ErrMissingEqualsSign indicates that key=value
	// pair doesn't contains = separator
	ErrMissingEqualsSign = errors.New("syntax error - can't find =. Must be of the form: name=value")
)

// Directives is the structure used during a build run to hold the state of
// parsing directives.
type Directives struct {
	EscapeToken rune // Current escape token
}

// helper to parse words (i.e space delimited or quoted strings) in a statement.
// The quotes are preserved as part of this function and they are stripped later
// as part of processWords().
func parseWords(raw string, d *Directives) []string {
	const (
		inSpaces = iota // looking for start of a word
		inWord
		inQuote
	)

	words := []string{}
	phase := inSpaces
	word := ""
	quote := '\000'
	blankOK := false

	var (
		ch      rune
		chWidth int
	)

	for pos := 0; pos <= len(raw); pos += chWidth {
		if pos != len(raw) {
			ch, chWidth = utf8.DecodeRuneInString(raw[pos:])
		}

		if phase == inSpaces { // Looking for start of word
			if pos == len(raw) { // end of input
				break
			}

			if unicode.IsSpace(ch) { // skip spaces
				continue
			}

			phase = inWord // found it, fall through
		}

		if (phase == inWord || phase == inQuote) && (pos == len(raw)) {
			if blankOK || len(word) > 0 {
				words = append(words, word)
			}

			break
		}

		if phase == inWord {
			if unicode.IsSpace(ch) {
				phase = inSpaces

				if blankOK || len(word) > 0 {
					words = append(words, word)
				}

				word = ""
				blankOK = false

				continue
			}

			if ch == '\'' || ch == '"' {
				quote = ch
				blankOK = true
				phase = inQuote
			}

			if ch == d.EscapeToken {
				if pos+chWidth == len(raw) {
					continue // just skip an escape token at end of line
				}
				// If we're not quoted and we see an escape token, then always just
				// add the escape token plus the char to the word, even if the char
				// is a quote.
				word += string(ch)
				pos += chWidth
				ch, chWidth = utf8.DecodeRuneInString(raw[pos:])
			}

			word += string(ch)

			continue
		}

		if phase == inQuote {
			if ch == quote {
				phase = inWord
			}
			// The escape token is special except for ' quotes - can't escape anything for '
			if ch == d.EscapeToken && quote != '\'' {
				if pos+chWidth == len(raw) {
					phase = inWord
					continue // just skip the escape token at end
				}

				pos += chWidth
				word += string(ch)
				ch, chWidth = utf8.DecodeRuneInString(raw[pos:])
			}

			word += string(ch)
		}
	}

	return words
}

// ParseRawWithDirectives is wrapper around parseRaw
func ParseRawWithDirectives(raw string, d *Directives, cb func(string, string)) error {
	return parseRaw(raw, d, cb)
}

// ParseRaw is wrapper around parseRaw with default Directives
func ParseRaw(raw string, cb func(string, string)) error {
	d := Directives{EscapeToken: '\\'}
	return parseRaw(raw, &d, cb)
}

// parseRaw environment like statements. Note that this does *not* handle
// variable interpolation, which will be handled in the evaluator.
func parseRaw(raw string, d *Directives, cb func(string, string)) error {
	// This is kind of tricky because we need to support the old
	// variant:   KEY name value
	// as well as the new one:    KEY name=value ...
	// The trigger to know which one is being used will be whether we hit
	// a space or = first.  space ==> old, "=" ==> new

	preCb := func(key, value string) {
		value = strings.Trim(value, `"`)
		value = strings.Trim(value, `'`)
		value = strings.ReplaceAll(value, `\ `, " ")
		cb(key, value)
	}

	if raw == "" {
		return nil
	}

	words := parseWords(raw, d)
	if len(words) == 0 {
		return nil
	}

	// Old format (KEY name value)
	if !strings.Contains(words[0], "=") {
		parts := reWhitespace.Split(raw, 2)
		if len(parts) < 2 {
			return ErrMustHaveTwoArguments
		}

		preCb(parts[0], parts[1])

		return nil
	}

	for _, word := range words {
		if !strings.Contains(word, "=") {
			return fmt.Errorf("error parse word %q: %w", word, ErrMissingEqualsSign)
		}

		parts := strings.SplitN(word, "=", 2)
		preCb(parts[0], parts[1])
	}

	return nil
}
