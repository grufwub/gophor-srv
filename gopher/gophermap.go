package gopher

import (
	"gophor/core"
	"os"
)

const (
	InvalidGophermapErr  core.ErrorCode = -21
	SubgophermapIsDirErr core.ErrorCode = -22
	SubgophermapSizeErr  core.ErrorCode = -23
)

var (
	subgophermapSizeMax int64
)

type GophermapSection interface {
	RenderAndWrite(*core.Client) core.Error
}

func readGophermap(fd *os.File, path *core.Path) ([]GophermapSection, core.Error) {
	// Create return slice
	sections := make([]GophermapSection, 0)

	// Create hidden files map now in case later requested
	hidden := map[string]bool{
		path.Relative(): true,
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
				hidden[path.JoinRelative(line[1:])] = true
				return true

			case typeSubGophermap:
				// Parse new Path and parameters
				request := parseInternalRequest(path, line[1:])
				if returnErr != nil {
					return false
				} else if request.Path().Relative() == "" || request.Path().Relative() == path.Relative() {
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
				dirPath := core.NewPath(path.Root(), path.Dir())
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

type TextSection struct {
	contents []byte
}

func (s *TextSection) RenderAndWrite(client *core.Client) core.Error {
	return client.Conn().WriteBytes(s.contents)
}

type DirectorySection struct {
	hidden map[string]bool
	path   *core.Path
}

func (s *DirectorySection) RenderAndWrite(client *core.Client) core.Error {
	fd, err := core.FileSystem.OpenFile(s.path)
	if err != nil {
		return err
	}

	// Slice to write
	dirContents := make([]byte, 0)

	// Scan directory and build lines
	err = core.FileSystem.ScanDirectory(fd, func(file os.FileInfo) {
		filePath := core.NewPath(s.path.Root(), s.path.JoinRelative(file.Name()))

		// Skip hidden or restricted files
		_, ok := s.hidden[filePath.Relative()]
		if ok || core.IsRestrictedPath(filePath) {
			return
		}

		// Handle file, directory, ignore others
		switch {
		case file.Mode()&os.ModeDir != 0:
			// Directory -- create directory entry
			dirContents = append(dirContents, buildLine(typeDirectory, file.Name(), filePath.Selector(), core.Hostname, core.FwdPort)...)

		case file.Mode()&os.ModeType == 0:
			// File -- get item type and create entry
			t := getItemType(filePath.Relative())
			dirContents = append(dirContents, buildLine(t, file.Name(), filePath.Selector(), core.Hostname, core.FwdPort)...)
		}

		return
	})
	if err != nil {
		return err
	}

	// Write dirContents to client
	return client.Conn().WriteBytes(dirContents)
}

type FileSection struct {
	path *core.Path
}

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

type SubgophermapSection struct {
	path *core.Path
}

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

type CGISection struct {
	request core.Request
}

func (s *CGISection) RenderAndWrite(client *core.Client) core.Error {
	return core.ExecuteCGIScript(client, s.request)
}
