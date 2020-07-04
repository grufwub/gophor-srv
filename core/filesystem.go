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
)

// FileSystemObject holds onto an LRUCacheMap and manages access to it, handless freshness checking and multi-threading
type FileSystemObject struct {
	cache *LRUCacheMap
	sync.RWMutex
}

// NewFileSystemObject returns a new FileSystemObject
func NewFileSystemObject(size int) *FileSystemObject {
	return &FileSystemObject{
		NewLRUCacheMap(size),
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

	fs.cache.Iterate(func(path string, file *File) {
		// If this is a generated file we skip
		if isGeneratedType(file) {
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
		if file.IsFresh() && file.LastRefresh() < lastMod {
			file.SetUnfresh()
		}
	})

	// Done! Unlock (:
	fs.Unlock()
}

// OpenFile opens a file for reading (read-only, world-readable)
func (fs *FileSystemObject) OpenFile(path *Path) (*os.File, Error) {
	fd, err := os.OpenFile(path.Absolute(), os.O_RDONLY, 0444)
	if err != nil {
		return nil, WrapError(FileOpenErr, err)
	}
	return fd, nil
}

// StatFile performs a file stat on a file at path
func (fs *FileSystemObject) StatFile(path *Path) (os.FileInfo, Error) {
	stat, err := os.Stat(path.Absolute())
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
	reader := bufio.NewReaderSize(fd, fileReadBufSize)

	// Iterate through file!
	for {
		// Line buffer
		b := make([]byte, 0)

		// Read until line-end, or file end!
		for {
			// Read a line
			line, isPrefix, err := reader.ReadLine()
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
func (fs *FileSystemObject) ScanDirectory(fd *os.File, iterator func(os.FileInfo)) Error {
	dirList, err := fd.Readdir(-1)
	if err != nil {
		return WrapError(DirectoryReadErr, err)
	}

	// Sort by name
	sort.Sort(byName(dirList))

	// Walk through the directory list using supplied iterator function
	for _, info := range dirList {
		iterator(info)
	}

	return nil
}

// HandleClient .
func (fs *FileSystemObject) HandleClient(client *Client, request Request, newFileContents func(*Path) FileContents, handleDirectory func(*FileSystemObject, *Client, *os.File, *Path) Error) Error {
	// If restricted, return error
	if IsRestrictedPath(request.Path()) {
		return NewError(RestrictedPathErr)
	}

	// Try remap request
	request = RemapRequest(request)

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
	stat, statErr := fd.Stat()
	if err != nil {
		// Unlock, return error
		fs.RUnlock()
		return WrapError(FileStatErr, statErr)
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
		return fs.fetchFile(client, fd, request.Path(), newFileContents)

	// Unsupported type
	default:
		return NewError(FileTypeErr)
	}
}

func (fs *FileSystemObject) fetchFile(client *Client, fd *os.File, path *Path, newFileContents func(*Path) FileContents) Error {
	// Get cache read lock, defer unlock
	fs.RLock()
	defer fs.RUnlock()

	// Now check for file in cache
	file, ok := fs.cache.Get(path.Absolute())
	if !ok {
		// Create new file contents with supplied function
		contents := newFileContents(path)

		// Wrap contents in file
		file = NewFile(contents)

		// Cache the file contents
		err := file.CacheContents(fd, path)
		if err != nil {
			// Unlock, return error
			return err
		}

		// Get cache write lock
		fs.RUnlock()
		fs.Lock()

		// Put file in cache
		fs.cache.Put(path.Absolute(), file)

		// Switch back to cache read lock, get file read lock
		fs.Unlock()
		fs.RLock()
		file.RLock()
	} else {
		// Get file read lock
		file.RLock()

		// Check for file freshness
		if !file.IsFresh() {
			// Switch to file write lock
			file.RUnlock()
			file.Lock()

			// Refresh file contents
			err := file.CacheContents(fd, path)
			if err != nil {
				// Unlock file, return error
				file.Unlock()
				return err
			}

			// Done! Switch back to read lock
			file.Unlock()
			file.RLock()
		}
	}

	// Defer file unlock, write to client
	defer file.RUnlock()
	return file.WriteToClient(client, path)
}
