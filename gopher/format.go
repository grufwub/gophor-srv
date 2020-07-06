package gopher

import (
	"gophor/core"
	"os"
	"strings"
)

// Gophermap line formatting constants
const (
	maxSelectorLen = 255
	nullHost       = "null.host"
	nullPort       = "0"
	errorSelector  = "/error_selector_length"
)

var (
	// pageWidth is the maximum set page width of a gophermap document to render to
	pageWidth int

	// footer holds the formatted footer text (if supplied), and gophermap last-line
	footer []byte
)

// formatName formats a gopher line name string
func formatName(name string) string {
	if len(name) > pageWidth {
		return name[:pageWidth-4] + "...\t"
	}
	return name + "\t"
}

// formatSelector formats a gopher line selector string
func formatSelector(selector string) string {
	if len(selector) > maxSelectorLen {
		return errorSelector + "\t"
	}
	return selector + "\t"
}

// formatHostPort formats a gopher line host + port
func formatHostPort(host, port string) string {
	return host + "\t" + port
}

// buildLine builds a gopher line string
func buildLine(t ItemType, name, selector, host, port string) []byte {
	return []byte(string(t) + formatName(name) + formatSelector(selector) + formatHostPort(host, port) + "\r\n")
}

// buildInfoLine builds a gopher info line string
func buildInfoLine(line string) []byte {
	return []byte(string(typeInfo) + formatName(line) + formatHostPort(nullHost, nullPort) + "\r\n")
}

// buildErrorLine builds a gopher error line string
func buildErrorLine(selector string) []byte {
	return []byte(string(typeError) + selector + "\r\n" + ".\r\n")
}

// appendFileListing formats and appends a new file entry as part of a directory listing
func appendFileListing(b []byte, file os.FileInfo, p *core.Path) []byte {
	switch {
	case file.Mode()&os.ModeDir != 0:
		return append(b, buildLine(typeDirectory, file.Name(), p.Selector(), core.Hostname, core.FwdPort)...)
	case file.Mode()&os.ModeType == 0:
		t := getItemType(p.Relative())
		return append(b, buildLine(t, file.Name(), p.Selector(), core.Hostname, core.FwdPort)...)
	default:
		return b
	}
}

// buildFooter formats a raw gopher footer ready to attach to end of gophermaps (including DOS line-end)
func buildFooter(raw string) []byte {
	ret := make([]byte, 0)

	if raw != "" {
		ret = append(ret, buildInfoLine(footerLineSeparator())...)

		for _, line := range strings.Split(raw, "\n") {
			ret = append(ret, buildInfoLine(line)...)
		}
	}

	return append(ret, []byte(".\r\n")...)
}

// footerLineSeparator is an internal function that generates a footer line separator string
func footerLineSeparator() string {
	ret := ""
	for i := 0; i < pageWidth; i++ {
		ret += "_"
	}
	return ret
}
