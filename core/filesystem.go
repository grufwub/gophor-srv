package core

import (
	"bufio"
	"io"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	// FileReadBufSize is the file read buffer size
	fileReadBufSize int

	// MonitorSleepTime is the duration the goroutine should periodically sleep before running file cache freshness checks
	monitorSleepTime time.Duration

	// FileSizeMax is the maximum file size that is alloewd to be cached
	fileSizeMax int64

	// FileSystem is the global FileSystem object
	FileSystem *FileSystemObject

	// userDir is the set subdir name to be looked for under user's home folders
	userDir string
)

// FileSystemObject holds onto an LRUCacheMap and manages access to it, handless freshness checking and multi-threading
type FileSystemObject struct {
	cache *lruCacheMap
	sync.RWMutex
}

// NewFileSystemObject returns a new FileSystemObject
func newFileSystemObject(size int) *FileSystemObject {
	return &FileSystemObject{
		newLRUCacheMap(size),
		sync.RWMutex{},
	}
}

// StartMonitor starts the FileSystemObject freshness check monitor in its own goroutine
func (fs *FileSystemObject) StartMonitor() {
	for {
		// Sleep to not take up all the precious CPU time :)
		time.Sleep(monitorSleepTime)

		// Check file cache freshness
		fs.checkCacheFreshness()
	}
}

// checkCacheFreshness iterates through FileSystemObject's cache and check for freshness
func (fs *FileSystemObject) checkCacheFreshness() {
	// Before anything get cache lock
	fs.Lock()

	fs.cache.Iterate(func(path string, f *file) {
		// If this is a generated file we skip
		if isGeneratedType(f) {
			return
		}

		// Check file still exists on disk
		stat, err := os.Stat(path)
		if err != nil {
			SystemLog.Error("Failed to stat file in cache: %s\n", path)
			fs.cache.Remove(path)
			return
		}

		// Get last mod time and check freshness
		lastMod := stat.ModTime().UnixNano()
		if f.IsFresh() && f.LastRefresh() < lastMod {
			f.SetUnfresh()
		}
	})

	// Done! Unlock (:
	fs.Unlock()
}

// OpenFile opens a file for reading (read-only, world-readable)
func (fs *FileSystemObject) OpenFile(p *Path) (*os.File, Error) {
	fd, err := os.OpenFile(p.Absolute(), os.O_RDONLY, 0444)
	if err != nil {
		return nil, WrapError(FileOpenErr, err)
	}
	return fd, nil
}

// StatFile performs a file stat on a file at path
func (fs *FileSystemObject) StatFile(p *Path) (os.FileInfo, Error) {
	stat, err := os.Stat(p.Absolute())
	if err != nil {
		return nil, WrapError(FileStatErr, err)
	}
	return stat, nil
}

// ReadFile reads a supplied file descriptor into a return byte slice, or error
func (fs *FileSystemObject) ReadFile(fd *os.File) ([]byte, Error) {
	// Return slice
	ret := make([]byte, 0)

	// Read buffer
	buf := make([]byte, fileReadBufSize)

	// Read through file until null bytes / error
	for {
		count, err := fd.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, WrapError(FileReadErr, err)
		}

		ret = append(ret, buf[:count]...)

		if count < fileReadBufSize {
			break
		}
	}

	return ret, nil
}

// ScanFile scans a supplied file at file descriptor, using iterator function
func (fs *FileSystemObject) ScanFile(fd *os.File, iterator func(string) bool) Error {
	// Buffered reader
	rdr := bufio.NewReaderSize(fd, fileReadBufSize)

	// Iterate through file!
	for {
		// Line buffer
		b := make([]byte, 0)

		// Read until line-end, or file end!
		for {
			// Read a line
			line, isPrefix, err := rdr.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				return WrapError(FileReadErr, err)
			}

			// Append to line buffer
			b = append(b, line...)

			// If not isPrefix, we can break-out
			if !isPrefix {
				break
			}
		}

		// Run scan iterator on this line, break-out if requested
		if !iterator(string(b)) {
			break
		}
	}

	return nil
}

// ScanDirectory reads the contents of a directory and performs the iterator function on each os.FileInfo entry returned
func (fs *FileSystemObject) ScanDirectory(fd *os.File, p *Path, iterator func(os.FileInfo, *Path)) Error {
	dirList, err := fd.Readdir(-1)
	if err != nil {
		return WrapError(DirectoryReadErr, err)
	}

	// Sort by name
	sort.Sort(byName(dirList))

	// Walk through the directory list using supplied iterator function
	for _, info := range dirList {
		// Make new Path object
		fp := p.JoinPath(info.Name())

		// Skip restricted files
		if IsRestrictedPath(fp) || WithinCGIDir(fp) {
			continue
		}

		// Perform iterator
		iterator(info, p.JoinPath(info.Name()))
	}

	return nil
}

// AddGeneratedFile adds a generated file content byte slice to the file cache, with supplied path as the key
func (fs *FileSystemObject) AddGeneratedFile(p *Path, b []byte) {
	// Get write lock, defer unlock
	fs.Lock()
	defer fs.Unlock()

	// Create new generatedFileContents
	contents := &generatedFileContents{b}

	// Wrap contents in File
	file := newFile(contents)

	// Add to cache!
	fs.cache.Put(p.Absolute(), file)
}

// HandleClient handles a Client, attempting to serve their request from the filesystem whether a regular file, gophermap, dir listing or CGI script
func (fs *FileSystemObject) HandleClient(client *Client, request *Request, newFileContents func(*Path) FileContents, handleDirectory func(*FileSystemObject, *Client, *os.File, *Path) Error) Error {
	// If restricted, return error
	if IsRestrictedPath(request.Path()) {
		return NewError(RestrictedPathErr)
	}

	// Try remap request, log if so
	ok := RemapRequest(request)
	if ok {
		client.LogInfo(requestRemappedStr, request.Path().Selector(), request.Params())
	}

	// First check for file on disk
	fd, err := fs.OpenFile(request.Path())
	if err != nil {
		// Get read-lock, defer unlock
		fs.RLock()
		defer fs.RUnlock()

		// Don't throw in the towel yet! Check for generated file in cache
		file, ok := fs.cache.Get(request.Path().Absolute())
		if !ok {
			return err
		}

		// We got a generated file! Close and send as-is
		return file.WriteToClient(client, request.Path())
	}
	defer fd.Close()

	// Get stat
	stat, goErr := fd.Stat()
	if goErr != nil {
		// Unlock, return error
		fs.RUnlock()
		return WrapError(FileStatErr, goErr)
	}

	switch {
	// Directory
	case stat.Mode()&os.ModeDir != 0:
		// Don't support CGI script dir enumeration
		if WithinCGIDir(request.Path()) {
			return NewError(RestrictedPathErr)
		}

		// Else enumerate dir
		return handleDirectory(fs, client, fd, request.Path())

	// Regular file
	case stat.Mode()&os.ModeType == 0:
		// Execute script if within CGI dir
		if WithinCGIDir(request.Path()) {
			return ExecuteCGIScript(client, request)
		}

		// Else just fetch
		return fs.FetchFile(client, fd, stat, request.Path(), newFileContents)

	// Unsupported type
	default:
		return NewError(FileTypeErr)
	}
}

// FetchFile attempts to fetch a file from the cache, using the supplied file stat, Path and serving client. Returns Error status
func (fs *FileSystemObject) FetchFile(client *Client, fd *os.File, stat os.FileInfo, p *Path, newFileContents func(*Path) FileContents) Error {
	// If file too big, write direct to client
	if stat.Size() > fileSizeMax {
		return client.Conn().WriteFrom(fd)
	}

	// Get cache read lock, defer unlock
	fs.RLock()
	defer fs.RUnlock()

	// Now check for file in cache
	f, ok := fs.cache.Get(p.Absolute())
	if !ok {
		// Create new file contents with supplied function
		contents := newFileContents(p)

		// Wrap contents in file
		f = newFile(contents)

		// Cache the file contents
		err := f.CacheContents(fd, p)
		if err != nil {
			// Unlock, return error
			return err
		}

		// Get cache write lock
		fs.RUnlock()
		fs.Lock()

		// Put file in cache
		fs.cache.Put(p.Absolute(), f)

		// Switch back to cache read lock, get file read lock
		fs.Unlock()
		fs.RLock()
		f.RLock()
	} else {
		// Get file read lock
		f.RLock()

		// Check for file freshness
		if !f.IsFresh() {
			// Switch to file write lock
			f.RUnlock()
			f.Lock()

			// Refresh file contents
			err := f.CacheContents(fd, p)
			if err != nil {
				// Unlock file, return error
				f.Unlock()
				return err
			}

			// Done! Switch back to read lock
			f.Unlock()
			f.RLock()
		}
	}

	// Defer file unlock, write to client
	defer f.RUnlock()
	return f.WriteToClient(client, p)
}
