import { Controller } from '@hotwired/stimulus'
import DataTable from 'datatables.net-bs5';
import 'datatables.net-buttons-bs5';
import 'datatables.net-buttons/js/buttons.html5.mjs';
import 'datatables.net-buttons/js/buttons.print.mjs';
import 'datatables.net-buttons/js/buttons.flash';
import 'jszip';
import 'pdfmake';
import { severityToBsClass } from '../ResultSummary';

export default class extends Controller<HTMLElement> {
    static targets = ['output']
    static values = { type: String }

    declare readonly outputTarget: HTMLTableElement
    declare readonly typeValue: string

    connect() {
        new DataTable(this.outputTarget, {
            dom: "<'card-header d-flex justify-content-between align-items-center'Bfr><'card-body't><'card-footer d-flex justify-content-between align-items-center' ip>",
            buttons: [
                'excel', 'csv', 'print', 'copy'
            ],
            columnDefs: [
                {
                    targets: 0,
                    render: (data) => {
                        return this.typeValue === 'vuln'
                            ? `<span class="badge bg-${severityToBsClass(data)}">${data}</span>`
                            : data
                    }
                },
            ],
        })
    }
}