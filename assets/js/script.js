var wsr = new WebSocket("ws://" + window.location.host + "/wsr");



window.onload = function () {


	function showContent() {

		var z = document.getElementsByName('a');
		for (var i = 0; i < z.length; i++)  {

    if  (z[i].checked) {

     var  petNumber = i; break;

    }

}

alert(s);



		var output = '';

		output += '<div class="naming"> <div class="percentage">Проценты</div><div class="recipe">Диагноз и рекомендации</div></div>'

		document.getElementById('content').innerHTML = output;

		var query = document.getElementById("theSearch").innerText;
		var req = '{"name": "'+petNumber'", "query":"' + query + '"}';
		wsr.send(JSON.parse(req))


	}

	document.getElementById('searchButton').onclick = showContent;







}


	/*
wsr.onmessage = function(event){

	var data=JSON.parse(event.data); //4 поля 
	/*data.name - имя болезни
	pets -
	symp- 
	therapy - лечение

}


*/














