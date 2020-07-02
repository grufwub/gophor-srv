package gopher

import "strings"

const (
	maxSelectorLen = 255
	nullHost       = "null.host"
	nullPort       = "0"
	errorSelector  = "/error_selector_length"
)

// formatName is an internal function to format a gopher line name string
func formatName(name string) string {
	if len(name) > pageWidth {
		return name[:pageWidth-4] + "...\t"
	}
	return name
}

// formatSelector is an internal function to format a gopher line selector string
func formatSelector(selector string) string {
	if len(selector) > maxSelectorLen {
		return errorSelector + "\t"
	}
	return selector + "\t"
}

// formatHostPort is an internal function to format a gopher line host + port
func formatHostPort(host, port string) string {
	return host + "\t" + port
}

// buildLine is an internal function that builds a gopher line string
func buildLine(t ItemType, name, selector, host, port string) []byte {
	return []byte(string(t) + formatName(name) + formatSelector(selector) + formatHostPort(host, port) + "\r\n")
}

// buildInfoLine is an internal function that builds a gopher info line string
func buildInfoLine(line string) []byte {
	return []byte(string(typeInfo) + formatName(line) + formatHostPort(nullHost, nullPort) + "\r\n")
}

// buildErrorLine is an internal function that builds a gopher error line string
func buildErrorLine(selector string) []byte {
	return []byte(string(typeError) + selector + "\r\n.")
}

// buildFooter is an internal function that formats a raw gopher footer ready to attach to end of gophermaps (including DOS line-end)
func buildFooter(raw string) []byte {
	ret := make([]byte, 0)

	if raw != "" {
		ret = append(ret, buildInfoLine("")...)
		ret = append(ret, buildInfoLine(footerLineSeparator())...)

		for _, line := range strings.Split(raw, "\n") {
			ret = append(ret, buildInfoLine(line)...)
		}
	}

	return append(ret, []byte("\r\n")...)
}

// footerLineSeparator is an internal function that generates a footer line separator string
func footerLineSeparator() string {
	ret := ""
	for i := 0; i < pageWidth; i++ {
		ret += "_"
	}
	return ret
}
