package core

import "path"

// Specific Path error codes
const (
	RestrictedPathErr ErrorCode = -17
)

// Path safely holds a file path
type Path struct {
	root string
	rel  string
}

// Root returns file's root directory
func (p *Path) Root() string {
	return p.root
}

// Relative returns the file's relative path
func (p *Path) Relative() string {
	return p.rel
}

// Absolute returns the file's absolute path
func (p *Path) Absolute() string {
	return path.Join(p.root, p.rel)
}

// Selector returns the file's selector path
func (p *Path) Selector() string {
	return formatSelector(p.rel)
}

// formatSelector formats a relative path to a selector path
func formatSelector(relPath string) string {
	switch len(relPath) {
	case 0:
		return "/"
	case 1:
		if relPath[0] == '.' {
			return "/"
		}
		return "/" + relPath
	default:
		if relPath[0] == '/' {
			return relPath
		}
		return "/" + relPath
	}
}
