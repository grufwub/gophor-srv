package gopher

var (
	pageWidth int
	footer    []byte
)

func configure(width uint, footerText string) {
	pageWidth = int(width)
	footer = buildFooter(footerText)
	gophermapRegex = compileGophermapRegex()
}
