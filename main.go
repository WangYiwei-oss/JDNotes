package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//index函数用来返回主页面
func index(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadFile("views/index.html")
	if err != nil {
		panic(err)
	}
	w.Write(buf)
}

func main() {
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("views/pages/"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static/"))))
	http.Handle("/notes/", http.StripPrefix("/notes/", http.FileServer(http.Dir("views/notes/"))))
	http.HandleFunc("/index", index)
	fmt.Println("Ip:0.0.0.0\nPort:8080")
	http.ListenAndServe(":8080", nil)
}
