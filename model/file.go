package model

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func CleanFiles(path string) {
	fileinfos, _ := ioutil.ReadDir(path)
	for _, fileinfo := range fileinfos {
		filename := path + "/" + fileinfo.Name()
		os.Remove(filename)
	}
}
func CleanIntegralFiles() {
	CleanFiles("/home/wangyiwei/JDNotes/views/static/temp/integral")
}
func CleanIQYFiles() {
	CleanFiles("/home/wangyiwei/JDNotes/views/static/temp/iqy")
}

func LineDataRead(line string) (int, []string) {
	datas := strings.Fields(line)
	return len(datas), datas
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func FormatUserNote(id int, firstclass, secondclass, content string) []byte {
	if secondclass == "" {
		res := []byte("<!--" + firstclass + "-->\n" + "<!--其他-->\n" + "<!--" + strconv.Itoa(id) + "-->\n" + "<pre>" + content + "</pre>")
		return res
	} else {
		res := []byte("<!--" + firstclass + "-->\n" + "<!--" + secondclass + "-->\n" + "<!--" + strconv.Itoa(id) + "-->\n" + "<pre>" + content + "</pre>")
		return res
	}
}
