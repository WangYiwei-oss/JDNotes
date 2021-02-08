package model

type ExpAssistant struct {
	Name     string
	Path     string
	Username string
}

//根据输入的信息来生成布局
func RenderContentByType(t string) string {
	switch t {
	case "cal_iqy_input":
		return "<p>样品光子数的输入格式为:<br/>样品1名称&nbsp样品1反射光子数&nbsp样品1发射光子数<br/>样品2名称&nbsp样品2反射光子数&nbsp样品2发射光子数<br/>...</p><input id='b1' type='text' placeholder='请输入标样反射光子数'><input id='b2' type='text' placeholder='请输入标样发射光子数'><br/><textarea id='samples' placeholder='请输入样品光子数'></textarea><br/><button id='btn'>计算</button>"
	case "cal_iqy_file":
	case "cal_integral":
	case "fit_lifetime":
	case "cal_exp1":
	case "cal_exp2":
	}
	return ""
}

//根据输入的信息来生成script
func RenderScriptByType(t string) string {
	switch t {
	case "cal_iqy_input":
		return "$(function(){$('#btn').click(function(){var parm={'b1':$('#b1').val(),'b2':$('#b2').val(),'samples':$('#samples').val()};$.post('/cal_iqy_input_post',parm,function(res){$('#result').html(res)});});});"
	case "cal_iqy_file":
	case "cal_integral":
	case "fit_lifttime":
	case "cal_exp1":
	case "cal_exp2":
	}
	return ""
}
