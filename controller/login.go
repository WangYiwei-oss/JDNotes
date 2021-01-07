package controller

import (
	"net/http"

	"../model"
	"../utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	user.Username = r.PostFormValue("username")
	user.Password = r.PostFormValue("password")
	sql := "SELECT id,name,password,number,email,imgpath FROM users WHERE name=? AND password=?"
	row := utils.Db.QueryRow(sql, user.Username, user.Password)
	row.Scan(&user.Id, &user.Username, &user.Password, &user.Number, &user.Email, &user.Img_path)
	if user.Id > 0 {
		model.WriteTemplate(w, "views/pages/user/login_success.html", user.Username)
	} else {
		model.WriteTemplate(w, "views/pages/user/login.html", "用户名或密码错误")
	}
}
