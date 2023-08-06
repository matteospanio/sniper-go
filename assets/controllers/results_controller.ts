import { Controller } from '@hotwired/stimulus'
import { ResultSummary } from '../ResultSummary'

interface Response {
    results: ResultSummary[]
}

export default class extends Controller {
    connect() {
        this.load()
    }

    load() {

        const element = document.createElement('ul')
        element.classList.add('list-group')
        fetch('/results')
            .then(response => response.json())
            .then( (data: Response) => {
                data.results.forEach(result => {
                    let li = document.createElement('li')
                    li.innerHTML = `${result.host} - ${result.score}`
                    element.appendChild(li)
                })
                this.element.appendChild(element)
            })
            .catch(error => console.error(error))
    }
}