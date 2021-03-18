package model

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var user User
	user.Username = "wangyiwei"
	user.Email = "123@qq.com"
	user.Password = "wwww"
	//user.SelectByUsername()
	//fmt.Println(user)
	//fmt.Println(user.CheckEmailFormat())
	//fmt.Println(user.CheckPassword())
	//ScanAllNotesToMysql("/home/wangyiwei/JDNotes/views/pages/notes")
	//GetAllNotes()
	//GetAllFirstclass()
	//a := SelectByFirstClass("Golang")
	//fmt.Println(a)
	//CleanStaticFiles()
	cal := NewCalculator()
	a := []string{"(", "a", "*", "(", "b", "-", "c", ")", ")", "+", "(", "d", ")"}
	datas := make(map[string]float64)
	datas["a"] = 1.2
	datas["b"] = 2.3
	datas["c"] = 3.4
	datas["d"] = 2
	cal.CalculateFunc("x", a, datas)
	b := []string{"x", "*", "(", "a", "+", "b", ")"}
	cal.CalculateFunc("y", b, datas)
	fmt.Println(datas)
}
