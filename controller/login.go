package controller

import (
	"net/http"

	"../model"
)

func Login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	user.Username = r.PostFormValue("username")
	user.Password = r.PostFormValue("password")
	user.SelectByNameAndPassword()
	if user.Id > 0 {
		model.WriteTemplate(w, "views/pages/user/login_success.html", user.Username)
	} else {
		model.WriteTemplate(w, "views/pages/user/login.html", "用户名或密码错误")
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	user.Username = r.PostFormValue("username")
	if user.Username == "" {
		model.WriteTemplate(w, "views/pages/user/register.html", "请输入用户名")
		return
	}
	user.SelectByUsername()
	if user.Id > 0 {
		model.WriteTemplate(w, "views/pages/user/register.html", "该用户已存在")
		return
	}
	user.Password = r.PostFormValue("password")
	password2 := r.PostFormValue("password2")
	if user.Password == "" || password2 == "" {
		model.WriteTemplate(w, "views/pages/user/register.html", "请输入密码")
		return
	}
	if user.Password != password2 {
		model.WriteTemplate(w, "views/pages/user/register.html", "两次密码输入不一致")
		return
	}
	user.Email = r.PostFormValue("email")
	b := user.CheckEmailFormat()
	if b == false {
		model.WriteTemplate(w, "views/pages/user/register.html", "邮箱格式错误")
		return
	}
	user.Number = r.PostFormValue("number")
	user.InsertIntoSql()
	model.WriteTemplate(w, "views/pages/user/register_success.html", user.Username)
}
