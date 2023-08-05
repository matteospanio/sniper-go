import { Application } from "@hotwired/stimulus"
import { definitionsFromContext } from "@hotwired/stimulus-webpack-helpers";

declare global {
    interface Window {
        Stimulus: Application
    }
}

window.Stimulus = Application.start()
// @ts-ignore
const context = require.context("./controllers", true, /\.ts$/);
window.Stimulus.load(definitionsFromContext(context));