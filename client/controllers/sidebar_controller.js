import Logo from '../images/MPAI-logo.png'
import { Controller } from '@hotwired/stimulus'

export default class extends Controller {
    static targets = ["logo", 'link']

    connect() {
        this.logoTarget.innerHTML = `<img class="img-fluid" src="${Logo}" />`
    }

    onLinkHover(e) {
        this.linkTargets.forEach(element => {
            element.classList.replace('active', 'link-dark')
        });

        e.currentTarget.classList.replace('link-dark', 'active')
    }
}