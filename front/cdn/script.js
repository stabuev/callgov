function doLogin() {
	recLoginUsed = true;

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