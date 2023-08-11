import { Controller } from '@hotwired/stimulus'
import { decode } from 'he';
import { websocket } from '../ws';
import ResultsController from './results_controller';
import { parseFormatAnsi } from '../utils';

export default class extends Controller {
    static targets = ['input', 'output', 'btn']
    static outlets = ['results']
    
    declare readonly btnTarget: HTMLButtonElement
    declare readonly inputTarget: HTMLInputElement
    declare readonly outputTarget: HTMLInputElement
    declare readonly resultsOutlet: ResultsController

    connect() {
        websocket.onmessage = (msg) => {
            const data = msg.data
            console.log(data)
            if (data === "[DONE]") {
                this.inputTarget.value = ''
                this.resultsOutlet.load()
                this.btnTarget.disabled = false
            } else {
                this.writeOutput(decode(data))
            }
        }
    }

    writeOutput(msg: string) {
        let lines = this.outputTarget.innerText.split('\n')
        if (lines.length === 800) {
            lines.shift()
            this.outputTarget.innerText = lines.join('\n')
        }
        this.outputTarget.innerHTML += `${parseFormatAnsi(msg)}\n`
        this.outputTarget.scrollTop = this.outputTarget.scrollHeight
    }

    sendInput() {
        const data = this.inputTarget.value
        this.btnTarget.disabled = true
        this.outputTarget.innerText += `\n~ ${data}\n`
        websocket.send(data)
    }
}