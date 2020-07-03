package gopher

import (
	"gophor/core"
	"os"
	"strings"
)

func serve(client *core.Client) {
	// First defer client close
	defer client.Conn().Close()

	// Receive line from client
	received, err := client.Conn().ReadLine()
	if err != nil {
		handleError(client, err)
		return
	}

	// Convert to string
	line := string(received)

	// If prefixed by 'URL:' send a redirect
	lenBefore := len(line)
	line = strings.TrimPrefix(line, "URL:")
	if len(line) < lenBefore {
		client.Conn().WriteBytes(generateHTMLRedirect(line))
		core.AccessLog.Info("Redirecting to: %s", line)
		return
	}

	// Parse supplied URL
	path, params, err := core.ParseSafeURL(line)
	if err != nil {
		handleError(client, err)
		return
	}

	// Create new request from path and params
	request := newSanitizedRequest(core.Root, path, params)

	// Check for remap
	request, _ = core.RemapRequest(request).(*Request)

	// Handle the request!
	err = core.FileSystem.HandleClient(
		client,
		request,
		func(path *core.Path) core.FileContents {
			if isGophermap(path) {
				return &GophermapContents{}
			}
			return &FileContents{}
		},
		func(fs *core.FileSystemObject, client *core.Client, fd *os.File, path *core.Path) core.Error {
			// Slice to write
			dirContents := make([]byte, 0)

			// Add directory heading + empty line
			dirContents = append(dirContents, buildLine(typeInfo, "[ "+core.Hostname+path.Selector()+" ]", "TITLE", nullHost, nullPort)...)
			dirContents = append(dirContents, buildInfoLine("")...)

			// Scan directory and build lines
			err := fs.ScanDirectory(
				fd,
				func(file os.FileInfo) {
					filePath := core.NewPath(path.Root(), path.JoinRelative(file.Name()))

					// Skip restricted files
					if core.IsRestrictedPath(filePath) {
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
				},
			)
			if err != nil {
				return err
			}

			// Add footer, write contents
			dirContents = append(dirContents, footer...)
			return client.Conn().WriteBytes(dirContents)
		},
	)

	// Final error handling
	if err != nil {
		handleError(client, err)
	}
}

func handleError(client *core.Client, err core.Error) {

}
