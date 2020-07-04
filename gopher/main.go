package gopher

import (
	"flag"
	"gophor/core"
)

var (
	pageWidth int
	footer    []byte
)

func setup() {
	pWidth := flag.Uint("page-width", 80, "Gopher page width")
	footerText := flag.String("footer-text", "Gophor, a gopher server in Go!", "Footer text (empty to disable)")
	core.ParseFlagsAndSetup()

	pageWidth = int(*pWidth)
	footer = buildFooter(*footerText)
	gophermapRegex = compileGophermapRegex()
}

func Run() {
	setup()
	core.Start(serve)
}
