const HOST = "0.0.0.0"
const PORT = 8080
const ADDRESS = `ws:/${HOST}:${PORT}/ws`

export const websocket = new WebSocket(ADDRESS);