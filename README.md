
## Simple Websocket Chat in Golang

Пример реализации простого websocket-chat сервера на Go.
Собственно разбираемся с гоу-рутинами и каналами.
Требует пакет `golang.org/x/net/websocket`.

## Front-end для ваших тестов

```html
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<script src="http://code.jquery.com/jquery-2.2.0.min.js"></script>
<title>Simple Websocket Chat</title>
</head>
<body>

<hr>
<ul class="list" id="chat-list"></ul>
<hr>

<form class="form" id="chat-form">
	<input type="text" name="message" value="" maxlength="200" autocomplete="off" id="chat-input-message">
	<input type="submit" value="Go!" id="chat-button-submit">
</form>

<script type="text/javascript">
var chatInput = $("#chat-input-message"),
	chatList = $("#chat-list"),
	chatForm = $("#chat-form").get(0),
	chatSubmit = $("#chat-button-submit").get(0);

if('WebSocket' in window) {
	var ws = new WebSocket("ws://localhost:5213");

	chatForm.onsubmit = function(e) {
		e.preventDefault();
		chatSubmit.click();
		return false;
	};

	chatSubmit.onclick = function(e) {
		e.preventDefault();
		var msg = chatInput.val();
		var data = {
			body: msg
		};
		ws.send(JSON.stringify(data));
		chatInput.val("");
		return false;
	};

	ws.onopen = function() {
		$('<li>').text("Соединение установленно...").appendTo(chatList);
	};

	ws.onmessage = function(e) {
		try {
			var data = JSON.parse(e.data);
			$('<li>').text(data.body).appendTo(chatList);
		} catch(e) {
			$('<li>').text("Не удается разобрать ответа сервера").appendTo(chatList);
		}
	};

	ws.onerror = function(e) {
		if (e.data) {
			$('<li>').text(e.data).appendTo(chatList);
		} else {
			$('<li>').text("Ошибка...").appendTo(chatList);
		}
	};

	ws.onclose = function(e) {
		if (e.wasClean) {
			$('<li>').text("Соединение закрыто...").appendTo(chatList);
		} else {
			$('<li>').text("Обрыв соединения...").appendTo(chatList);
		}
	};
} else {
	alert("Браузер не поддерживает WebSocket");
}
</script>

</body>
</html>
```