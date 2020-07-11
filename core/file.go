package core

import (
	"os"
	"sync"
	"time"
)

// isGeneratedType just checks if a file's contents implemented is GeneratedFileContents
func isGeneratedType(f *file) bool {
	switch f.contents.(type) {
	case *generatedFileContents:
		return true
	default:
		return false
	}
}

// file provides a structure for managing a cached file including freshness, last refresh time etc
type file struct {
	contents    FileContents
	lastRefresh int64
	isFresh     bool
	sync.RWMutex
}

// newFile returns a new File based on supplied FileContents
func newFile(contents FileContents) *file {
	return &file{
		contents,
		0,
		true,
		sync.RWMutex{},
	}
}

// IsFresh returns files freshness status
func (f *file) IsFresh() bool {
	return f.isFresh
}

// SetFresh sets the file as fresh
func (f *file) SetFresh() {
	f.isFresh = true
}

// SetUnfresh sets the file as unfresh
func (f *file) SetUnfresh() {
	f.isFresh = false
}

// LastRefresh gets the time in nanoseconds of last refresh
func (f *file) LastRefresh() int64 {
	return f.lastRefresh
}

// UpdateRefreshTime updates the lastRefresh time to the current time in nanoseconds
func (f *file) UpdateRefreshTime() {
	f.lastRefresh = time.Now().UnixNano()
}

// CacheContents caches the file contents using the supplied file descriptor
func (f *file) CacheContents(fd *os.File, path *Path) Error {
	f.contents.Clear()

	// Load the file contents into cache
	err := f.contents.Load(fd, path)
	if err != nil {
		return err
	}

	// Set the cache freshness
	f.UpdateRefreshTime()
	f.SetFresh()
	return nil
}

// WriteToClient writes the cached file contents to the supplied client
func (f *file) WriteToClient(client *Client, path *Path) Error {
	return f.contents.WriteToClient(client, path)
}
