import { Controller } from '@hotwired/stimulus'

export default class extends Controller<HTMLElement> {
    static targets = ['input', 'btn']

    declare readonly inputTarget: HTMLInputElement
    declare readonly btnTarget: HTMLInputElement

    toggle() {
        if (this.inputTarget.classList.contains('d-none')) {
            this.inputTarget.classList.remove('d-none')
            this.btnTarget.classList.add('d-none')
            this.inputTarget.focus()
        } else {
            this.inputTarget.classList.add('d-none')
            this.btnTarget.classList.remove('d-none')
        }
    }
}