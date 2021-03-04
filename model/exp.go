package model

import (
	"fmt"
	"io/ioutil"
	"strings"

	"../utils"
)

type ExpAssistant struct {
	Id       int
	Title    string
	Funcname string
}

//返回数据库中所有的笔记信息
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
	iqy := (a1 - b1) / (b2 - a2)
	return iqy, abs
}
