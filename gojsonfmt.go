// Copyright 2025 Chris Peterson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package gojsonfmt implements functions for formatting JSON data similar to
// Go source code.
package gojsonfmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// FormatJSONString formats a JSON string with Go source code like formatting
// rules.
func FormatJSONString(data string) (string, error) {
	output, err := readAndFormatJSON(strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return output, nil
}

// FormatJSONBytes decodes the byte array as a JSON string and formats it with
// Go source code like formatting rules.
func FormatJSONBytes(data []byte) ([]byte, error) {
	output, err := readAndFormatJSON(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return []byte(output), nil
}

// readAndFormatJSON reads JSON data from the reader and formats it with Go
// source code like formatting rules.
func readAndFormatJSON(r io.Reader) (string, error) {
	dec := json.NewDecoder(r)
	dec.UseNumber()

	var buf strings.Builder
	// err := formatObject(dec, &buf, 0, false)
	err := formatJSON(dec, &buf, 0)
	if err != nil && err != io.EOF {
		return "", err
	}
	output := buf.String()
	output = strings.TrimSpace(output)
	return output, nil
}

func formatJSON(dec *json.Decoder, buf *strings.Builder, indent int) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}
	switch tok := t.(type) {
	case json.Delim:
		switch tok {
		case '{':
			buf.WriteString("{\n")
			err := formatObject(dec, buf, indent+1)
			if err != nil && err != io.EOF {
				return err
			}
			buf.WriteString("}")
		default:
			return fmt.Errorf("expected '{' at start of document but got: %v", tok)
		}
	default:
		return fmt.Errorf("expected '{' at start of document but got: %v", t)
	}

	return nil
}

func formatArray(dec *json.Decoder, buf *strings.Builder, indent int) (string, error) {
	currentIndent := indent
	startedWithDelim := false
	lastItemWasObjectOrArray := true
	closing := "]"
	isFirstItem := true
	for {
		t, err := dec.Token()
		if err != nil {
			return "", err
		}
		switch tok := t.(type) {
		case json.Delim:
			switch tok {
			case '{':
				if lastItemWasObjectOrArray && !isFirstItem {
					buf.WriteString(" ")
				}
				isFirstItem = false
				startedWithDelim = true
				lastItemWasObjectOrArray = true
				// TODO: Handle one-liners
				if !dec.More() {
					t, err = dec.Token()
					tok, ok := t.(json.Delim)
					if !ok {
						return "", fmt.Errorf("expected '}' but got %T", t)
					}
					if tok != '}' {
						return "", fmt.Errorf("expected '} but got %q", tok)
					}
					buf.WriteString("{}")
					if dec.More() {
						buf.WriteString(",")
					} else {
						buf.WriteString("")
					}
					continue
				}
				buf.WriteString("{\n")
				err := formatObject(dec, buf, currentIndent+1)
				if err != nil {
					return "", err
				}
				writeIndent(buf, currentIndent)
				if dec.More() {
					buf.WriteString("},")
				} else {
					buf.WriteString("}")
				}
			case '[':
				startedWithDelim = true
				if lastItemWasObjectOrArray && !isFirstItem {
					buf.WriteString(" ")
				}
				isFirstItem = false
				lastItemWasObjectOrArray = true
				// Make sure empty JSON arrays are one-liners.
				if !dec.More() {
					t, err = dec.Token()
					tok, ok := t.(json.Delim)
					if !ok {
						return "", fmt.Errorf("expected ']' but got %T", t)
					}
					if tok != ']' {
						return "", fmt.Errorf("expected ']' but got %q", tok)
					}
					buf.WriteString("[]")
					if dec.More() {
						buf.WriteString(",")
					} else {
						buf.WriteString("")
					}
					continue
				}
				buf.WriteString("[")
				innerClose, err := formatArray(dec, buf, currentIndent)
				if err != nil {
					return "", err
				}
				if dec.More() {
					// If the last item in the array was an
					// empty list or object, the buffer
					// will still be on that line. If not,
					// then we want to align the closing
					// braces to the regular item indent.
					lastChar := last(buf)
					if lastChar != ']' && lastChar != '}' {
						writeIndent(buf, currentIndent)
					}

					buf.WriteString(innerClose)
					buf.WriteString(",")
				} else {
					closing = innerClose + "]"
				}
			case '}':
				return "", fmt.Errorf("unexpected '}' in array")
			case ']':
				return closing, nil
			}
		case json.Number, bool, string, nil:
			if !startedWithDelim {
				currentIndent += 1
				startedWithDelim = true
				buf.WriteString("\n")
				writeIndent(buf, currentIndent)
			} else if lastItemWasObjectOrArray {
				buf.WriteString("\n")
				writeIndent(buf, currentIndent)
			}
			lastItemWasObjectOrArray = false
			var out string
			switch tok := t.(type) {
			case json.Number:
				out = fmt.Sprintf("%s", tok.String())
			case bool:
				out = fmt.Sprintf("%t", tok)
			case string:
				out = fmt.Sprintf("%q", tok)
			case nil:
				out = "null"
			}
			fmt.Fprintf(buf, "%s", out)
			if dec.More() {
				buf.WriteString(",\n")
				writeIndent(buf, currentIndent)
			} else {
				buf.WriteString("\n")
			}
		default:
			return "", fmt.Errorf("cannot parse unknown token type: %T", tok)

		}
	}
}

func formatObject(dec *json.Decoder, buf *strings.Builder, indent int) error {
	nextValueIsKey := true
	for {
		t, err := dec.Token()
		if err != nil {
			return err
		}
		switch tok := t.(type) {
		case json.Delim:
			switch tok {
			case '{':
				if nextValueIsKey {
					return fmt.Errorf("expected string key but got: {")
				}
				nextValueIsKey = true
				// Make sure empty JSON objects are one-liners.
				if !dec.More() {
					t, err = dec.Token()
					tok, ok := t.(json.Delim)
					if !ok {
						return fmt.Errorf("expected '}' but got %T", t)
					}
					if tok != '}' {
						return fmt.Errorf("expected '} but got %q", tok)
					}
					buf.WriteString("{}")
					if dec.More() {
						buf.WriteString(",\n")
					} else {
						buf.WriteString("\n")
					}
					continue
				}
				buf.WriteString("{\n")
				err := formatObject(dec, buf, indent+1)
				if err != nil {
					return err
				}
				writeIndent(buf, indent)
				if dec.More() {
					buf.WriteString("},\n")
				} else {
					buf.WriteString("}\n")
				}
			case '[':
				if nextValueIsKey {
					return fmt.Errorf("expected string key but got: [")
				}
				nextValueIsKey = true
				// Make sure empty JSON arrays are one-liners.
				if !dec.More() {
					t, err = dec.Token()
					tok, ok := t.(json.Delim)
					if !ok {
						return fmt.Errorf("expected ']' but got %T", t)
					}
					if tok != ']' {
						return fmt.Errorf("expected ']' but got %q", tok)
					}
					buf.WriteString("[]")
					if dec.More() {
						buf.WriteString(",\n")
					} else {
						buf.WriteString("\n")
					}
					continue
				}
				buf.WriteString("[")
				innerClose, err := formatArray(dec, buf, indent)
				if err != nil {
					return err
				}

				// If the last item in the array was an empty
				// list or object, the buffer will still be
				// on that line. If not, then we want to
				// align the closing braces to the key with
				// an indent here.
				lastChar := last(buf)
				if lastChar != ']' && lastChar != '}' {
					writeIndent(buf, indent)
				}
				buf.WriteString(innerClose)

				if dec.More() {
					buf.WriteString(",\n")
				} else {
					buf.WriteString("\n")
				}
			case '}':
				return nil
			case ']':
				return fmt.Errorf("unexpected ']' in object")
			}
		case string, json.Number, bool, nil:
			if nextValueIsKey {
				switch tok := t.(type) {
				case string:
					nextValueIsKey = false
					writeIndent(buf, indent)
					fmt.Fprintf(buf, "%q: ", tok)
				default:
					return fmt.Errorf("expected string key but got: %T", tok)
				}
			} else {
				nextValueIsKey = true
				var out string
				switch tok := t.(type) {
				case json.Number:
					out = fmt.Sprintf("%s", tok.String())
				case bool:
					out = fmt.Sprintf("%t", tok)
				case string:
					out = fmt.Sprintf("%q", tok)
				case nil:
					out = "null"
				}

				fmt.Fprintf(buf, "%s", out)
				if dec.More() {
					buf.WriteString(",\n")
				} else {
					buf.WriteString("\n")
				}
			}
		default:
			return fmt.Errorf("cannot parse unknown token type: %T", tok)

		}
	}
}

func last(buf *strings.Builder) byte {
	s := buf.String()
	return s[len(s)-1]
}

func writeIndent(buf *strings.Builder, n int) {
	for range n {
		buf.WriteByte('\t')
	}
}
