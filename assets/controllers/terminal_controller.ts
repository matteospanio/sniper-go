import { Controller } from '@hotwired/stimulus'
import { websocket } from '../ws';
import ListController from './list_controller';
import { parseAnsi } from '../utils';

export default class extends Controller {
    static targets = ['input', 'output']
    static outlets = ['list']

    declare readonly inputTarget: HTMLInputElement
    declare readonly outputTarget: HTMLInputElement
    declare readonly listOutlets: ListController[]


    connect() {
        websocket.onmessage = (msg) => {
            const data = msg.data
            if (data === "[DONE]") {
                this.inputTarget.value = ''
                this.loadOutlets()
            } else {
                console.log(parseAnsi(data))
            }
        }
    }

    writeOutput(msg: string) {
        let lines = this.outputTarget.innerText.split('\n')
        if (lines.length === 200) {
            lines.shift()
            this.outputTarget.innerText = lines.join('\n')
        }
        this.outputTarget.innerHTML += `${parseAnsi(msg)}\n`
        this.outputTarget.scrollTop = this.outputTarget.scrollHeight
    }

    sendInput() {
        const data = this.inputTarget.value
        // this.outputTarget.innerText += `\n~ ${data}\n`
        websocket.send(data)
        setTimeout(() => this.loadOutlets(), 5000)
    }

    private loadOutlets() {
        this.listOutlets.forEach((outlet) => outlet.load())
    }
}