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

	// Generate request
	request, err := parseRequest(line)
	if err != nil {
		handleError(client, err)
		return
	}

	// Check for remap
	request = core.RemapRequest(
		request,
		func(remap *core.PathRemap, request core.Request) core.Request {
			//
			return nil
		},
	)

	// Handle the request!
	err = core.FileSystem.Fetch(
		request.Path(),
		func(path *core.Path) core.FileContents {
			//
			return nil
		},
		func(file *core.File) core.Error {
			//
			return nil
		},
	)
}

func handleError(client *core.Client, err core.Error) {

}
