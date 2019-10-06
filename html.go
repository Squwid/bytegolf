// This file controls all of the view templates in the container for now until
// getting a javascript front end to do more extensive testing

package main

import "net/http"

func indexHTML(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", struct {
		MOTD string
	}{
		MOTD: "Bugs not birds...",
	})
}
