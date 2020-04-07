$(document).ready(function() {
	$('#cernonamequery').css('display', 'none');
	$('#entityidquery').css('display', 'block');
})

function changeop() {
	var op = $('#op option:selected');
	switch (op.val()) {
		case '身份证号码':
			$('#cernonamequery').css('display', 'none');
			$('#entityidquery').css('display', 'block');
			break;
		case '证书编号和姓名':
			$('#cernonamequery').css('display', 'block');
			$('#entityidquery').css('display', 'none');
			break;
	}
}
