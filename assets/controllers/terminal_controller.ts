import { Controller } from '@hotwired/stimulus'
import { cleanStringColors } from '../utils';
import { websocket } from '../ws';


export default class extends Controller {
    static targets = ['input', 'output']
    
    declare readonly inputTarget: HTMLInputElement
    declare readonly outputTarget: HTMLInputElement

    connect() {
        websocket.onmessage = (msg) => {
            this.writeOutput(cleanStringColors(msg.data))
        }
    }

    writeOutput(msg: string) {
        this.outputTarget.innerText += `${msg}\n`
    }

    sendInput() {
        const data = this.inputTarget.value
        this.outputTarget.innerText += `\n~ ${data}\n`
        websocket.send(data)
    }
}