package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Db  *sql.DB
	err error
)

func init() {
	Db, err = sql.Open("mysql", "wangyiwei:123@tcp(127.0.0.1:3306)/JDNotes")
	if err != nil {
		panic(err)
	}
	if err = Db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("数据库连接成功")
}
