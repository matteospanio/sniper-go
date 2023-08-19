export interface ResultSummary {
    critical: number
    high: number
    medium: number
    low: number
    info: number
    score: number
}

export function severityToBsClass(severity: string): string {
    const severityLower = severity.toLowerCase()
    switch (severityLower) {
        case "critical":
            return "dark"
        case "high":
            return "danger"
        case "medium":
            return "warning"
        case "low":
            return "success"
        case "info":
            return "info"
        default:
            return "secondary"
    }
}