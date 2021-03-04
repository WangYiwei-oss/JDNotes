package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"./controller"
	"./model"
)

//index函数用来返回主页面
func index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./views/index.html"))
	t.Execute(w, model.CheckLogin(r))
}

func ScanNotes() {
	for {
		startTime := time.Now()
		model.ScanAllNotesToMysql("./views/pages/notes")
		fmt.Printf("%s:更新笔记完成,用时[%s]\n", startTime.Format("2006-01-02 15:04:05"), time.Since(startTime))
		time.Sleep(time.Minute * 1)
	}
}

func main() {
	go ScanNotes()
	//静态文件系统
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("views/pages/"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static/"))))
	http.Handle("/notes/", http.StripPrefix("/notes/", http.FileServer(http.Dir("views/notes/"))))
	//
	http.HandleFunc("/index", index)                          //主页
	http.HandleFunc("/tologin", controller.ToLogin)           //登陆页面
	http.HandleFunc("/login", controller.Login)               //登陆
	http.HandleFunc("/toregister", controller.ToRegister)     //注册页面
	http.HandleFunc("/register", controller.Register)         //注册
	http.HandleFunc("/mypage", controller.Mypage)             //我的主页
	http.HandleFunc("/unlogin", controller.Unlogin)           //退出登陆
	http.HandleFunc("/expassistant", controller.Expassistant) //实验助手
	//Ajax用的请求
	http.HandleFunc("/checkUsername", controller.CheckUsername)
	http.HandleFunc("/checkPassword", controller.CheckPassword)
	http.HandleFunc("/checkEmail", controller.CheckEmail)

	http.HandleFunc("/cal_iqy_input_post", controller.CalIqyInput)
	http.HandleFunc("/cal_integral", controller.CalIntegral)
	//
	http.HandleFunc("/mynotes", controller.MyNotes) //我的笔记页面
	fmt.Println("Ip:0.0.0.0\nPort:8080")
	http.ListenAndServe(":8080", nil)

}
