import { Controller } from '@hotwired/stimulus'

export default class extends Controller<HTMLElement> {
    static targets = ['input', 'output']
    static values = { query: String }

    declare queryValue: string
    declare readonly inputTarget: HTMLInputElement
    declare readonly outputTarget: HTMLInputElement

    connect() {
        this.load()
    }

    queryValueChanged() {
        this.load(this.queryValue)
    }

    filter() {
        const query = this.inputTarget.value.trim()
        this.queryValue = query
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