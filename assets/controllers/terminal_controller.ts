import { Controller } from '@hotwired/stimulus'
import { cleanStringColors } from '../utils';
import { getWebSocket } from '../ws';
import ResultsController from './results_controller';

export default class extends Controller {
    static targets = ['input', 'output']
    static outlets = ['results']
    
    declare readonly inputTarget: HTMLInputElement
    declare readonly outputTarget: HTMLInputElement
    declare readonly resultsOutlet: ResultsController

    connect() {
        const websocket = getWebSocket()

        websocket.onmessage = (msg) => {
            const data = msg.data
            if (data === '[DONE]') {
                websocket.send('exit')
                websocket.close()
                this.inputTarget.value = ''
                this.resultsOutlet.load()
            } else {
                this.writeOutput(cleanStringColors(data))
            }
        }
    }

    writeOutput(msg: string) {
        let lines = this.outputTarget.innerText.split('\n')
        if (lines.length === 800) {
            lines.shift()
            this.outputTarget.innerText = lines.join('\n')
        }
        this.outputTarget.innerText += `${msg}\n`
    }

    sendInput() {
        const data = this.inputTarget.value
        this.outputTarget.innerText += `\n~ ${data}\n`
        getWebSocket().send(data)
    }
}