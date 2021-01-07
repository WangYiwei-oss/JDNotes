package model

import (
	"io/ioutil"
	"net/http"
	"text/template"
)

func WriteHTML(w http.ResponseWriter, path string) (err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	w.Write(buf)
	return nil
}

func WriteTemplate(w http.ResponseWriter, path string, data interface{}) {
	t := template.Must(template.ParseFiles(path))
	t.Execute(w, data)
}
