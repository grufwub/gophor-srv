package gopher

import (
	"gophor/core"
	"os"
	"strings"
)

// serve is the global gopher server's serve function
func serve(client *core.Client) {
	// Receive line from client
	received, err := client.Conn().ReadLine()
	if err != nil {
		client.LogError(clientReadFailStr)
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
		client.LogInfo(clientRedirectFmtStr, line)
		return
	}

	// Parse new request
	request, err := core.ParseURLEncodedRequest(line)
	if err != nil {
		client.LogError(clientRequestParseFailStr)
		handleError(client, err)
		return
	}

	// Handle the request!
	err = core.FileSystem.HandleClient(
		client,
		request,
		newFileContents,
		func(fs *core.FileSystemObject, client *core.Client, fd *os.File, p *core.Path) core.Error {
			// First check for gophermap, create gophermap Path object
			gophermap := p.JoinPath("gophermap")

			// If gophermap exists, we fetch this
			fd2, err := fs.OpenFile(gophermap)
			if err == nil {
				stat, osErr := fd2.Stat()
				if osErr == nil {
					return fs.FetchFile(client, fd2, stat, gophermap, newFileContents)
				}

				// Else, just close fd2
				fd2.Close()
			}

			// Slice to write
			dirContents := make([]byte, 0)

			// Add directory heading + empty line
			dirContents = append(dirContents, buildLine(typeInfo, "[ "+core.Hostname+p.Selector()+" ]", "TITLE", nullHost, nullPort)...)
			dirContents = append(dirContents, buildInfoLine("")...)

			// Scan directory and build lines
			err = fs.ScanDirectory(
				fd,
				p,
				func(file os.FileInfo, fp *core.Path) {
					// Append new formatted file listing (if correct type)
					dirContents = appendFileListing(dirContents, file, fp)
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
		client.LogError(clientServeFailStr, request.Path().Absolute())
	} else {
		client.LogInfo(clientServedStr, request.Path().Absolute())
	}
}

// handleError determines whether to send an error response to the client, and logs to system
func handleError(client *core.Client, err core.Error) {
	response, ok := generateErrorResponse(err.Code())
	if ok {
		client.Conn().WriteBytes(response)
	}
	core.SystemLog.Error(err.Error())
}

// newFileContents returns a new FileContents object
func newFileContents(p *core.Path) core.FileContents {
	if isGophermap(p) {
		return &GophermapContents{}
	}
	return &core.RegularFileContents{}
}
