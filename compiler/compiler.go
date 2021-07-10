package compiler

import (
	"github.com/Squwid/bytegolf/models"
)

func Compile(rc models.RemoteCompiler) {
	ch := rc.ResponseChan()

	req, err := rc.Request()
	if err != nil {
		ch <- models.RemoteCompilerOutput{Err: err}
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := rc.Client().Do(req)
	if err != nil {
		ch <- models.RemoteCompilerOutput{Err: err}
		return
	}

	ch <- models.RemoteCompilerOutput{
		StatusCode: resp.StatusCode,
		Body:       resp.Body,
	}
}
