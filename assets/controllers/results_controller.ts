import { Controller } from '@hotwired/stimulus'

export default class extends Controller<HTMLElement> {
    connect() {
        this.load()
    }

    load() {
        fetch('/results')
            .then(response => response.text())
            .then(data => {
                this.element.innerHTML = data
            })
            .catch(error => console.error(error))
    }
}