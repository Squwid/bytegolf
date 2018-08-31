package main

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestStore(t *testing.T) {
	s := Storage{
		Logs:     "true",
		Location: "aws",
	}
	c := Configuration{
		Port:    "8000",
		Storage: s,
	}

	bs, err := yaml.Marshal(c)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	t.Log("\n" + string(bs))
}
