$(function(){
        $('#btn').click(function(){
		var formData = new FormData();
		for(var i=0; i<$('#fileinput')[0].files.length;i++){
                	formData.append('file[]', $('#fileinput')[0].files[i]);
                }
		formData.append('startwl',$('#startwl').val())
		formData.append('endwl',$('#endwl').val())
		$.ajax({
			url:"/cal_integral",
			type:'POST',
			data:formData,
			contentType: false,    // ajax上传文件时这两个参数需要设置成false
             		processData: false,
			success:function(res){
				$('#result').css('display','block');
				$('#result').text(res);
				$('#btn_save').css('display','block');
			},
			error:function(){
			}
		});
        });
});
