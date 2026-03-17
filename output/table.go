package output

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/alphabatem/flux_cli/dto"
)

func printTable(response *dto.CLIResponse) {
	if !response.Success {
		fmt.Fprintf(os.Stderr, "Error [%s]: %s\n", response.Error.Code, response.Error.Message)
		return
	}

	if response.Data == nil {
		fmt.Println("No data")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Try to render as a slice of maps or structs
	switch data := response.Data.(type) {
	case []interface{}:
		printSliceTable(w, data)
	case map[string]interface{}:
		printMapTable(w, data)
	default:
		// For complex types, marshal to map first
		b, err := json.Marshal(data)
		if err != nil {
			fmt.Fprintf(w, "%v\n", data)
			return
		}
		var m interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			fmt.Fprintf(w, "%s\n", string(b))
			return
		}
		switch v := m.(type) {
		case []interface{}:
			printSliceTable(w, v)
		case map[string]interface{}:
			printMapTable(w, v)
		default:
			fmt.Fprintf(w, "%v\n", data)
		}
	}
}

func printMapTable(w *tabwriter.Writer, m map[string]interface{}) {
	fmt.Fprintf(w, "KEY\tVALUE\n")
	fmt.Fprintf(w, "---\t-----\n")
	for k, v := range m {
		fmt.Fprintf(w, "%s\t%v\n", k, formatValue(v))
	}
}

func printSliceTable(w *tabwriter.Writer, items []interface{}) {
	if len(items) == 0 {
		fmt.Println("No data")
		return
	}

	// Extract headers from first item
	first, ok := items[0].(map[string]interface{})
	if !ok {
		for _, item := range items {
			fmt.Fprintf(w, "%v\n", item)
		}
		return
	}

	var headers []string
	for k := range first {
		headers = append(headers, k)
	}

	fmt.Fprintf(w, "%s\n", strings.Join(headers, "\t"))
	fmt.Fprintf(w, "%s\n", strings.Repeat("---\t", len(headers)))

	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		var vals []string
		for _, h := range headers {
			vals = append(vals, fmt.Sprintf("%v", formatValue(m[h])))
		}
		fmt.Fprintf(w, "%s\n", strings.Join(vals, "\t"))
	}
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Map {
		b, _ := json.Marshal(v)
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}
