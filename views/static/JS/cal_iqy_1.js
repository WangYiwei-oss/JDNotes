$(function(){
	$('#btn').click(function(){
		var parm={'b1':$('#b1').val(),'b2':$('#b2').val(),'samples':$('#samples').val()};
		result=$('#result');
		result.css("display","block")
		$('#btn_save').css("display","block")
		$.post('/cal_iqy_input_post',parm,function(res){
			$('#result').text(res)
		});
	});
});
