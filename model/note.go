package model

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"../utils"
)

//储存笔记的信息
type Note struct {
	Id          int
	Title       string
	Filepath    string
	Firstclass  string
	Secondclass string
	Thirdclass  string
}

//返回dir目录下所有的文件信息
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}
	return entries
}

//返回dir文件的一级二级三级分类
func getClass(dir string) (res [3]string) {
	file, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(file)
	count := 0
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err || count > 2 {
			break
		} else {
			if strings.HasPrefix(line, "<!-->") {
				res[count] = line[5 : len(line)-5]
			}
			count++
		}
	}
	return res
}

//扫描dir目录下的所有的文件信息到[]Note中
func walkDir(dir string, ns *[]Note) {
	for _, path := range dirents(dir) {
		if path.IsDir() {
			subdir := filepath.Join(dir, path.Name())
			walkDir(subdir, ns)
		} else {
			note := Note{}
			note.Id = 1
			note.Title = strings.TrimSuffix(path.Name(), ".html")
			note.Filepath = filepath.Join(dir, path.Name())
			class := getClass(note.Filepath)
			if class[0] == "" {
				note.Firstclass = "其他"
			} else {
				note.Firstclass = class[0]
			}
			note.Secondclass = class[1]
			note.Thirdclass = class[2]
			*ns = append(*ns, note)
		}
	}
}

//walkDir函数包装成可导出的
func ScanAllNotes(dir string) *[]Note {
	var ns []Note
	walkDir(dir, &ns)
	return &ns
}

//将dir下的所有的文件信息写入mysql
func ScanAllNotesToMysql(dir string) {
	ns := ScanAllNotes(dir)
	for _, n := range *ns {
		utils.Db.Exec("INSERT INTO notes(id,title,filepath,firstclass,secondclass,thirdclass)VALUES(?,?,?,?,?,?)", n.Id, n.Title, n.Filepath, n.Firstclass, n.Secondclass, n.Thirdclass)
	}
}

//获取笔记的类型，用于填充导航栏
func GetAllFirstclass() *[]string {
	sql := "SELECT firstclass FROM notes GROUP BY firstclass"
	rows, _ := utils.Db.Query(sql)
	ss := make([]string, 1)
	for rows.Next() {
		var class string
		rows.Scan(&class)
		ss = append(ss, class)
	}
	return &ss
}

//获取所有笔记的所有信息
func GetAllNotes() *[]Note {
	sql := "SELECT id,title,filepath,firstclass,secondclass,thirdclass FROM notes"
	rows, err := utils.Db.Query(sql)
	if err != nil {
		return nil
	}
	var ns []Note
	for rows.Next() {
		var n Note
		rows.Scan(&n.Id, &n.Title, &n.Filepath, &n.Firstclass, &n.Secondclass, &n.Thirdclass)
		ns = append(ns, n)
	}
	return &ns
}

//获取一级分类下所有的笔记标题信息
func SelectByFirstClass(firstclass string) *map[string][]Note {
	sql := "SELECT title,firstclass,secondclass,thirdclass FROM notes WHERE firstclass=?"
	rows, _ := utils.Db.Query(sql, firstclass)
	group_ns := make(map[string][]Note)
	for rows.Next() {
		var n Note
		rows.Scan(&n.Title, &n.Firstclass, &n.Secondclass, &n.Thirdclass)
		group_ns[n.Secondclass] = append(group_ns[n.Secondclass], n)
	}
	return &group_ns
}

//根据title打开并返回笔记内容
func OpenByTitle(title string) (buf []byte) {
	sql := "SELECT filepath FROM notes WHERE title=?"
	var (
		filepath string
		err      error
	)
	utils.Db.QueryRow(sql, title).Scan(&filepath)
	buf, err = ioutil.ReadFile(filepath)
	if err != nil {
		return nil
	}
	return
}
