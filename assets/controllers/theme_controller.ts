import { Controller } from '@hotwired/stimulus'

export default class extends Controller<HTMLElement> {
    static values = {
        theme: String,
        icon: String
    }

    static targets = ['icon']

    declare themeValue: string
    declare iconValue: string
    declare readonly iconTarget: HTMLInputElement

    connect() {
        this.element.dataset.bsTheme=this.themeValue
        this.iconTarget.classList.add(this.iconValue)
    }

    themeValueChanged() {
        this.element.dataset.bsTheme = this.themeValue
    }

    iconValueChanged() {
        this.iconTarget.classList.remove(this.getOldIcon())
        this.iconTarget.classList.add(this.iconValue)
    }

    switchTheme() {
        this.themeValue = this.themeValue === 'light' ? 'dark' : 'light'
        this.iconValue = this.themeValue === 'light' ? 'fa-moon' : 'fa-sun'
    }

    private getOldIcon() {
        return this.iconValue === 'light' ? 'fa-sun' : 'fa-moon'
    }
}