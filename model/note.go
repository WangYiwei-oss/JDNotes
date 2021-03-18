package model

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
}

//将笔记信息插入mysql
func (n *Note) InsertIntoMysql() {
	if n.Secondclass == "" {
		utils.Db.Exec("INSERT INTO notes(id,title,filepath,firstclass,secondclass) VALUES(?,?,?,?,?)", n.Id, n.Title, n.Filepath, n.Firstclass, "其他")
	} else {
		utils.Db.Exec("INSERT INTO notes(id,title,filepath,firstclass,secondclass) VALUES(?,?,?,?,?)", n.Id, n.Title, n.Filepath, n.Firstclass, n.Secondclass)
	}
}

//返回dir目录下所有的文件信息
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}
	return entries
}

//返回dir文件的一级\二级\用户ID
func getClass(dir string) (res [3]string) {
	file, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(file)
	count := 0
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err || count > 3 {
			break
		} else {
			if strings.HasPrefix(line, "<!--") {
				res[count] = line[4 : len(line)-4]
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
			note.Title = strings.TrimSuffix(path.Name(), ".html")
			note.Filepath = filepath.Join(dir, path.Name())
			class := getClass(note.Filepath)
			if class[0] == "" {
				note.Firstclass = "其他"
			} else {
				note.Firstclass = class[0]
			}
			note.Secondclass = class[1]
			note.Id, _ = strconv.Atoi(class[2])
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
		utils.Db.Exec("INSERT INTO notes(id,title,filepath,firstclass,Secondclass)VALUES(?,?,?,?,?)", n.Id, n.Title, n.Filepath, n.Firstclass, n.Secondclass)
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
func GetAllFirstclassWithCookie(r *http.Request) *[]string {
	cookie, err := r.Cookie("_cookie")
	ss := make([]string, 1)
	if err != nil {
		sql := "SELECT firstclass FROM notes WHERE id=1 GROUP BY firstclass"
		rows, _ := utils.Db.Query(sql)
		for rows.Next() {
			var class string
			rows.Scan(&class)
			ss = append(ss, class)
		}
		return &ss
	} else {
		user := User{}
		user.Username = cookie.Value
		user.SelectByUsername()
		sql := "SELECT firstclass FROM notes WHERE id=1 OR id=? GROUP BY firstclass"
		rows, _ := utils.Db.Query(sql, user.Id)
		for rows.Next() {
			var class string
			rows.Scan(&class)
			ss = append(ss, class)
		}
		return &ss
	}
}

//获取所有笔记的所有信息
func GetAllNotes() *[]Note {
	sql := "SELECT id,title,filepath,firstclass,secondclass FROM notes"
	rows, err := utils.Db.Query(sql)
	if err != nil {
		return nil
	}
	var ns []Note
	for rows.Next() {
		var n Note
		rows.Scan(&n.Id, &n.Title, &n.Filepath, &n.Firstclass, &n.Secondclass)
		ns = append(ns, n)
	}
	return &ns
}

//获取一级分类下所有的笔记标题信息
func SelectByFirstClass(firstclass string) *map[string][]Note {
	sql := "SELECT title,firstclass,secondclass FROM notes WHERE firstclass=?"
	rows, _ := utils.Db.Query(sql, firstclass)
	group_ns := make(map[string][]Note)
	for rows.Next() {
		var n Note
		rows.Scan(&n.Title, &n.Firstclass, &n.Secondclass)
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

//根据Id和Title查找笔记
func (n *Note) SelectNoteByIdAndTitle() {
	sql := "SELECT filepath,firstclass,secondclass FROM notes WHERE id=? AND title=?"
	row := utils.Db.QueryRow(sql, n.Id, n.Title)
	row.Scan(&n.Filepath, &n.Firstclass, &n.Secondclass)
}

//根据Id和Title删除笔记
func (n *Note) DeleteNoteByIdAndTitle() {
	n.SelectNoteByIdAndTitle()
	os.Remove("/home/wangyiwei/JDNotes/" + n.Filepath)
	sql := "DELETE FROM notes WHERE id=? AND title=?"
	utils.Db.Exec(sql, n.Id, n.Title)
}

//parse文件格式
func ParseUsernote(f *os.File) string {
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	scanner.Scan()
	scanner.Scan()
	var res string
	for scanner.Scan() {
		res += scanner.Text() + "\n"
	}
	res = strings.TrimPrefix(res, "<pre>")
	res = strings.TrimSuffix(res, "</pre>\n")
	return res
}
