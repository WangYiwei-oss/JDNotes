CREATE TABLE IF NOT EXISTS userexps
(
        id MEDIUMINT UNSIGNED NOT NULL,
	title varchar(50) NOT NULL,
        parmnames varchar(200) NOT NULL,
        parms varchar(200) NOT NULL,
	funcs VARCHAR(500) NOT NULL,
	content VARCHAR(1000) NOT  NULL,
        FOREIGN KEY(id) REFERENCES users(id)
)ENGINE=InnoDB;
