CREATE TABLE IF NOT EXISTS users
(
	id int NOT NULL AUTO_INCREMENT,
	name varchar(50)NOT NULL,
	password varchar(50)NOT NULL,
	number varchar(11) NULL,
	email varchar(50)NULL,
	imgpath varchar(50) NULL DEFAULT "view/static/img/defaultlogo.png",
	PRIMARY KEY(id)
)
	
