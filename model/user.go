package model

import (
	"../utils"
)

type User struct {
	Id       int
	Username string `db:"name"`
	Password string
	Number   string
	Email    string
	Img_path string `db:"imgpath"`
}

func (u *User) SelectByUsername() (err error) {
	sql := "SELECT id,name,password,number,email,imgpath FROM users WHERE name=?"
	err = utils.Db.QueryRow(sql, u.Username).Scan(&u.Id, &u.Username, &u.Password, &u.Number, &u.Email, &u.Img_path)
	return
}

func (u *User) SelectByNameAndPassword() (err error) {
	sql := "SELECT id,name,password,number,email,imgpath FROM users WHERE name=? AND password=?"
	err = utils.Db.QueryRow(sql, u.Username, u.Password).Scan(&u.Id, &u.Username, &u.Password, &u.Number, &u.Email, &u.Img_path)
	return
}

func (u *User) InsertIntoSql() (err error) {
	sql := "INSERT INTO users(name,password,number,email,imgpath) VALUES (?,?,?,?,?)"
	err = utils.Db.QueryRow(sql, u.Username, u.Password, u.Number, u.Email, u.Img_path).Scan(&u.Id)
	return
}
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
