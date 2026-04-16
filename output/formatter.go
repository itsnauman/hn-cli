package output

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	toon "github.com/toon-format/toon-go"
)

type Format string

const (
	FormatTOON Format = "toon"
	FormatJSON Format = "json"
)

// Render marshals data in the given format and writes it to w.
func Render(w io.Writer, data any, format Format) {
	var out []byte
	var err error

	switch format {
	case FormatJSON:
		out, err = json.MarshalIndent(data, "", "  ")
	default:
		out, err = toon.Marshal(data, toon.WithLengthMarkers(true))
	}
	if err != nil {
		fmt.Fprintf(w, "marshal error: %v\n", err)
		return
	}

	w.Write(out)
	fmt.Fprintln(w)
}

// RenderError renders a structured error to the writer.
func RenderError(w io.Writer, errOut any, format Format) {
	Render(w, errOut, format)
}

// SelectFields filters a struct to include only the named fields.
// Uses the toon/json struct tag to match field names.
func SelectFields(item any, fields []string) map[string]any {
	if len(fields) == 0 {
		return nil
	}
	fieldSet := make(map[string]bool, len(fields))
	for _, f := range fields {
		fieldSet[strings.TrimSpace(f)] = true
	}

	result := make(map[string]any)
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		tag := sf.Tag.Get("json")
		if tag == "" {
			tag = sf.Tag.Get("toon")
		}
		// Strip options like ",omitempty"
		if idx := strings.Index(tag, ","); idx != -1 {
			tag = tag[:idx]
		}
		if tag == "" || tag == "-" {
			continue
		}
		if fieldSet[tag] {
			result[tag] = v.Field(i).Interface()
		}
	}
	return result
}

// SelectFieldsList applies SelectFields to a slice of items.
func SelectFieldsList[T any](items []T, fields []string) []map[string]any {
	if len(fields) == 0 {
		return nil
	}
	result := make([]map[string]any, len(items))
	for i, item := range items {
		result[i] = SelectFields(item, fields)
	}
	return result
}
