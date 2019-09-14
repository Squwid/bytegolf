package compiler

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestAPI(t *testing.T) {
	var exe = Execute{
		Script:       `print("hello world!")`,
		Language:     "python3",
		VersionIndex: "2",
	}
	resp, err := exe.Post()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Println(resp)
}
