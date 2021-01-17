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
	fmt.Println(user.CheckPassword())
}
