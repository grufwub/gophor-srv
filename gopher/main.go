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
	admin := flag.String(adminFlagStr, "", adminDescStr)
	desc := flag.String(descFlagStr, "", descDescStr)
	geo := flag.String(geoFlagStr, "", geoDescStr)
	core.ParseFlagsAndSetup(generateErrorMessage)

	// Setup gopher specific global variables
	subgophermapSizeMax = int64(1048576.0 * *subgopherSizeMax) // convert float to megabytes
	pageWidth = int(*pWidth)
	footer = buildFooter(*footerText)
	gophermapRegex = compileGophermapRegex()

	// Generate capability files
	capsTxt := generateCapsTxt(*desc, *admin, *geo)
	robotsTxt := generateRobotsTxt()

	// Add generated files to cache
	core.FileSystem.AddGeneratedFile(core.NewPath(core.Root, "caps.txt"), capsTxt)
	core.FileSystem.AddGeneratedFile(core.NewPath(core.Root, "robots.txt"), robotsTxt)
}

// Run does as says :)
func Run() {
	setup()
	core.Start(serve)
}
