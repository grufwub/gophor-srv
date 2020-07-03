package gopher

import (
	"gophor/core"
	"os"
)

// FileContents is the simplest implementation of core.FileContents for regular files
type FileContents struct {
	contents []byte
}

// WriteToClient writes the current contents of FileContents to the client
func (fc *FileContents) WriteToClient(client *core.Client, path *core.Path) core.Error {
	return client.Conn().WriteBytes(fc.contents)
}

// Load takes an open FD and loads the file contents into FileContents memory
func (fc *FileContents) Load(fd *os.File, path *core.Path) core.Error {
	var err core.Error
	fc.contents, err = core.FileSystem.ReadFile(fd)
	return err
}

// Clear empties currently cached FileContents memory
func (fc *FileContents) Clear() {
	fc.contents = nil
}

// GophermapContents .
type GophermapContents struct {
	sections []GophermapSection
}

// WriteToClient writes the current contents of FileContents to the client
func (gc *GophermapContents) WriteToClient(client *core.Client, path *core.Path) core.Error {
	for _, section := range gc.sections {
		err := section.RenderAndWrite(client)
		if err != nil {
			return err
		}
	}

	// Finally, write the footer (including last-line)
	return client.Conn().WriteBytes(footer)
}

// Load takes an open FD and loads the file contents into FileContents memory
func (gc *GophermapContents) Load(fd *os.File, path *core.Path) core.Error {
	var err core.Error
	gc.sections, err = readGophermap(fd, path)
	return err
}

// Clear empties currently cached FileContents memory
func (gc *GophermapContents) Clear() {
	gc.sections = nil
}
