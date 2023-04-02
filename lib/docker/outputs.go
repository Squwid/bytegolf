package docker

import (
	"bytes"
	"io"
	"strings"

	"github.com/docker/docker/pkg/stdcopy"
)

const maxReadBytes = 1024 * 30

// ContainerOutputs contain information about a container
// that was just started.
type ContainerOutputs struct {
	// ID is the ID of the container that was started.
	ID string

	StdOut io.Writer
	StdErr io.Writer
}

type LogWriter struct {
	// input. Once output hits max bytes for either output then
	// the reader should be closed.
	input io.ReadCloser

	buf bytes.Buffer
}

func NewLogWriter(reader io.ReadCloser) *LogWriter {
	return &LogWriter{
		input: reader,
		buf:   *bytes.NewBuffer([]byte{}),
	}
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	bw, err := lw.buf.Write(p)
	if lw.buf.Len() >= maxReadBytes {
		lw.input.Close()
	}

	return bw, err
}

// Output returns the buffer value after writing limited
// by the max number of bytes returned.
func (lw LogWriter) Output() []byte {
	if lw.buf.Len() < maxReadBytes {
		return lw.buf.Bytes()
	}
	return lw.buf.Bytes()[0:maxReadBytes]
}

// ReadLogOutputs reads log outputs from the docker container and closes the
// reader once the max number of bytes has been read.
func ReadLogOutputs(r io.ReadCloser) (*LogWriter, *LogWriter, error) {
	stdOut, stdErr := NewLogWriter(r), NewLogWriter(r)
	_, err := stdcopy.StdCopy(stdOut, stdErr, r)
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return stdOut, stdErr, nil
		}
		return nil, nil, err
	}

	return stdOut, stdErr, nil
}
