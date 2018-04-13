var wsr = new WebSocket("ws://" + window.location.host + "/wsr");

var dict = { 0: "собака", 1: "кошка" };

window.onload = function () {


	function showContent() {

		var z = document.getElementsByName('a');
		for (var i = 0; i < z.length; i++) {

			if (z[i].checked) {

				var pet = dict[i]; break;

			}

		}


		var output = '';

		output += '<div class="naming"> <div class="percentage">Проценты</div><div class="recipe">Диагноз и рекомендации</div></div>';

		document.getElementById('content').innerHTML = output;

		var query = document.getElementById("theSearch").value;
		var req = '{ "name":"' + pet + '", "query": "' + query + '" }';
		wsr.send(req);

	}

	document.getElementById('searchButton').onclick = showContent;







}



wsr.onmessage = function (event) {

	//var data=JSON.parse(event.data); 
	var dataName = 'имя болезни';
	var dataTherapy = 'Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.';

	var output = '<br><div class="naming"> <div class="percentage">Проценты</div><div class="recipe">' + dataName + '  ' + dataTherapy + '</div></div>';

	document.getElementById('content').innerHTML += output;
}
