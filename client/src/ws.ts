import { HOST, PORT } from './constants'

const ADDRESS = `ws:/${HOST}:${PORT}/ws`

export const websocket = new WebSocket(ADDRESS);