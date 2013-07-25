package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const echoPath = "/echo"

func main() {
	http.HandleFunc("/", Home)
	http.Handle(echoPath, logHeader{websocket.Handler(Echo)})
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatalln(err)
	}
}

type logHeader struct {
	h http.Handler
}

func (h logHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		log.Println(k, v)
	}
	h.h.ServeHTTP(w, r)
}

func Echo(ws *websocket.Conn) {
	defer ws.Close()
	log.Println("got conn", ws)
	buf := make([]byte, 64*1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			return
		}
		_, err = ws.Write(buf[:n])
		if err != nil {
			return
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	body := strings.NewReader(home)
	http.ServeContent(w, r, "index.html", time.Time{}, body)
}

// From http://www.websocket.org/echo.html
const home = `
<!DOCTYPE html>
<meta charset="utf-8" />
<title>WebSocket Test</title>
<script language="javascript" type="text/javascript">
var wsUri = "ws://" + document.location.host + "` + echoPath + `";
var output, msg;
function init() {
	output = document.getElementById("output");
	msg = document.getElementById("msg");
	document.getElementById("theform").addEventListener("submit", sendmsg);
	testWebSocket();
}

function testWebSocket() {
	websocket = new WebSocket(wsUri);
	websocket.onopen = function(evt) {
		onOpen(evt)
	};
	websocket.onclose = function(evt) {
		onClose(evt)
	};
	websocket.onmessage = function(evt) {
		onMessage(evt)
	};
	websocket.onerror = function(evt) {
		onError(evt)
	};
}

function onOpen(evt) {
	writeToScreen("CONNECTED", "black");
	doSend("WebSocket rocks");
}

function onClose(evt) {
	writeToScreen("DISCONNECTED", "black");
}

function onMessage(evt) {
	writeToScreen("RESPONSE: " + evt.data, "blue");
}

function onError(evt) {
	writeToScreen("ERROR: " + evt.data, "red");
}

function doSend(message) {
	writeToScreen("SENT: " + message, "black");
	websocket.send(message);
}

function sendmsg(ev) {
	ev.preventDefault();
	doSend(msg.value);
}

function writeToScreen(message, color) {
	var p = document.createElement("p");
	p.style.wordWrap = "break-word";
	p.style.color = color;
	p.textContent = message;
	output.appendChild(p);
}


</script>
<h2>WebSocket Test</h2>
<form id=theform>
<input id=msg>
</form>
<div id=output></div>
<script>
window.addEventListener("load", init);
</script>
</html>
`
