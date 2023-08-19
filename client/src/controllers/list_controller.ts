import {Controller} from '@hotwired/stimulus'

export default class extends Controller<HTMLElement> {
    static targets = ['input', 'output']
    static values = { query: String, url: String}

    declare queryValue: string
    declare urlValue: string
    declare readonly inputTarget: HTMLInputElement
    declare readonly outputTarget: HTMLInputElement

    queryValueChanged() {
        this.load(this.queryValue)
    }

    filter() {
        this.queryValue = this.inputTarget.value.trim()
    }

    load(query: string = '') {
        fetch(`${this.urlValue}?query=${query}`)
            .then(response => response.text())
            .then(data => {
                this.outputTarget.innerHTML = data
            })
            .catch(error => console.error(error))
    }
}