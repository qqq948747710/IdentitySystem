//设置动态背景图片
$(document).ready(function() {
			var imgurl = "../static/images/o_bg" + Math.floor(Math.random() * 23) + ".jpg";
			$('.bg').css('background-image', "url(\"" + imgurl + "\")");
		})
