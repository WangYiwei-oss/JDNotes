package model

type User struct {
	Id       int
	Username string `db:"name"`
	Password string
	Number   string
	Email    string
	Img_path string `db:"imgpath"`
}
