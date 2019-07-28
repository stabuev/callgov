function doLogin() {
	var login = document.getElementById("loginName").value;
	var password = document.getElementById("loginPassword").value;

	var req = new XMLHttpRequest();

	req.onload = function(e) {
				window.location.replace("/");
	};
	req.onerror = function(e) {
		document.getElementById("loginWarn").innerText = "Service temporarily unavailable.";
	};

	req.open("POST", "/login");
	req.setRequestHeader("Content-Type", "application/json");
	req.send('{"login":"' + login + '","password":"' + password + '"}');
}

function createObr() {

    // Запрос: {"id":uint32, "title":"TITLE", "content":"CONTENT", "public":0|1,  "state":"draft|sign|post","address":"ADDRESS"}


	var id = 0;
    var title = document.getElementById("title").value;
    var content = document.getElementById("content").value;
    var publ = document.getElementById("public").value;
    var state = document.getElementById("state").value;
    var address = document.getElementById("address").value;

	var req = new XMLHttpRequest();
    if (publ == "public"){
        publ = 1;
    } else {
        publ = 0;
    }
    console.log(publ);
	req.onload = function(e) {
				window.location.replace("/");
	};
	req.onerror = function(e) {
		document.getElementById("loginWarn").innerText = "Service temporarily unavailable.";
	};

	req.open("POST", "/json/obr");
	req.setRequestHeader("Content-Type", "application/json");
	req.send('{"id":"' + id + '","title":"' + title + '","title":"' + content + '","public":"' + publ + '", "state":"' + state + '", "address":"' + address + '"}');
}