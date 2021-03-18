package model

import (
	"unicode"

	"../utils"
)

//用户类型
type User struct {
	Id       int
	Username string `db:"name"`
	Password string
	Number   string
	Email    string
	Img_path string `db:"imgpath"`
}

//用户类型方法
//根据用户名在sql中查找
func (u *User) SelectByUsername() (err error) {
	sql := "SELECT id,name,password,number,email,imgpath FROM users WHERE name=?"
	err = utils.Db.QueryRow(sql, u.Username).Scan(&u.Id, &u.Username, &u.Password, &u.Number, &u.Email, &u.Img_path)
	return
}

//根据id在sql中查找
func (u *User) SelectById() (err error) {
	sql := "SELECT id,name,password,number,email,imgpath FROM users WHERE id=?"
	err = utils.Db.QueryRow(sql, u.Id).Scan(&u.Id, &u.Username, &u.Password, &u.Number, &u.Email, &u.Img_path)
	return
}

//根据用户名和密码在sql中查找
func (u *User) SelectByNameAndPassword() (err error) {
	sql := "SELECT id,name,password,number,email,imgpath FROM users WHERE name=? AND password=?"
	err = utils.Db.QueryRow(sql, u.Username, u.Password).Scan(&u.Id, &u.Username, &u.Password, &u.Number, &u.Email, &u.Img_path)
	return
}

//将用户信息插入mysql
func (u *User) InsertIntoSql() (err error) {
	sql := "INSERT INTO users(name,password,number,email,imgpath) VALUES (?,?,?,?,?)"
	err = utils.Db.QueryRow(sql, u.Username, u.Password, u.Number, u.Email, u.Img_path).Scan(&u.Id)
	return
}

//检查邮箱格式是否正确
func (u *User) CheckEmailFormat() (b bool) {
	if u.Email == "" {
		return true
	}
	count := 0
	for _, a := range u.Email {
		if a == int32('@') && count > 0 && count != len(u.Email)-1 {
			return true
		}
		count++
	}
	return false
}

//检查密码是否符合要求：由英文和数字组成
func (u *User) CheckPassword() (b bool) {
	state := 0 //0:初始状态,1:字母,2:数字
	count := 1
	for _, c := range u.Password {
		if unicode.IsLetter(c) {
			if state == 0 {
				state = 1
			}
			if state == 2 {
				state = 1
				count++
			}
		}
		if unicode.IsNumber(c) {
			if state == 0 {
				state = 2
			}
			if state == 1 {
				state = 2
				count++
			}
		}
	}
	if count > 1 {
		return true
	} else {
		return false
	}
}
func SearchUsernotes(userid int) []Note {
	sql := "SELECT title,firstclass,secondclass FROM notes WHERE id=?"
	rows, _ := utils.Db.Query(sql, userid)
	ns := make([]Note, 0)
	for rows.Next() {
		var n Note
		rows.Scan(&n.Title, &n.Firstclass, &n.Secondclass)
		n.Id = userid
		ns = append(ns, n)
	}
	return ns
}

type Expid struct {
	Id    int
	Title string
}

func (e *Expid) DeleteExpByIdAndTitle() {
	sql := "DELETE FROM userexps WHERE id=? AND title=?"
	utils.Db.Exec(sql, e.Id, e.Title)
}

func SearchUserexps(userid int) []Expid {
	sql := "SELECT title FROM userexps WHERE id=?"
	rows, _ := utils.Db.Query(sql, userid)
	ns := make([]Expid, 0)
	for rows.Next() {
		title := Expid{}
		rows.Scan(&title.Title)
		title.Id = userid
		ns = append(ns, title)
	}
	return ns
}
