package models

import "net/http"

type RemoteCompiler interface {
	Client() *http.Client
	Request() (*http.Request, error)
	ResponseChan() chan RemoteCompilerOutput
}

type RemoteCompilerOutput struct {
	Out map[string]interface{}
	Err error
}
