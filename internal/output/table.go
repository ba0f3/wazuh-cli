package output

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const maxColWidth = 60

// RenderTable writes rows as a GitHub-Flavored Markdown table.
// Columns are auto-detected from the union of all keys in rows.
func RenderTable(w io.Writer, rows []map[string]interface{}) error {
	if len(rows) == 0 {
		_, err := fmt.Fprintln(w, "_No results_")
		return err
	}

	// Collect and sort column names for deterministic output
	colSet := make(map[string]struct{})
	for _, row := range rows {
		for k := range row {
			colSet[k] = struct{}{}
		}
	}
	cols := make([]string, 0, len(colSet))
	for c := range colSet {
		cols = append(cols, c)
	}
	sort.Strings(cols)

	// Compute column widths (header vs content)
	widths := make(map[string]int, len(cols))
	for _, c := range cols {
		widths[c] = len(c)
	}
	for _, row := range rows {
		for _, c := range cols {
			v := cellValue(row[c])
			if len(v) > widths[c] {
				if len(v) > maxColWidth {
					widths[c] = maxColWidth
				} else {
					widths[c] = len(v)
				}
			}
		}
	}

	// Header row
	headerParts := make([]string, len(cols))
	for i, c := range cols {
		headerParts[i] = pad(c, widths[c])
	}
	if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(headerParts, " | ")); err != nil {
		return err
	}

	// Separator row
	sepParts := make([]string, len(cols))
	for i, c := range cols {
		sepParts[i] = strings.Repeat("-", widths[c])
	}
	if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(sepParts, " | ")); err != nil {
		return err
	}

	// Data rows
	for _, row := range rows {
		cellParts := make([]string, len(cols))
		for i, c := range cols {
			v := cellValue(row[c])
			if len(v) > maxColWidth {
				v = v[:maxColWidth-3] + "..."
			}
			cellParts[i] = pad(v, widths[c])
		}
		if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(cellParts, " | ")); err != nil {
			return err
		}
	}

	return nil
}

// cellValue converts any interface{} to a printable string for a table cell.
func cellValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		// Integers come as float64 in JSON
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case []interface{}:
		parts := make([]string, len(val))
		for i, item := range val {
			parts[i] = cellValue(item)
		}
		return strings.Join(parts, ", ")
	case map[string]interface{}:
		// Flatten nested objects to key:value pairs
		parts := make([]string, 0, len(val))
		for k, v2 := range val {
			parts = append(parts, k+":"+cellValue(v2))
		}
		sort.Strings(parts)
		return "{" + strings.Join(parts, ", ") + "}"
	default:
		return fmt.Sprintf("%v", val)
	}
}

// pad right-pads s to the given width with spaces.
func pad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
