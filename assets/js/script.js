window.onload = function(){


function showContent(){
	var output = '';

	output += '<div class="naming"> <div class="percentage">Проценты</div><div class="recipe">Диагноз и рекомендации</div></div>'
    
	document.getElementById('content').innerHTML = output;
	

}

document.getElementById('searchButton').onclick =showContent;

}