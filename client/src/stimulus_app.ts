import { Application } from "@hotwired/stimulus"
import { definitionsFromContext } from "@hotwired/stimulus-webpack-helpers";

declare global {
    interface Window {
        Stimulus: Application
    }
}

const application = Application.start()
application.debug = false
window.Stimulus = application
// @ts-ignore
const context = require.context("./controllers", true, /\.ts$/);
window.Stimulus.load(definitionsFromContext(context));