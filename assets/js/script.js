var wsr = new WebSocket("ws://" + window.location.host + "/wsr");

var dict = { 0: "собака", 1: "кошка" };

window.onload = function() {
    function showContent() {
        var z = document.getElementsByName("a");
        for (var i = 0; i < z.length; i++) {
            if (z[i].checked) {
                var pet = dict[i];
                break;
            }
        }

        var output = "";

        output +=
            '<div class="naming columns"> <div class="percentage column is-1">Вероятность болезни</div><div class="recipe column is-9">Диагноз и рекомендации</div></div>';

        document.getElementById("content").innerHTML = output;

        var query = document.getElementById("theSearch").value;
        var req = '{ "name":"' + pet + '", "query": "' + query + '" }';
        wsr.send(req);
    }

    document.getElementById("searchButton").onclick = showContent;
};

wsr.onmessage = function(event) {
    var data = JSON.parse(event.data);

    var output =
        '<div class="naming"><div class="percentage">*(тут будут проценты)</div><div class="recipe">' +
        data.name +
        ":  " +
        data.therapy +
        "</div></div>";

    document.getElementById("content").innerHTML += output;
};
