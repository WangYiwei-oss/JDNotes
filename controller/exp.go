package controller

import (
	"../model"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	stlpath "path"
	"strconv"
	"strings"
)

func Expassistant(w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("funcname")
	arg := make(map[string]interface{})
	arg["login"] = model.CheckLogin(r)
	funcclass := model.GetAllExps()
	arg["funcclass"] = funcclass
	content, _ := model.RenderContentByfuncname(t)
	arg["content"] = content
	script, _ := model.RenderScriptByfuncname(t)
	arg["script"] = script
	model.WriteTemplate(w, "views/pages/expassistant.html", arg)
}

func CalIqyInput(w http.ResponseWriter, r *http.Request) {
	bt1 := r.PostFormValue("b1")
	bt2 := r.PostFormValue("b2")
	if bt1 == "" || bt2 == "" {
		w.Write([]byte("请输入标样信息"))
		return
	}
	b1, err := strconv.ParseFloat(bt1, 64)
	if err != nil {
		w.Write([]byte("请输入标样不正确"))
		return
	}
	b2, err := strconv.ParseFloat(bt2, 64)
	if err != nil {
		w.Write([]byte("请输入标样不正确"))
		return
	}
	samples := r.PostFormValue("samples")
	if samples == "" {
		w.Write([]byte("请输入样品信息"))
		return
	}
	sss := strings.Split(samples, "\n")
	res := "样品 量子效率 蓝光吸收率\n"
	for _, ss := range sss {
		s := strings.Split(ss, " ")
		if len(s) != 3 {
			w.Write([]byte("样品格式不正确\n样品名称 反射光子数 发射光子数"))
			return
		}
		a1, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			w.Write([]byte("样品光子数不能是字符"))
			return
		}
		a2, err := strconv.ParseFloat(s[2], 64)
		if err != nil {
			w.Write([]byte("样品光子数不能是字符"))
			return
		}
		iqy, abs := model.CalculateIQY(b1, b2, a1, a2)
		res += (s[0] + " " + strconv.FormatFloat(iqy, 'f', 4, 64) + " " + strconv.FormatFloat(abs, 'f', 4, 64) + "\n")
	}
	model.CleanIQYFiles()
	f, _ := os.Create("/home/wangyiwei/JDNotes/views/static/temp/iqy/iqy.txt")
	defer f.Close()
	f.Write([]byte(res))
	w.Write([]byte(res))
}

func CalIntegral(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	files := r.MultipartForm.File["file[]"]
	if len(files) == 0 {
		w.Write([]byte("请上传文件"))
		return
	}
	startwl := r.MultipartForm.Value["startwl"][0]
	endwl := r.MultipartForm.Value["endwl"][0]
	if startwl == "" {
		startwl = "-10000"
	}
	if endwl == "" {
		endwl = "10000"
	}
	swl, _ := strconv.ParseFloat(startwl, 64)
	ewl, _ := strconv.ParseFloat(endwl, 64)
	if ewl < swl {
		w.Write([]byte("请输入正确的范围"))
		return
	}
	if files == nil {
		return
	}
	model.CleanIntegralFiles()
	for i := 0; i < len(files); i++ {
		file, err := files[i].Open()
		if err != nil {
			fmt.Println(err)
			return
		}
		filename := "views/static/temp/integral" + files[i].Filename
		cur, _ := os.Create(filename)
		io.Copy(cur, file)
		cur.Close()
		file.Close()
	}
	entries, err := ioutil.ReadDir("/home/wangyiwei/JDNotes/views/static/temp/integral")
	if err != nil {
		fmt.Println(err)
	}
	content := "起始波长:" + startwl + "结束波长:" + endwl + "\n" + "样品 积分面积\n"
	for _, path := range entries {
		fi, _ := os.Open("/home/wangyiwei/JDNotes/views/static/temp/integral" + path.Name())
		br := bufio.NewReader(fi)
		counttwo := 0
		var tempdata1, tempdata2 [2]float64
		startread := false
		datax := make([]float64, 0)
		datay := make([]float64, 0)
		for {
			a, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			n, datas := model.LineDataRead(string(a))
			if n != 2 {
				counttwo = 0
				continue
			}
			tempx, err := strconv.ParseFloat(datas[0], 64)
			if err != nil {
				counttwo = 0
				continue
			}
			tempy, err := strconv.ParseFloat(datas[1], 64)
			if err != nil {
				counttwo = 0
				continue
			}
			if n == 2 {
				if counttwo == 0 && startread == false {
					tempdata1[0] = tempx
					tempdata1[1] = tempy
				}
				if counttwo == 1 && startread == false {
					tempdata2[0] = tempx
					tempdata2[1] = tempy
				}
				if counttwo == 2 && startread == false {
					startread = true
				}
				counttwo++
			} else {
				counttwo = 0
			}
			if startread == true && n != 2 {
				break
			}
			if startread == true {
				if counttwo == 3 {
					datax = append(datax, tempdata1[0])
					datax = append(datax, tempdata2[0])
					datay = append(datay, tempdata1[1])
					datay = append(datay, tempdata2[1])
					datax = append(datax, tempx)
					datay = append(datay, tempy)
				} else {
					datax = append(datax, tempx)
					datay = append(datay, tempy)
				}
			}
		}
		res := model.CalculateIntegral(datax, datay, swl, ewl)
		content += strings.TrimSuffix(path.Name(), stlpath.Ext(path.Name())) + " " + strconv.FormatFloat(res, 'f', 4, 64) + "\n"
	}
	res, _ := os.Create("views/static/temp/integral/integral.txt")
	defer res.Close()
	res.Write([]byte(content))
	w.Write([]byte(content))
}
