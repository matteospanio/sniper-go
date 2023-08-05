import { Controller } from '@hotwired/stimulus'
import { websocket } from '../ws';

export default class extends Controller {
    static targets = ['input', 'output']

    declare readonly hasInputTarget: boolean
    declare readonly inputTarget: HTMLInputElement
    declare readonly inputTargets: HTMLInputElement[]

    declare readonly hasOutputTarget: boolean
    declare readonly outputTarget: HTMLInputElement
    declare readonly outputTargets: HTMLInputElement[]

    connect() {
        console.log('socket controller');
    }

    sendInput() {
        let data = this.inputTarget.value
        websocket.send(data)
    }
}