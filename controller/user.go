package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"../model"
	"../utils"
)

func ToLogin(w http.ResponseWriter, r *http.Request) {
	model.WriteHTML(w, "views/pages/user/login.html")
}
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
			MaxAge:   360000,
		}
		//写入cookie
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/index", 302)
	} else {
		model.WriteTemplate(w, "views/pages/user/login.html", "用户名或密码错误")
	}
}
func ToRegister(w http.ResponseWriter, r *http.Request) {
	model.WriteHTML(w, "views/pages/user/register.html")
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
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    user.Username,
		HttpOnly: true,
		MaxAge:   3600,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/index", 302)
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
	if model.CheckLogin(r) {
		model.WriteHTML(w, "views/pages/user/mypage.html")
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}
func Unlogin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/index", 302)
	}
}
func AddNote(w http.ResponseWriter, r *http.Request) {
	model.WriteHTML(w, "views/pages/user/addnote.html")
}
func WriteNote(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		w.Write([]byte("请登陆"))
		http.Redirect(w, r, "/tologin", 302)
		return
	}
	username := cookie.Value
	dirpath := "views/pages/usernote/" + username
	if !model.PathExists(dirpath) {
		os.Mkdir(dirpath, 0777)
	}
	note := model.Note{}
	note.Title = r.FormValue("title")
	filename := dirpath + "/" + note.Title
	if model.PathExists(filename) {
		w.Write([]byte("提交了重复的文件，请修改标题"))
		return
	}
	f, err := os.Create(filename)
	defer f.Close()
	user := model.User{}
	user.Username = username
	user.SelectByUsername()
	note.Id = user.Id
	note.Firstclass = r.FormValue("firstclass")
	note.Secondclass = r.FormValue("secondclass")
	content := model.FormatUserNote(note.Id, note.Firstclass, note.Secondclass, r.FormValue("content"))
	f.Write(content)
	note.Filepath = strings.TrimPrefix(filename, "/home/wangyiwei/JDNotes/")
	note.InsertIntoMysql()
	w.Write([]byte("提交成功"))
}
func AlertUsernote(w http.ResponseWriter, r *http.Request) {
	prenote := model.Note{}
	prenote.Id, _ = strconv.Atoi(r.FormValue("id"))
	prenote.Title = r.FormValue("oldtitle")
	prenote.SelectNoteByIdAndTitle()
	//删除原来的文件
	prenote.DeleteNoteByIdAndTitle()
	note := model.Note{}
	note.Id = prenote.Id
	user := model.User{}
	user.Id = note.Id
	user.SelectById()
	username := user.Username
	dirpath := "views/pages/usernote/" + username
	if !model.PathExists(dirpath) {
		os.Mkdir(dirpath, 0777)
	}
	note.Title = r.FormValue("title")
	note.Firstclass = r.FormValue("firstclass")
	note.Secondclass = r.FormValue("secondclass")
	filename := dirpath + "/" + note.Title
	if model.PathExists(filename) {
		w.Write([]byte("提交了重复的文件，请修改标题"))
		return
	}
	f, _ := os.Create(filename)
	defer f.Close()
	content := model.FormatUserNote(note.Id, note.Firstclass, note.Secondclass, r.FormValue("content"))
	f.Write(content)
	note.Filepath = strings.TrimPrefix(filename, "/home/wangyiwei/JDNotes/")
	note.InsertIntoMysql()
	w.Write([]byte("修改成功"))
}
func UserNote(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		user := model.User{}
		user.Username = cookie.Value
		user.SelectByUsername()
		ns := model.SearchUsernotes(user.Id)
		args := make(map[string]interface{})
		args["notes"] = ns
		model.WriteTemplate(w, "views/pages/user/usernotes.html", args)
	} else {
		w.Write([]byte("请登陆"))
		http.Redirect(w, r, "/login", 302)
	}
}

//删除我的笔记中的笔记
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	note := model.Note{}
	note.Id, _ = strconv.Atoi(r.FormValue("id"))
	note.Title = r.FormValue("title")
	if note.Id != 1 && note.Title != "" {
		note.DeleteNoteByIdAndTitle()
		http.Redirect(w, r, "/mypage_mynotes", 302)
	}
}
func AlertNote(w http.ResponseWriter, r *http.Request) {
	note := model.Note{}
	note.Id, _ = strconv.Atoi(r.FormValue("id"))
	note.Title = r.FormValue("title")
	note.SelectNoteByIdAndTitle()
	f, _ := os.Open(note.Filepath)
	defer f.Close()
	content := model.ParseUsernote(f)
	args := make(map[string]interface{})
	args["title"] = note.Title
	args["id"] = note.Id
	args["content"] = content
	model.WriteTemplate(w, "views/pages/user/alertnote.html", args)
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	note := model.Note{}
	note.Id, _ = strconv.Atoi(r.FormValue("id"))
	note.Title = r.FormValue("title")
	note.SelectNoteByIdAndTitle()
	type Notejs struct {
		Title       string
		Firstclass  string
		Secondclass string
	}
	js := Notejs{}
	js.Title = note.Title
	js.Firstclass = note.Firstclass
	js.Secondclass = note.Secondclass
	data, _ := json.Marshal(js)
	w.Write(data)
}

func AddExp(w http.ResponseWriter, r *http.Request) {
	model.WriteHTML(w, "views/pages/user/addexp.html")
}
func AddUserExp(w http.ResponseWriter, r *http.Request) {
	parmsJson := r.FormValue("parms")
	parmNameJson := r.FormValue("parmname")
	ressJson := r.FormValue("ress")
	funcsJson := r.FormValue("funcs")
	output := r.FormValue("output")
	title := r.FormValue("title")
	t := 0
	err := utils.Db.QueryRow("SELECT id FROM userexps WHERE title=?", title).Scan(&t)
	if err == nil {
		w.Write([]byte("标题已被占用，请更改"))
		return
	}
	parms := make([]string, 0)
	ress := make([]string, 0)
	funcs := make([]string, 0)
	json.Unmarshal([]byte(parmsJson), &parms)
	json.Unmarshal([]byte(ressJson), &ress)
	json.Unmarshal([]byte(funcsJson), &funcs)
	if model.CheckRepeatString(parms) {
		w.Write([]byte("有重复的参数,请检查"))
		return
	}
	if model.CheckRepeatString(ress) {
		w.Write([]byte("有重复的结果,请检查"))
		return
	}
	if model.CheckRepeatString(funcs) {
		w.Write([]byte("有重复的方程,请检查"))
		return
	}
	for _, f := range funcs {
		sf := model.FuncSplit(f)
		if !model.CheckFuncFormat(sf) {
			w.Write([]byte("方程格式错误，请检查"))
			return
		}
	}
	if !model.CheckFuncInclueAllRes(ress, funcs) {
		w.Write([]byte("方程未包含所有的结果，请检查"))
		return
	}
	a, _ := model.ParseOutput(output)
	if !model.CheckFuncContentInclueAllParm(parms, funcs, a) {
		w.Write([]byte("方程结果与输出没有包含所有的参数，或输出的格式有误，请检查"))
		return
	}
	user := model.User{}
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		w.Write([]byte("请登录"))
		http.Redirect(w, r, "/login", 302)
		return
	}
	user.Username = cookie.Value
	user.SelectByUsername()
	model.WriteUserExpToMysql(user.Id, title, parmNameJson, parmsJson, funcsJson, output)
	w.Write([]byte("收到"))
}
func UserExp(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		user := model.User{}
		user.Username = cookie.Value
		user.SelectByUsername()
		ns := model.SearchUserexps(user.Id)
		args := make(map[string]interface{})
		args["exps"] = ns
		model.WriteTemplate(w, "views/pages/user/userexps.html", args)
	} else {
		w.Write([]byte("请登陆"))
		http.Redirect(w, r, "/login", 302)
	}
}
func DeleteExp(w http.ResponseWriter, r *http.Request) {
	note := model.Expid{}
	note.Id, _ = strconv.Atoi(r.FormValue("id"))
	note.Title = r.FormValue("title")
	if note.Id != 1 && note.Title != "" {
		note.DeleteExpByIdAndTitle()
		http.Redirect(w, r, "/mypage_myexps", 302)
	}
}
