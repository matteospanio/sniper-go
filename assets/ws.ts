export const websocket = new WebSocket("ws://localhost:8080/ws");

websocket.onmessage = function (msg) {
    console.log(msg)
}
