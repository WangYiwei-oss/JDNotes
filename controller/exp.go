package controller

import (
	"fmt"
	"net/http"
	"strings"

	"../model"
)

func Expassistant(w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("type")
	arg := make(map[string]interface{})
	content := model.RenderContentByType(t)
	arg["content"] = content
	script := model.RenderScriptByType(t)
	arg["script"] = script
	model.WriteTemplate(w, "views/pages/expassistant.html", arg)
}

func CalIqyInput(w http.ResponseWriter, r *http.Request) {
	b1 := r.PostFormValue("b1")
	b2 := r.PostFormValue("b2")
	fmt.Println(b1, b2)
	samples := r.PostFormValue("samples")
	sss := strings.Split(samples, "\n")
	res := "样品 量子效率 蓝光吸收率"
	for _, ss := range sss {
		s := strings.Split(ss, " ")
		res += ("<br/>" + s[0] + " " + s[1] + " " + s[2])
	}
	w.Write([]byte(res))
}
