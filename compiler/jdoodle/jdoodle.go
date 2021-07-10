package jdoodle

import "os"

var clientID, clientSecret string

func init() {
	clientID = os.Getenv("JDOODLE_CLIENT")
	clientSecret = os.Getenv("JDOODLE_SECRET")

	if clientID == "" || clientSecret == "" {
		panic("JDOODLE_CLIENT and JDOODLE_SECRET cannot be blank")
	}
}
