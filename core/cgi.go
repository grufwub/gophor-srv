package core

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	// cgiEnv holds the global slice of constant CGI environment variables
	cgiEnv []string

	// maxCGIRunTime specifies the maximum time a CGI script can run for
	maxCGIRunTime time.Duration

	// httpPrefixBufSize
	httpPrefixBufSize int

	// ExecuteCGIScript is a pointer to the currently set CGI execution function
	ExecuteCGIScript func(client *Client, request Request) Error
)

func setupInitialCGIEnv(safePath string) []string {
	env := make([]string, 0)

	SystemLog.Info("CGI safe path: %s", safePath)
	env = append(env, "PATH="+safePath)
	env = append(env, "SERVER_NAME="+Hostname)
	env = append(env, "SERVER_PORT="+FwdPort)
	env = append(env, "DOCUMENT_ROOT="+Root)

	return env
}

func generateCGIEnv(client *Client, request Request) []string {
	env := cgiEnv

	env = append(env, "REMOTE_ADDR="+client.IP())
	env = append(env, "QUERY_STRING="+request.Params())
	env = append(env, "SCRIPT_NAME="+request.Path().Relative())
	env = append(env, "SCRIPT_FILENAME="+request.Path().Absolute())
	env = append(env, "SELECTOR="+request.Path().Selector())
	env = append(env, "REQUEST_URI="+request.Path().Selector())

	return env
}

func executeCGIScriptNoHTTP(client *Client, request Request) Error {
	return execute(client.Conn().Writer(), request.Path(), generateCGIEnv(client, request))
}

func executeCGIScriptStripHTTP(client *Client, request Request) Error {
	// Create new httpStripWriter
	httpWriter := newhttpStripWriter(client.Conn().Writer())

	// Begin executing script
	err := execute(httpWriter, request.Path(), generateCGIEnv(client, request))

	// Parse HTTP headers (if present). Return error or continue letting output of script -> client
	cgiStatusErr := httpWriter.FinishUp()
	if cgiStatusErr != nil {
		return cgiStatusErr
	}
	return err
}

func execute(writer io.Writer, path *Path, env []string) Error {
	// Create cmd object
	cmd := exec.Command(path.Absolute())

	// Set new process group id
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Setup cmd environment
	cmd.Env, cmd.Dir = env, path.Root()

	// Setup cmd out writer
	cmd.Stdout = writer

	// Start executing
	err := cmd.Start()
	if err != nil {
		return WrapError(CGIStartErr, err)
	}

	// Setup goroutine to kill cmd after maxCGIRunTime
	go func() {
		// At least let the script try to finish...
		time.Sleep(maxCGIRunTime)

		// We've already finished
		if cmd.ProcessState != nil {
			return
		}

		// Get process group id
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err != nil {
			SystemLog.Fatal("Process unfinished, PGID not found!")
		}

		// Kill process group!
		err = syscall.Kill(-pgid, syscall.SIGTERM)
		if err != nil {
			SystemLog.Fatal("Error stopping process group %d: %s", pgid, err.Error())
		}
	}()

	// Wait for command to finish, get exit code
	err = cmd.Wait()
	exitCode := 0
	if err != nil {
		// Error, try to get exit code
		exitError, ok := err.(*exec.ExitError)
		if ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			exitCode = waitStatus.ExitStatus()
		} else {
			exitCode = 1
		}
	} else {
		// No error! Get exit code directly from command process state
		waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = waitStatus.ExitStatus()
	}

	// Non-zero exit code? Return error
	if exitCode != 0 {
		SystemLog.Error("Exit executing: %s [%d]", path.Absolute(), exitCode)
		return NewError(CGIExitCodeErr)
	}

	// Exit fine!
	return nil
}

type httpStripWriter struct {
	/* Wrapper to io.Writer that reads a predetermined amount into a buffer
	 * then parses the buffer for valid HTTP headers and status code, deciding
	 * whether to strip these headers or returning with an HTTP status code.
	 */
	writer     io.Writer
	skipBuffer []byte
	skipIndex  int
	err        Error

	/* We set underlying write function with a variable, so that each call
	 * to .Write() doesn't have to perform a check every time whether we need
	 * to keep checking for headers to skip.
	 */
	WriteFunc func(*httpStripWriter, []byte) (int, error)
}

func newhttpStripWriter(writer io.Writer) *httpStripWriter {
	return &httpStripWriter{
		writer,
		make([]byte, httpPrefixBufSize),
		0,
		nil,
		writeCheckForHeaders,
	}
}

// addToSkipBuffer adds supplied bytes to the skip buffer, returning number added
func (w *httpStripWriter) addToSkipBuffer(data []byte) int {
	/* Figure out how much data we need to add */
	toAdd := len(w.skipBuffer) - w.skipIndex
	if len(data) < toAdd {
		toAdd = len(data)
	}

	/* Add the data to the skip buffer! */
	copy(w.skipBuffer[w.skipIndex:], data[:toAdd])
	w.skipIndex += toAdd
	return toAdd
}

// parseHTTPHeaderSection checks if we've received a valid HTTP header section, and determine if we should continue writing
func (w *httpStripWriter) parseHTTPHeaderSection() (bool, bool) {
	validHeaderSection, shouldContinue := false, true
	for _, header := range strings.Split(string(w.skipBuffer), "\r\n") {
		header = strings.ToLower(header)

		// Try look for status header
		lenBefore := len(header)
		header = strings.TrimPrefix(header, "status:")
		if len(header) < lenBefore {
			// Ensure no spaces + just number
			header = strings.Split(header, " ")[0]

			// Ignore 200
			if header == "200" {
				continue
			}

			// Any other value indicates error, should not continue
			shouldContinue = false

			// Parse error code
			code := CGIStatusUnknownErr
			switch header {
			case "400":
				code = CGIStatus400Err
			case "401":
				code = CGIStatus401Err
			case "403":
				code = CGIStatus403Err
			case "404":
				code = CGIStatus404Err
			case "408":
				code = CGIStatus408Err
			case "410":
				code = CGIStatus410Err
			case "500":
				code = CGIStatus500Err
			case "501":
				code = CGIStatus501Err
			case "503":
				code = CGIStatus503Err
			}

			// Set error code
			w.err = NewError(code)
			continue
		}

		// Found a content-type header, this is a valid header section
		if strings.Contains(header, "content-type:") {
			validHeaderSection = true
		}
	}

	return validHeaderSection, shouldContinue
}

// writeSkipBuffer writes contents of skipBuffer to the underlying writer if necessary
func (w *httpStripWriter) writeSkipBuffer() (bool, error) {
	// Defer resetting skipIndex
	defer func() {
		w.skipIndex = 0
	}()

	// First try parse the headers, determine next steps
	validHeaders, shouldContinue := w.parseHTTPHeaderSection()

	// Valid headers received, don't bother writing. Return the shouldContinue value
	if validHeaders {
		return shouldContinue, nil
	}

	// Default is to write skip buffer contents, shouldContinue only means something with valid headers
	_, err := w.writer.Write(w.skipBuffer[:w.skipIndex])
	return true, err
}

func (w *httpStripWriter) FinishUp() Error {
	/* If SkipBuffer still has contents, in case of data written being less
	 * than w.Size() --> check this data for HTTP headers to strip, parse
	 * any status codes and write this content with underlying writer if
	 * necessary.
	 */
	if w.skipIndex > 0 {
		w.writeSkipBuffer()
	}

	/* Return HttpStripWriter error code if set */
	return w.err
}

func (w *httpStripWriter) Write(data []byte) (int, error) {
	/* Write using whatever write function is currently set */
	return w.WriteFunc(w, data)
}

func writeRegular(w *httpStripWriter, data []byte) (int, error) {
	/* Regular write function */
	return w.writer.Write(data)
}

func writeCheckForHeaders(w *httpStripWriter, data []byte) (int, error) {
	split := bytes.Split(data, []byte("\r\n\r\n"))
	if len(split) == 1 {
		/* Try add these to skip buffer */
		added := w.addToSkipBuffer(data)

		if added < len(data) {
			defer func() {
				/* Having written skipbuffer after this if clause, set write to regular */
				w.WriteFunc = writeRegular
			}()

			doContinue, err := w.writeSkipBuffer()
			if !doContinue {
				return len(data), io.EOF
			} else if err != nil {
				return added, err
			}

			/* Write remaining data not added to skip buffer */
			count, err := w.writer.Write(data[added:])
			if err != nil {
				return added + count, err
			}
		}

		return len(data), nil
	} else {
		defer func() {
			/* No use for skip buffer after this clause, set write to regular */
			w.WriteFunc = writeRegular
			w.skipIndex = 0
		}()

		/* Try add what we can to skip buffer */
		added := w.addToSkipBuffer(append(split[0], []byte("\r\n\r\n")...))

		/* Write skip buffer data if necessary, check if we should continue */
		doContinue, err := w.writeSkipBuffer()
		if !doContinue {
			return len(data), io.EOF
		} else if err != nil {
			return added, err
		}

		/* Write remaining data not added to skip buffer */
		count, err := w.writer.Write(data[added:])
		if err != nil {
			return added + count, err
		}

		return len(data), nil
	}
}
