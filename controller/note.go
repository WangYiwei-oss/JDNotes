package controller

import (
	"net/http"

	"../model"
)

func MyNotes(w http.ResponseWriter, r *http.Request) {
	//处理导航栏
	firstcs := model.GetAllFirstclassWithCookie(r)
	arg := make(map[string]interface{})
	arg["login"] = model.CheckLogin(r)
	arg["firstclass"] = firstcs
	//处理导航栏跳转
	firstclass := r.FormValue("firstclass")
	ns := model.SelectByFirstClass(firstclass) //当前一级标题下所有笔记信息
	arg["titles"] = ns
	//处理笔记内容
	title := r.FormValue("title")
	buf := model.OpenByTitle(title)
	arg["text"] = string(buf)
	model.WriteTemplate(w, "views/pages/mynotes.html", arg)
}
