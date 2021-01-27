package main

import (
	"io/ioutil"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile("./Ajax-get.html")
	w.Write(f)
}

func f1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ajaxtest", f1)
	http.ListenAndServe(":8080", nil)
}
