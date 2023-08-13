import { Controller } from "@hotwired/stimulus";
import { Chart } from "chart.js/auto";
import { getData } from "../utils";

export default class extends Controller<HTMLCanvasElement> {
    static values = { host: String }

    declare readonly hostValue: string

    connect(): void {
        this.renderChart()
    }

    renderChart(): void {
        getData(this.hostValue).then((data) => {
                new Chart(this.element, {
                    type: 'pie',
                    data: {
                        labels: [
                            'Critical',
                            'High',
                            'Medium',
                            'Low',
                            'Info'
                        ],
                        datasets: [{
                            label: `${this.hostValue} report`,
                            data: [data.critical, data.high, data.medium, data.low, data.info],
                        }]
                    },
                    options: {
                        maintainAspectRatio: false,
                    }
                })
            })
            .catch(err => console.log(err))
    }
}