package compiler

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	if resp.StatusCode != http.StatusOK {
		ch <- models.RemoteCompilerOutput{Err: fmt.Errorf("got bad status code %v", resp.StatusCode)}
		return
	}

	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		ch <- models.RemoteCompilerOutput{Err: err}
		return
	}
	defer resp.Body.Close()

	ch <- models.RemoteCompilerOutput{Out: out}
}
