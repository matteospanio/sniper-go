const HOST = "0.0.0.0"
const PORT = 8080
const ADDRESS = `ws:/${HOST}:${PORT}/ws`

let websocket = new WebSocket(ADDRESS);

export function getWebSocket(): WebSocket {
    if (websocket.readyState === WebSocket.CLOSED)
        websocket = new WebSocket(ADDRESS);
    return websocket;
}