package gopher

import (
	"flag"
	"gophor/core"
)

// setup parses gopher specific flags, and all core flags, preparing server for .Run()
func setup() {
	pWidth := flag.Uint(PageWidthFlagStr, 80, PageWidthDescStr)
	footerText := flag.String(FooterTextFlagStr, "Gophor, a gopher server in Go!", FooterTextDescStr)
	subgopherSizeMax := flag.Float64(SubgopherSizeMaxFlagStr, 1.0, SubgopherSizeMaxDescStr)
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
