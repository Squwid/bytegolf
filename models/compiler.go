package models

import (
	"io"
	"net/http"
)

type RemoteCompiler interface {
	Client() *http.Client
	Request() (*http.Request, error)
	ResponseChan() chan RemoteCompilerOutput
}

type RemoteCompilerOutput struct {
	StatusCode int
	Body       io.ReadCloser
	Err        error
}
