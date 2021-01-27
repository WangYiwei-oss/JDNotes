CREATE TABLE IF NOT EXISTS users
(
	id MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
	name varchar(20)NOT NULL,
	password varchar(30)NOT NULL,
	number varchar(11) NULL,
	email varchar(20)NULL,
	imgpath varchar(100) NOT NULL DEFAULT "views/static/img/defaultPhoto.png",
	PRIMARY KEY(id)
)ENGINE=InnoDB;
	
