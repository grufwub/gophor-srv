package gopher

import (
	"gophor/core"
	"strings"
)

func serve(client *core.Client) {
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
	request := NewSanitizedRequest(core.Root, path, params)

	// Check for remap
	request = core.RemapRequest(request)

	// Handle the request!
	err = core.FileSystem.Fetch(
		request.Path(),
		func(path *core.Path) core.FileContents {
			// Return FileContents or GophermapContents depending on file name
			if isGophermap(path) {
				return &GophermapContents{}
			}
			return &FileContents{}
		},
		func(file *core.File) core.Error {
			//
			return nil
		},
	)
}

func handleError(client *core.Client, err core.Error) {

}
