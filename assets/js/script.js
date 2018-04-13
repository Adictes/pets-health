var wsr = new WebSocket("ws://" + window.location.host + "/wsr");



window.onload = function () {


	function showContent() {
		var output = '';

		output += '<div class="naming"> <div class="percentage">Проценты</div><div class="recipe">Диагноз и рекомендации</div></div>'

		document.getElementById('content').innerHTML = output;

		var query = document.getElementById("theSearch").innerText;
		var req = '{"name": "todo", "query":"' + query + '"}';
		wsr.send(JSON.parse(req))
	}

	document.getElementById('searchButton').onclick = showContent;

}