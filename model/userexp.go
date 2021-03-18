package model

import (
	"../utils"
)

type Userfunc struct {
	Id    int
	Title string
}

func WriteUserExpToMysql(id int, title, parmnames, parms, funcs, content string) {
	sql := "insert into userexps(id,title,parmnames,parms,funcs,content)VALUES(?,?,?,?,?,?)"
	utils.Db.Exec(sql, id, title, parmnames, parms, funcs, content)
}

func GetAllUserExps() []Userfunc {
	userfuncs := make([]Userfunc, 0)
	sql := "SELECT id,title FROM userexps"
	rows, _ := utils.Db.Query(sql)
	for rows.Next() {
		var f Userfunc
		rows.Scan(&f.Id, &f.Title)
		userfuncs = append(userfuncs, f)
	}
	return userfuncs
}
