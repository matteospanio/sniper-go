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

export const cleanStringColors = (str: string): string => {
    for (const color of COLORS) {
        str = str.replace(color, '')
    }
    return str
}

const ansiRegex = /\x1b\[(\d+)(;(\d+))*m/g;

// replace ansi codes with bootstrap classes in span tags
export const parseFormatAnsi = (str: string): string => {
    let match = ansiRegex.exec(str)
    while (match !== null) {
        const code = match[1]
        const style = match[3]
        if (code === '0') {
            str = str.replace(match[0], '</span>')
        } else {
            let className = ''
            switch (code) {
                case '1':
                    className = 'fw-bold'
                    break
                case '2':
                    className = 'fw-light'
                    break
                case '3':
                    className = 'fst-italic'
                    break
                case '4':
                    className = 'text-decoration-underline'
                    break
                case '5':
                    className = 'text-decoration-line-through'
                    break
                case '7':
                    className = 'text-decoration-underline'
                    break
                case '30':
                    className = 'text-dark'
                    break
                case '31':
                    className = 'text-danger'
                    break
                case '32':
                    className = 'text-success'
                    break
                case '33':
                    className = 'text-warning'
                    break
                case '90':
                    className = 'text-secondary'
                    break
                case '91':
                    className = 'text-danger'
                    break
                case '92':
                    className = 'text-success'
                    break
                case '93':
                    className = 'text-warning'
                    break
                case '94':
                    className = 'text-primary'
                    break
                case '95':
                    className = 'text-purple'
                    break
                case '96':
                    className = 'text-info'
                    break
                case '97':
                    className = 'text-white'
                    break
            }
            str = str.replace(match[0], `<span class="${className}">`)
        }
        match = ansiRegex.exec(str)
    }
    return str    
}