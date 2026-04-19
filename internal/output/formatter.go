// Package output handles formatting and writing wazuh-cli responses
// in JSON, Markdown, or raw formats.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Format constants.
const (
	FormatJSON     = "json"
	FormatMarkdown = "markdown"
	FormatRaw      = "raw"
)

// Formatter writes structured data to an output stream.
type Formatter struct {
	format    string
	pretty    bool
	writer    io.Writer
	errWriter io.Writer
}

// New creates a Formatter for the given format.
func New(format string, pretty bool) *Formatter {
	return &Formatter{
		format:    format,
		pretty:    pretty,
		writer:    os.Stdout,
		errWriter: os.Stderr,
	}
}

// Write outputs data to stdout in the configured format.
// data should be json.RawMessage or any JSON-serializable value.
func (f *Formatter) Write(data interface{}) error {
	switch f.format {
	case FormatRaw:
		return f.writeRaw(data)
	case FormatMarkdown:
		return f.writeMarkdown(data)
	default:
		return f.writeJSON(data)
	}
}

// WriteError outputs a machine-readable error to stdout (JSON) or stderr (other formats).
func (f *Formatter) WriteError(code int, message, detail string) {
	errObj := map[string]interface{}{
		"error":   true,
		"code":    code,
		"message": message,
	}
	if detail != "" {
		errObj["detail"] = detail
	}

	switch f.format {
	case FormatJSON:
		enc := json.NewEncoder(f.writer)
		enc.SetIndent("", "  ")
		_ = enc.Encode(errObj)
	default:
		fmt.Fprintf(f.errWriter, "Error %d: %s\n", code, message)
		if detail != "" {
			fmt.Fprintf(f.errWriter, "Detail: %s\n", detail)
		}
	}
}

// writeJSON encodes data as JSON to stdout.
func (f *Formatter) writeJSON(data interface{}) error {
	enc := json.NewEncoder(f.writer)
	if f.pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(data)
}

// writeRaw writes the raw bytes (or JSON) directly.
func (f *Formatter) writeRaw(data interface{}) error {
	switch v := data.(type) {
	case []byte:
		_, err := f.writer.Write(v)
		if err != nil {
			return err
		}
		// Ensure newline at end
		if len(v) > 0 && v[len(v)-1] != '\n' {
			_, err = fmt.Fprintln(f.writer)
		}
		return err
	case json.RawMessage:
		_, err := f.writer.Write(v)
		if err != nil {
			return err
		}
		if len(v) > 0 && v[len(v)-1] != '\n' {
			_, err = fmt.Fprintln(f.writer)
		}
		return err
	default:
		return f.writeJSON(data)
	}
}

// writeMarkdown renders the data as a Markdown table.
func (f *Formatter) writeMarkdown(data interface{}) error {
	// Convert to []map[string]interface{} for table rendering
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Try as array first
	var rows []map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &rows); err == nil && len(rows) > 0 {
		return RenderTable(f.writer, rows)
	}

	// Try as single object
	var row map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &row); err == nil {
		return RenderTable(f.writer, []map[string]interface{}{row})
	}

	// Fall back to JSON
	return f.writeJSON(data)
}

// ExitCode translates an error to a wazuh-cli exit code.
//
//	0 = success
//	1 = client error (bad flags, network, config)
//	2 = API error
//	3 = auth error
//	4 = permission denied
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	switch err.(type) {
	case interface{ IsAuth() bool }:
		return 3
	case interface{ IsPermission() bool }:
		return 4
	case interface{ IsAPI() bool }:
		return 2
	default:
		return 1
	}
}

// PrintInfo writes an informational message to stderr (respects --quiet via quiet flag).
func PrintInfo(quiet bool, format string, args ...interface{}) {
	if !quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
