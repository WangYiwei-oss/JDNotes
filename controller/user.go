package controller

import (
	"net/http"

	"../model"
)

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := &model.User{}
	user.Username = r.PostFormValue("username")
	user.Password = r.PostFormValue("password")
	user.SelectByNameAndPassword()
	if user.Id > 0 {
		//cookie结构
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    user.Username,
			HttpOnly: true,
			MaxAge:   3600,
		}
		//写入cookie
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/index", 302)
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

func CheckUsername(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	user.Username = r.PostFormValue("username")
	user.SelectByUsername()
	if user.Id > 0 {
		w.Write([]byte("用户名已存在"))
	} else {
		w.Write([]byte("<font style='color:green'>用户名可用</font>"))
	}
}
func CheckPassword(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	user.Password = r.PostFormValue("password")
	password2 := r.PostFormValue("password2")
	if user.Password == password2 {
		if user.CheckPassword() {
			w.Write([]byte("true"))
		} else {
			w.Write([]byte("密码强度太低，必须包含英文和数字"))
		}
	} else {
		w.Write([]byte("两次密码输入不一致"))
	}
}
func CheckEmail(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	user.Email = r.PostFormValue("email")
	if user.CheckEmailFormat() {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("邮箱格式不正确"))
	}
}
func Mypage(w http.ResponseWriter, r *http.Request) {
	model.WriteHTML(w, "views/pages/user/mypage.html")
}
func Unlogin(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("_cookie")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/index", 302)
}
