package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"../utils"
)

type ExpAssistant struct {
	Id       int
	Title    string
	Funcname string
}

//返回数据库中所有的实验信息
func GetAllExps() *[]ExpAssistant {
	sql := "SELECT id,title,funcname FROM exps"
	rows, _ := utils.Db.Query(sql)
	ss := make([]ExpAssistant, 1)
	for rows.Next() {
		var exp ExpAssistant
		rows.Scan(&exp.Id, &exp.Title, &exp.Funcname)
		ss = append(ss, exp)

	}
	return &ss
}

//根据函数名返回html
func RenderContentByfuncname(t string) (string, error) {
	path := "views/pages/exps/"
	path += strings.ToLower(t)
	path += ".html"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

//根据函数名返回JS
func RenderScriptByfuncname(t string) (string, error) {
	path := "views/static/JS/"
	path += strings.ToLower(t)
	path += ".js"
	js, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

//根据函数名返回代码路径
func pathByfuncname(t string) string {
	path := "views/pages/exps"
	path += t
	path += ".html"
	return path
}

//根据x和y的数据切片计算积分面积
func CalculateIntegral(xs, ys []float64, startwl, endwl float64) float64 {
	if len(xs) != len(ys) {
		fmt.Println("xy数据长度不一致")
		return 0
	}
	var res float64
	for i := 1; i < len(xs); i++ {
		if xs[i] < startwl {
			continue
		}
		if xs[i] > endwl {
			break
		}
		res += ((ys[i-1] + ys[i]) * 0.5 * (xs[i] - xs[i-1]))
	}
	return res
}

//根据四个光子数计算量子效率和蓝光吸收率
func CalculateIQY(b1, b2, a1, a2 float64) (float64, float64) {
	abs := (b1 - a1) / b1
	iqy := (b2 - a2) / (a1 - b1)
	return iqy, abs
}

//判断切片是否有重复的元素
func CheckRepeatString(ss []string) bool {
	n := make(map[string]bool)
	for _, value := range ss {
		if _, ok := n[value]; ok {
			return true
		}
		n[value] = true
	}
	return false
}

//方程的split
func FuncSplit(f string) []string {
	var temp string
	res := make([]string, 0)
	for _, a := range f {
		b := string(a)
		if b == "=" || b == "+" || b == "-" || b == "*" || b == "/" || b == "(" || b == ")" {
			if temp != "" {
				res = append(res, temp)
			}
			temp = ""
			res = append(res, b)
		} else {
			temp += b
		}
	}
	if temp != "" {
		res = append(res, temp)
	}
	return res
}

//检查方程的格式是否正确
func CheckFuncFormat(splitedfunc []string) bool {
	if splitedfunc[1] == "=" {
		return true
	} else {
		return false
	}
}

//判断方程是否包含了所有的结果
func CheckFuncInclueAllRes(ress, funcs []string) bool {
	res := make(map[string]bool)
	for _, r := range ress {
		res[r] = true
	}
	for _, f := range funcs {
		sf := FuncSplit(f)
		if _, ok := res[sf[0]]; ok {
			delete(res, sf[0])
		}
	}
	if len(res) == 0 {
		return true
	}
	return false
}

//解析output
func ParseOutput(output string) ([]string, []string) {
	parm := make([]string, 0)
	content := make([]string, 0)
	tempparm := ""
	tempcontent := ""
	count := 0 //0:读content，1读到第一个{，2读parm，3读到第一个}
	for _, b := range output {
		if b == '}' && count == 2 {
			parm = append(parm, tempparm)
			content = append(content, tempparm)
			count++
			tempparm = ""
		}
		if b == '}' && count == 3 {
			count = 0
		}
		if b == '{' && count == 0 {
			content = append(content, tempcontent)
			count++
			tempcontent = ""
		}
		if b == '{' && count == 1 {
			count++
		}
		if count == 0 && b != '}' {
			tempcontent += string(b)
		}
		if count == 2 && b != '{' {
			tempparm += string(b)
		}
	}
	content = append(content, tempcontent)
	return parm, content
}

//判断参数是否全部出现在方程或者输出中
func CheckFuncContentInclueAllParm(parms, funcs, parsedcontent []string) bool {
	allparm := make(map[string]bool)
	for _, p := range parms {
		allparm[p] = true
	}
	for _, f := range funcs {
		sf := FuncSplit(f)
		for _, parm := range sf {
			if _, ok := allparm[parm]; ok {
				delete(allparm, parm)
			}
		}
	}
	for _, c := range parsedcontent {
		if _, ok := allparm[c]; ok {
			delete(allparm, c)
		}
	}
	if len(allparm) == 0 {
		return true
	} else {
		return false
	}
}

//处理parmnamesJSON,parmsJSON,funcsJSON,contentJSON用的结构体
type FuncPair struct {
	result string
	parms  []string
}

func (f *FuncPair) GetResult() string {
	return f.result
}
func (f *FuncPair) GetParms() []string {
	return f.parms
}

type UserExpDealer struct {
	parms        [][2]string
	funcs        []FuncPair
	contentparms map[string]string
	content      []string
	datas        map[string]float64
}

func (ued *UserExpDealer) GetFuncs() []FuncPair {
	return ued.funcs
}
func (ued *UserExpDealer) GetDatas() map[string]float64 {
	return ued.datas
}
func (ued *UserExpDealer) WriteDatas(datasJson string) {
	ued.datas = make(map[string]float64)
	datas := make([]string, 0)
	json.Unmarshal([]byte(datasJson), &datas)
	for i := 0; i < len(ued.parms); i++ {
		ued.datas[ued.parms[i][1]], _ = strconv.ParseFloat(datas[i], 64)
	}
}
func (ued *UserExpDealer) InitUserExpDealer(parmnamesJSON, parmsJSON, funcsJSON, contentJSON string) {
	ued.funcs = make([]FuncPair, 0)
	parmnames := make([]string, 0)
	ued.contentparms = make(map[string]string)
	err := json.Unmarshal([]byte(parmnamesJSON), &parmnames)
	if err != nil {
		fmt.Println("parmnamesjson解析失败", err)
		return
	}
	parms := make([]string, 0)
	err = json.Unmarshal([]byte(parmsJSON), &parms)
	if err != nil {
		fmt.Println("parmsjson解析失败", err)
		return
	}
	funcs := make([]string, 0)
	err = json.Unmarshal([]byte(funcsJSON), &funcs)
	if err != nil {
		fmt.Println("parmsjson解析失败", err)
		return
	}
	if len(parmnames) != len(parms) {
		fmt.Println("参数和参数名长度不相同")
		return
	}
	for i := 0; i < len(parmnames); i++ {
		var temp [2]string
		temp[0] = parmnames[i]
		temp[1] = parms[i]
		ued.parms = append(ued.parms, temp)
	}
	outputparms, contents := ParseOutput(contentJSON)
	ued.content = contents
	for _, p := range outputparms {
		ued.contentparms[p] = ""
	}
	for _, fs := range funcs {
		f := FuncSplit(fs)
		fp := FuncPair{}
		fp.result = f[0]
		fp.parms = f[2:]
		ued.funcs = append(ued.funcs, fp)
	}
}
func sliceToString(s []string) string {
	res := ""
	for _, a := range s {
		res += a
	}
	return res
}
func (ued *UserExpDealer) RenderParmHtml() string {
	html := "<button id='showinfo'>查看实验信息&gt</button><br/>"
	for i := 0; i < len(ued.parms); i++ {
		html = html + "<p class='parmlabel' >" + ued.parms[i][0] + ":</p>"
		html = html + "<input type='text' class='parms'><br/>"
	}
	parm := ""
	for _, p := range ued.parms {
		parm = parm + p[0] + ":" + p[1] + "<br/>"
	}
	funcs := ""
	for _, fp := range ued.funcs {
		funcs = funcs + fp.result + "=" + sliceToString(fp.parms) + "<br/>"
	}
	content := ""
	for _, a := range ued.content {
		if _, ok := ued.contentparms[a]; ok {
			content += "<font style='color:red'> " + a + " </font>"
		} else {
			content += a
		}
	}
	html = html + "\n<div id='info'><p>你的参数为:<br/>" + parm + "你的结果及方程为：<br/>" + funcs + "输出格式为：<br/>" + content + "</p></div>\n"
	html += "<button id='submit'>计算</button>"
	return html
}
func (ued *UserExpDealer) RenderOutput() []byte {
	content := ""
	for _, a := range ued.content {
		if b, ok := ued.GetDatas()[a]; ok {
			content += strconv.FormatFloat(b, 'f', 4, 64)
		} else {
			content += a
		}
	}
	return []byte(content)
}
