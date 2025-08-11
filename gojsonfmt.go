package gojsonfmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func FormatJSONString(data string) (string, error) {
	dataAsBytes, err := io.ReadAll(strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return FormatJSONStream(dataAsBytes)
}

func FormatJSONStream(data []byte) (string, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	var buf strings.Builder
	err := formatObject(dec, &buf, 0, false, false)
	if err != nil && err != io.EOF {
		return "", err
	}
	output := buf.String()
	output = strings.TrimSpace(output)
	return output, nil
}

func formatObject(dec *json.Decoder, buf *strings.Builder, indent int, inArray, inObject bool) error {
	nextTokenIsValue := false
	valueNeedsNewline := inArray
	for {
		t, err := dec.Token()
		if err != nil {
			return err
		}
		switch tok := t.(type) {
		case json.Delim:
			switch tok {
			case '{':
				nextTokenIsValue = false
				buf.WriteString("{")
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
				} else {
					buf.WriteString("\n")
					if err := formatObject(dec, buf, indent+1, false, true); err != nil {
						return err
					}
					writeIndent(buf, indent)
				}
				if inArray {
					if dec.More() {
						buf.WriteString("}, ")
					} else {
						buf.WriteString("}")
					}
				} else {
					if dec.More() {
						buf.WriteString("},\n")
					} else {
						buf.WriteString("}\n")
					}
				}
			case '[':
				nextTokenIsValue = false
				buf.WriteString("[")

				// Make sure empty JSON lists are one-liners.
				if !dec.More() {
					t, err = dec.Token()
					tok, ok := t.(json.Delim)
					if !ok {
						return fmt.Errorf("expected ']' but got %T", t)
					}
					if tok != ']' {
						return fmt.Errorf("expected '] but got %q", tok)
					}
					buf.WriteString("]")
					if inArray {
						if dec.More() {
							buf.WriteString(", ")
						}
						// If not more, don't write.
					} else if inObject {
						if dec.More() {
							buf.WriteString(",\n")
						} else {
							buf.WriteString("\n")
						}
					}
					continue
				}

				if err := formatObject(dec, buf, indent, true, false); err != nil {
					return err
				}
				s := buf.String()
				lastChar := s[len(s)-1]
				if lastChar != '}' && lastChar != ']' {
					writeIndent(buf, indent)
				}
				if inArray {
					if dec.More() {
						buf.WriteString("], ")
					} else {
						buf.WriteString("]")
					}
				} else {
					if dec.More() {
						buf.WriteString("],\n")
					} else {
						buf.WriteString("]\n")
					}
				}
			case '}':
				return nil
			case ']':
				return nil
			}
		case string:
			if valueNeedsNewline {
				buf.WriteString("\n")
				valueNeedsNewline = false
			}
			if inObject {
				if nextTokenIsValue {
					fmt.Fprintf(buf, "%q", tok)
					if dec.More() {
						buf.WriteString(",\n")
					} else {
						buf.WriteString("\n")
					}
					nextTokenIsValue = false
				} else {
					writeIndent(buf, indent)
					fmt.Fprintf(buf, "%q: ", tok)
					nextTokenIsValue = true
				}
			} else {
				writeIndent(buf, indent+1)
				fmt.Fprintf(buf, "%q", tok)
				if dec.More() {
					buf.WriteString(",\n")
				} else {
					buf.WriteString("\n")
				}
			}
		default: // These are always value types.
			if valueNeedsNewline {
				buf.WriteString("\n")
				valueNeedsNewline = false
			}
			if inArray {
				writeIndent(buf, indent+1)
			}
			switch tok := t.(type) {
			case json.Number:
				fmt.Fprintf(buf, "%s", tok.String())
			case bool:
				fmt.Fprintf(buf, "%t", tok)
			case nil:
				buf.WriteString("null")
			}
			if dec.More() {
				buf.WriteString(",\n")
			} else {
				buf.WriteString("\n")
			}
			nextTokenIsValue = false
		}
	}
}

func writeIndent(buf *strings.Builder, n int) {
	for range n {
		buf.WriteByte('\t')
	}
}
