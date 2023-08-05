import { Controller } from '@hotwired/stimulus'

export default class extends Controller<HTMLElement> {
    static values = {
        theme: String
    }

    declare themeValue: string
    declare readonly hasThemeValue: boolean

    connect() {
        this.element.dataset.bsTheme=this.themeValue
    }

    themeValueChanged() {
        this.element.dataset.bsTheme = this.themeValue
    }

    switchTheme() {
        if (this.themeValue === 'light') {
            this.themeValue = 'dark'
        } else {
            this.themeValue = 'light'
        }
    }
}