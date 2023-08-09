import { Controller } from '@hotwired/stimulus'

export default class extends Controller<HTMLElement> {
    static targets = ['input', 'output']

    declare readonly inputTarget: HTMLInputElement
    declare readonly outputTarget: HTMLInputElement

    connect() {
        this.load()
    }

    filter() {
        const query = this.inputTarget.value.trim()
        this.load(query)
    }

    load(query: string = '') {
        fetch(`/results?query=${query}`)
            .then(response => response.text())
            .then(data => {
                this.outputTarget.innerHTML = data
            })
            .catch(error => console.error(error))
    }
}