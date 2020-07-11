package gopher

import (
	"flag"
	"gophor/core"
)

// setup parses gopher specific flags, and all core flags, preparing server for .Run()
func setup() {
	pWidth := flag.Uint(pageWidthFlagStr, 80, pageWidthDescStr)
	footerText := flag.String(footerTextFlagStr, "Gophor, a gopher server in Go!", footerTextDescStr)
	subgopherSizeMax := flag.Float64(subgopherSizeMaxFlagStr, 1.0, subgopherSizeMaxDescStr)
	core.ParseFlagsAndSetup(generateErrorMessage)

	subgophermapSizeMax = int64(1048576.0 * *subgopherSizeMax) // convert float to megabytes
	pageWidth = int(*pWidth)
	footer = buildFooter(*footerText)
	gophermapRegex = compileGophermapRegex()
}

// Run does as says :)
func Run() {
	setup()
	core.Start(serve)
}
