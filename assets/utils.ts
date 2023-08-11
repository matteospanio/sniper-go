const COLORS = [
    /\[0m/g,     // reset
    /\[1m/g,     // bold
    /\[2m/g,     // dim
    /\[3m/g,     // italic
    /\[4m/g,     // underline
    /\[5m/g,     // blink
    /\[6m/g,     // blink fast
    /\[7m/g,     // reverse
    /\[30m/g,    // black
    /\[31m/g,    // red
    /\[32m/g,    // green
    /\[33m/g,    // yellow
    /\[90m/g,    // gray
    /\[91m/g,    // red
    /\[92m/g,    // green
    /\[93m/g,    // orange
    /\[94m/g,    // blue
    /\[95m/g,    // purple
    /\[96m/g,    // cyan
    /\[97m/g,    // white
]

export const cleanStringColors = (str: string) => {
    for (const color of COLORS) {
        str = str.replace(color, '')
    }
    return str
}
