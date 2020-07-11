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

// NewPath returns a new Path structure based on supplied root and relative path
func newPath(root, rel string) *Path {
	return &Path{root, rel, formatSelector(rel)}
}

// NewSanitizedPath returns a new sanitized Path structure based on root and relative path
func newSanitizedPath(root, rel string) *Path {
	return newPath(root, sanitizeRawPath(root, rel))
}

// Remap remaps a Path to a new relative path, keeping previous selector
func (p *Path) Remap(newRel string) {
	p.rel = sanitizeRawPath(p.root, newRel)
}

// Root returns file's root directory
func (p *Path) Root() string {
	return p.root
}

// Relative returns the relative path
func (p *Path) Relative() string {
	return p.rel
}

// Absolute returns the absolute path
func (p *Path) Absolute() string {
	return path.Join(p.root, p.rel)
}

// Selector returns the formatted selector path
func (p *Path) Selector() string {
	return p.sel
}

// RelativeDir returns the residing dir of the relative path
func (p *Path) RelativeDir() string {
	return path.Dir(p.rel)
}

// SelectorDir returns the residing dir of the selector path
func (p *Path) SelectorDir() string {
	return path.Dir(p.sel)
}

// Dir returns a Path object at the residing dir of the calling object (keeping separate selector intact)
func (p *Path) Dir() *Path {
	return &Path{p.root, p.RelativeDir(), p.SelectorDir()}
}

// JoinRelative returns a string appended to the current relative path
func (p *Path) JoinRelative(newRel string) string {
	return path.Join(p.rel, newRel)
}

// JoinPath appends the supplied string to the Path's relative and selector paths
func (p *Path) JoinPath(toJoin string) *Path {
	return &Path{p.root, path.Join(p.rel, toJoin), path.Join(p.sel, toJoin)}
}

// formatSelector formats a relative path to a valid selector path
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

// sanitizerUserRoot takes a generated user root directory and sanitizes it, returning a bool as to whether it's safe
func sanitizeUserRoot(root string) (string, bool) {
	root = path.Clean(root)
	if !strings.HasPrefix(root, "/home/") && strings.HasSuffix(root, "/"+userDir) {
		return "", false
	}
	return root, true
}
