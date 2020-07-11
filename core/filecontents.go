package core

import "os"

// FileContents provides an interface for caching, rendering and getting cached contents of a file
type FileContents interface {
	WriteToClient(*Client, *Path) Error
	Load(*os.File, *Path) Error
	Clear()
}

// generatedFileContents is a simple FileContents implementation for holding onto a generated (virtual) file contents
type generatedFileContents struct {
	content []byte
}

// WriteToClient writes the generated file contents to the client
func (fc *generatedFileContents) WriteToClient(client *Client, path *Path) Error {
	return client.Conn().WriteBytes(fc.content)
}

// Load does nothing
func (fc *generatedFileContents) Load(fd *os.File, path *Path) Error { return nil }

// Clear does nothing
func (fc *generatedFileContents) Clear() {}

// RegularFileContents is the simplest implementation of core.FileContents for regular files
type RegularFileContents struct {
	contents []byte
}

// WriteToClient writes the current contents of FileContents to the client
func (fc *RegularFileContents) WriteToClient(client *Client, path *Path) Error {
	return client.Conn().WriteBytes(fc.contents)
}

// Load takes an open FD and loads the file contents into FileContents memory
func (fc *RegularFileContents) Load(fd *os.File, path *Path) Error {
	var err Error
	fc.contents, err = FileSystem.ReadFile(fd)
	return err
}

// Clear empties currently cached FileContents memory
func (fc *RegularFileContents) Clear() {
	fc.contents = nil
}
