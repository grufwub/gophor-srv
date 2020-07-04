package core

import (
	"path"
	"strings"
)

// Path safely holds a file path
type Path struct {
	root string // root dir
	rel  string // relative path
	sel  string // selector path
}

// NewPath returns a new Path structure
func NewPath(root, rel string) *Path {
	return &Path{root, rel, formatSelector(rel)}
}

// NewSanitizedPath returns a new sanitized Path structure
func NewSanitizedPath(root, rel string) *Path {
	return NewPath(root, sanitizeRawPath(root, rel))
}

// Remap remaps a Path to a new relative path, keeping previous selector
func (p *Path) Remap(newRel string) *Path {
	newPath := NewPath(p.root, sanitizeRawPath(p.root, newRel))
	newPath.sel = p.sel
	return newPath
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

func (p *Path) Dir() string {
	dir := p.rel
	last := len(dir) - 1

	for last > 0 {
		dir = dir[:last-1]
		if dir[last] == '/' {
			break
		}
		last--
	}

	return dir
}

// JoinRelative .
func (p *Path) JoinRelative(newRel string) string {
	return path.Join(p.rel, newRel)
}

// formatSelector formats a relative path to a selector path
func formatSelector(rel string) string {
	switch len(rel) {
	case 0:
		return "/"
	case 1:
		if rel[0] == '.' {
			return "/"
		}
		return "/" + rel
	default:
		if rel[0] == '/' {
			return rel
		}
		return "/" + rel
	}
}

// sanitizeRawPath takes a root and relative path, and returns a sanitized relative path
func sanitizeRawPath(root, rel string) string {
	// Start by cleaning
	rel = path.Clean(rel)

	if path.IsAbs(rel) {
		// Absolute path, try trimming root and leading '/'
		rel = strings.TrimPrefix(strings.TrimPrefix(rel, root), "/")
	} else {
		// Relative path, if back dir traversal give them server root
		if strings.HasPrefix(rel, "..") {
			rel = ""
		}
	}

	return rel
}
