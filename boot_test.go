package main

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestStore(t *testing.T) {
	s := Storage{
		SaveSubmissions: true,
		SaveLogs:        true,
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

func TestGet(t *testing.T) {
	// t.Skip("this test is disabled")
	c, err := ParseConfiguration()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	fmt.Println(*c)

}
