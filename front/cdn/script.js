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
	req.onload = function(e) {
		window.location.replace("/");
	};
	req.onerror = function(e) {
		document.getElementById("loginWarn").innerText = "Service temporarily unavailable.";
    };
    

	req.open("POST", "/json/obr");
	req.setRequestHeader("Content-Type", "application/json");
	req.send('{"id":' + id + ',"title":"' + title + '","content":"' + content + '","public":' + publ + ', "state":"' + state + '", "address":"' + address + '"}');
}

function getObrList(){
    var content = document.getElementById("obrList")
	while (content.firstChild) {
		content.removeChild(content.firstChild);
	}
	var req = new XMLHttpRequest();
	req.open("GET", "/json/obrlist", true);
	req.onload = function (e) {
            var msg = JSON.parse(req.responseText);
            for(var obr of msg.obr){
                var msgSign = "";
                if (obr[10]==0){
                    msgSign = ', <span class="text-danger">не подписано</span>'
                } else if (obr[10]==1){
                    msgSign = ', <span class="text-success">подписано</span>'
                }
                var div=document.createElement("div");
                div.classList.add("card", "m-3")
                div.innerHTML = 
            '<div class="row no-gutters">' +
              '<div class="col-md-12">' +
                '<div class="card-body">' +
                  '<a href="/detail?id=' + obr[0] +'"><h5 class="card-title">' + obr[1] + '</h5></a>' +
                  '<p class="card-text">' + obr[2].slice(0, 250) + '... </p>' +
                  '<div class="row no-gutters">' +
                      '<div class="col-md-6">' +
                  '<p class="card-text">' +
                      '<small>Автор: <span class="text-muted">' + obr[3]+ '</span></small><br />' +
                    '<small>Адресат: <span class="text-muted">'+ obr[6] + '</span></small></div>' +
                    '<div class="col-md-6">' +
                        '<p class="card-text">' +
   '<small>Подписей: <span class="text-muted">'+ obr[9] + msgSign + '</span></small></br>' +
                  '<small>Комментариев: <span class="text-muted">'+ obr[11] +'</span></small></p>' +
                  '</div></div>' +
                '</div>' +
              '</div>' +
            '</div>'
            content.appendChild(div);

                
            }
    }    
    req.send();       
}

function getObr(id){
    var content = document.getElementById("obrList")
	while (content.firstChild) {
		content.removeChild(content.firstChild);
	}
	var req = new XMLHttpRequest();
    req.open("POST", "/json/obrlist");
    req.setRequestHeader("Content-Type", "application/json");
    req.send('{"id":' + id + '}'); 
	req.onload = function (e) {
            var msg = JSON.parse(req.responseText);
            for(var obr of msg.obr){
                var msgSign = "";
                var buttonSign = ""
                if (obr[10]==0){
                    msgSign = ', <span class="text-danger">не подписано</span>'
                    buttonSign = '<button type="button" class="btn btn-success m-2" onclick="signObr('+obr[0]+')">Подписать обращение</button>'
                } else if (obr[10]==1){
                    msgSign = ', <span class="text-success">подписано</span>'
                }
                var div=document.createElement("div");
                div.classList.add("card", "m-3")
                div.innerHTML = 
            '<div class="row no-gutters">' +
              '<div class="col-md-12">' +
                '<div class="card-body">' +
                  '<h2 class="card-title">' + obr[1] + '</h2>' +
                  '<p class="card-text">' + obr[2] + '</p>' +
                  '<div class="row no-gutters">' +
                      '<div class="col-md-6">' +
                  '<p class="card-text">' +
                      '<small>Автор: <span class="text-muted">' + obr[3]+ '</span></small><br />' +
                    '<small>Адресат: <span class="text-muted">'+ obr[6] + '</span></small></div>' +
                    '<div class="col-md-6">' +
                        '<p class="card-text">' +
   '<small>Подписей: <span class="text-muted">'+ obr[9] + msgSign + '</span></small></br>' +
                  '<small>Комментариев: <span class="text-muted">'+ obr[11] +'</span></small></p>' +
                  '</div></div>' +
                  buttonSign +
                '</div>' +
              '</div>' +
            '</div>'
            content.appendChild(div);

                
        }
    }          
}

function signObr(id){
    var req = new XMLHttpRequest();
    req.open("POST", "/json/sign");
    req.setRequestHeader("Content-Type", "application/json");
    req.send('{"id":' + id + '}'); 
    getObr(id);
}

function getSign(id){
    var content = document.getElementById("signList")
	while (content.firstChild) {
		content.removeChild(content.firstChild);
	}
	var req = new XMLHttpRequest();
    req.open("POST", "/json/signlist");
    req.setRequestHeader("Content-Type", "application/json");
    req.send('{"id":' + id + '}'); 
	req.onload = function (e) {
            var msg = JSON.parse(req.responseText);
            for(var sign of msg.sign){
            var div = document.createElement("div")
            div.innerText = sign[1];
            content.appendChild(div);
    }
}
}

function createComment(id) {
    console.log(id)

    // Запрос: {"id":uint32, "title":"TITLE", "content":"CONTENT", "public":0|1,  "state":"draft|sign|post","address":"ADDRESS"}


    var content = document.getElementById("comm-content").value;

	var req = new XMLHttpRequest();
	req.onload = function(e) {
		window.location.replace("/detail?id="+ id);
	};
    

	req.open("POST", "/json/comment");
	req.setRequestHeader("Content-Type", "application/json");
	req.send('{"id":' + id + ',"content":"' + content + '"}');
}

function getComments(id) {
    var content = document.getElementById("commentList")
	while (content.firstChild) {
		content.removeChild(content.firstChild);
	}
	var req = new XMLHttpRequest();
    req.open("POST", "/json/commentlist");
    req.setRequestHeader("Content-Type", "application/json");
    req.send('{"id":' + id + '}'); 
	req.onload = function (e) {
            var msg = JSON.parse(req.responseText);
            for(var comment of msg.comment){
            var div = document.createElement("div")
            div.classList.add("m-3")
            var div2 = document.createElement("div")
            div2.innerText = comment[0];
            var small = document.createElement("small");
            small.innerText = comment[1] + " // " + comment[2]
            div.append(div2)
            div.appendChild(small)
            content.appendChild(div);
    }
}

}