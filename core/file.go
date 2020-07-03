package core

import (
	"os"
	"sync"
	"time"
)

// FileContents provides an interface for cacheing, rendering and getting cache'd contents of a file
type FileContents interface {
	WriteToClient(*Client, *Path) Error
	Load(*os.File, *Path) Error
	Clear()
}

// GeneratedFileContents is a simple FileContents implementation for holding onto a generated (virtual) file contents
type GeneratedFileContents struct {
	content []byte
}

// WriteToClient writes the generated file contents to the client
func (fc *GeneratedFileContents) WriteToClient(client *Client, path *Path) Error {
	return client.Conn().WriteBytes(fc.content)
}

// Load does nothing
func (fc *GeneratedFileContents) Load(fd *os.File, path *Path) Error { return nil }

// Clear does nothing
func (fc *GeneratedFileContents) Clear() {}

// isGeneratedType just checks if a file's contents implemented is GeneratedFileContents
func isGeneratedType(file *File) bool {
	switch file.contents.(type) {
	case *GeneratedFileContents:
		return true
	default:
		return false
	}
}

// File provides a structure for managing a cached file including freshness, last refresh time etc
type File struct {
	contents    FileContents
	lastRefresh int64
	isFresh     bool
	sync.RWMutex
}

// NewFile returns a new File based on supplied FileContents
func NewFile(contents FileContents) *File {
	return &File{
		contents,
		0,
		true,
		sync.RWMutex{},
	}
}

// IsFresh returns files freshness status
func (f *File) IsFresh() bool {
	return f.isFresh
}

// SetFresh sets the file as fresh
func (f *File) SetFresh() {
	f.isFresh = true
}

// SetUnfresh sets the file as unfresh
func (f *File) SetUnfresh() {
	f.isFresh = false
}

// LastRefresh gets the time in nanoseconds of last refresh
func (f *File) LastRefresh() int64 {
	return f.lastRefresh
}

// UpdateRefreshTime updates the lastRefresh time to the current time in nanoseconds
func (f *File) UpdateRefreshTime() {
	f.lastRefresh = time.Now().UnixNano()
}

// CacheContents caches the file contents using the supplied file descriptor
func (f *File) CacheContents(fd *os.File, path *Path) Error {
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
func (f *File) WriteToClient(client *Client, path *Path) Error {
	return f.contents.WriteToClient(client, path)
}
