$(document).ready(function(){
	var licolor='list-group-item list-group-item-dark';
	var number=1;
	$('.edulist li').each(function(){
		if(number%2==0){
			licolor='list-group-item list-group-item-dark';
		}else{
			licolor='list-group-item list-group-item-light';
		}
		$(this).prop('class',licolor)
		number++;
	});
})