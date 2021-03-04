package model

import (
	"io/ioutil"
	"os"
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
