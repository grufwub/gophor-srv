package gopher

import (
	"gophor/core"
	"os"
)

var (
	// subgophermapSizeMax specifies the maximum size of an included subgophermap
	subgophermapSizeMax int64
)

// GophermapSection is an interface that specifies individually renderable (and writeable) sections of a gophermap
type GophermapSection interface {
	RenderAndWrite(*core.Client) core.Error
}

// readGophermap reads a FD and Path as gophermap sections
func readGophermap(fd *os.File, p *core.Path) ([]GophermapSection, core.Error) {
	// Create return slice
	sections := make([]GophermapSection, 0)

	// Create hidden files map now in case later requested
	hidden := map[string]bool{
		p.Relative(): true,
	}

	// Error setting within nested function below
	var returnErr core.Error

	// Perform scan of gophermap FD
	titleAlready := false
	scanErr := core.FileSystem.ScanFile(
		fd,
		func(line string) bool {
			// Parse the line item type and handle
			lineType := parseLineType(line)
			switch lineType {
			case typeInfoNotStated:
				// Append TypeInfo to beginning of line
				sections = append(sections, &TextSection{buildInfoLine(line)})
				return true

			case typeTitle:
				// Reformat title line to send as info line with appropriate selector
				if !titleAlready {
					sections = append(sections, &TextSection{buildLine(typeInfo, line[1:], "TITLE", nullHost, nullPort)})
					titleAlready = true
					return true
				}
				returnErr = core.NewError(InvalidGophermapErr)
				return false

			case typeComment:
				// ignore this line
				return true

			case typeHiddenFile:
				// Add to hidden files map
				hidden[p.JoinRelative(line[1:])] = true
				return true

			case typeSubGophermap:
				// Parse new Path and parameters
				request := core.ParseInternalRequest(p, line[1:])
				if returnErr != nil {
					return false
				} else if request.Path().Relative() == "" || request.Path().Relative() == p.Relative() {
					returnErr = core.NewError(InvalidGophermapErr)
					return false
				}

				// Open FD
				var subFD *os.File
				subFD, returnErr = core.FileSystem.OpenFile(request.Path())
				if returnErr != nil {
					return false
				}

				// Get stat
				stat, err := subFD.Stat()
				if err != nil {
					returnErr = core.WrapError(core.FileStatErr, err)
					return false
				} else if stat.IsDir() {
					returnErr = core.NewError(SubgophermapIsDirErr)
					return false
				}

				// Handle CGI script
				if core.WithinCGIDir(request.Path()) {
					sections = append(sections, &CGISection{request})
					return true
				}

				// Error out if file too big
				if stat.Size() > subgophermapSizeMax {
					returnErr = core.NewError(SubgophermapSizeErr)
					return false
				}

				// Handle regular file
				if !isGophermap(request.Path()) {
					sections = append(sections, &FileSection{})
					return true
				}

				// Handle gophermap
				sections = append(sections, &SubgophermapSection{})
				return true

			case typeEnd:
				// Last line, break-out!
				return false

			case typeEndBeginList:
				// Append DirectorySection object then break, as-with typeEnd
				dirPath := p.Dir()
				sections = append(sections, &DirectorySection{hidden, dirPath})
				return false

			default:
				// Default is appending to sections slice as TextSection
				sections = append(sections, &TextSection{[]byte(line + "\r\n")})
				return true
			}
		},
	)

	// Check the scan didn't exit with error
	if returnErr != nil {
		return nil, returnErr
	} else if scanErr != nil {
		return nil, scanErr
	}

	return sections, nil
}

// TextSection is a simple implementation that holds line's byte contents as-is
type TextSection struct {
	contents []byte
}

// RenderAndWrite simply writes the byte slice to the client
func (s *TextSection) RenderAndWrite(client *core.Client) core.Error {
	return client.Conn().WriteBytes(s.contents)
}

// DirectorySection is an implementation that holds a dir path, and map of hidden files, to later list a dir contents
type DirectorySection struct {
	hidden map[string]bool
	path   *core.Path
}

// RenderAndWrite scans and renders a list of the contents of a directory (skipping hidden or restricted files)
func (s *DirectorySection) RenderAndWrite(client *core.Client) core.Error {
	fd, err := core.FileSystem.OpenFile(s.path)
	if err != nil {
		return err
	}

	// Slice to write
	dirContents := make([]byte, 0)

	// Scan directory and build lines
	err = core.FileSystem.ScanDirectory(fd, func(file os.FileInfo) {
		p := s.path.JoinPath(file.Name())

		// Skip hidden or restricted files
		_, ok := s.hidden[p.Relative()]
		if ok || core.IsRestrictedPath(p) || core.WithinCGIDir(p) {
			return
		}

		// Append new formatted file listing (if correct type)
		dirContents = appendFileListing(dirContents, file, p)
	})
	if err != nil {
		return err
	}

	// Write dirContents to client
	return client.Conn().WriteBytes(dirContents)
}

// FileSection is an implementation that holds a file path, and writes the file contents to client
type FileSection struct {
	path *core.Path
}

// RenderAndWrite simply opens, reads and writes the file contents to the client
func (s *FileSection) RenderAndWrite(client *core.Client) core.Error {
	// Open FD for the file
	fd, err := core.FileSystem.OpenFile(s.path)
	if err != nil {
		return err
	}

	// Read the file contents into memory
	b, err := core.FileSystem.ReadFile(fd)
	if err != nil {
		return err
	}

	// Write the file contents to the client
	return client.Conn().WriteBytes(b)
}

// SubgophermapSection is an implementation to hold onto a gophermap path, then read, render and write contents to a client
type SubgophermapSection struct {
	path *core.Path
}

// RenderAndWrite reads, renders and writes the contents of the gophermap to the client
func (s *SubgophermapSection) RenderAndWrite(client *core.Client) core.Error {
	// Get FD for gophermap
	fd, err := core.FileSystem.OpenFile(s.path)
	if err != nil {
		return err
	}

	// Read gophermap into sections
	sections, err := readGophermap(fd, s.path)
	if err != nil {
		return err
	}

	// Write each of the sections (AAAA COULD BE RECURSIONNNNN)
	for _, section := range sections {
		err := section.RenderAndWrite(client)
		if err != nil {
			return err
		}
	}

	return nil
}

// CGISection is an implementation that holds onto a built request, then processing as a CGI request on request
type CGISection struct {
	request *core.Request
}

// RenderAndWrite takes the request, and executes the associated CGI script with parameters
func (s *CGISection) RenderAndWrite(client *core.Client) core.Error {
	return core.ExecuteCGIScript(client, s.request)
}
