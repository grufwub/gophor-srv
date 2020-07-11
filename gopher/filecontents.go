package gopher

import (
	"gophor/core"
	"os"
)

// gophermapContents is an implementation of core.FileContents that holds individually renderable sections of a gophermap
type gophermapContents struct {
	sections []gophermapSection
}

// WriteToClient renders each cached section of the gophermap, and writes them to the client
func (gc *gophermapContents) WriteToClient(client *core.Client, path *core.Path) core.Error {
	for _, section := range gc.sections {
		err := section.RenderAndWrite(client)
		if err != nil {
			return err
		}
	}

	// Finally, write the footer (including last-line)
	return client.Conn().WriteBytes(footer)
}

// Load takes an open FD and loads the gophermap contents into memory as different renderable sections
func (gc *gophermapContents) Load(fd *os.File, path *core.Path) core.Error {
	var err core.Error
	gc.sections, err = readGophermap(fd, path)
	return err
}

// Clear empties currently cached GophermapContents memory
func (gc *gophermapContents) Clear() {
	gc.sections = nil
}
