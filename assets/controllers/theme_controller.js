import { Controller } from '@hotwired/stimulus'

export default class extends Controller {
    static values={
        theme: String
    }

    // connect() {
    //     this.element.dataset.bsTheme=this.themeValue
    // }

    themeValueChanged() {
        this.element.dataset.bsTheme=this.themeValue
    }

    onClick() {
        if (this.themeValue === 'light') {
            this.themeValue = 'dark'
        } else {
            this.themeValue = 'light'
        }
    }
}