package main

import (
	"./controller"
	"./model"
	"./utils"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

//index函数用来返回主页面
func index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./views/index.html"))
	args := make(map[string]interface{})
	args["login"] = model.CheckLogin(r)
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		username := cookie.Value
		args["username"] = username
	} else {
		args["username"] = ""
	}
	t.Execute(w, args)
}

func ScanNotes() {
	for {
		startTime := time.Now()
		utils.Db.Exec("DELETE FROM notes")
		model.ScanAllNotesToMysql("./views/pages/notes")
		model.ScanAllNotesToMysql("./views/pages/usernote")
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
	http.HandleFunc("/index", index)                            //主页
	http.HandleFunc("/tologin", controller.ToLogin)             //登陆页面
	http.HandleFunc("/login", controller.Login)                 //登陆
	http.HandleFunc("/toregister", controller.ToRegister)       //注册页面
	http.HandleFunc("/register", controller.Register)           //注册
	http.HandleFunc("/mypage", controller.Mypage)               //我的主页
	http.HandleFunc("/unlogin", controller.Unlogin)             //退出登陆
	http.HandleFunc("/expassistant", controller.Expassistant)   //实验助手
	http.HandleFunc("/addnote", controller.AddNote)             //写笔记页面
	http.HandleFunc("/writenote", controller.WriteNote)         //提交笔记
	http.HandleFunc("/mypage_mynotes", controller.UserNote)     //我的笔记页面
	http.HandleFunc("/deletenote", controller.DeleteNote)       //删除我的笔记页面中的笔记
	http.HandleFunc("/alternote", controller.AlertNote)         //修改笔记页面
	http.HandleFunc("/alterusernote", controller.AlertUsernote) //提交笔记修改
	http.HandleFunc("/addexp", controller.AddExp)
	http.HandleFunc("/adduserexp", controller.AddUserExp)
	http.HandleFunc("/mypage_myexps", controller.UserExp)
	http.HandleFunc("/deleteexp", controller.DeleteExp)
	//Ajax用的请求
	http.HandleFunc("/checkUsername", controller.CheckUsername)
	http.HandleFunc("/checkPassword", controller.CheckPassword)
	http.HandleFunc("/checkEmail", controller.CheckEmail)

	http.HandleFunc("/cal_iqy_input_post", controller.CalIqyInput)
	http.HandleFunc("/cal_integral", controller.CalIntegral)
	http.HandleFunc("/getnote", controller.GetNote)
	http.HandleFunc("/calexp", controller.CalExp)
	//
	http.HandleFunc("/mynotes", controller.MyNotes) //我的笔记页面
	fmt.Println("Ip:0.0.0.0\nPort:8080")
	http.ListenAndServe(":8080", nil)

}
