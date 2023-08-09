const HOST = "0.0.0.0"
const PORT = 8080

let websocket = new WebSocket(`ws:/${HOST}/:${PORT}/ws`);

export function getWebSocket(): WebSocket {
    if (websocket.readyState === WebSocket.CLOSED)
        websocket = new WebSocket(`ws:/${HOST}/:${PORT}/ws`);
    return websocket;
}