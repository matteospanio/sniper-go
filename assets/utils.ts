import ansiRegex from 'ansi-regex';

export function parseAnsi(str: string): string {
    const regex = ansiRegex({onlyFirst: true})
    let match = str.match(regex)
    let content = str;
    while (match !== null) {
        const code = match[0]
        content = content.replace(code, "")
        match = str.match(regex)
    }

    return str
}


// replace ansi codes with bootstrap classes in span tags
export function parseFormatAnsi(str: string): string {
    let match = str.match(ansiRegex())
    while (match !== null) {
        const code = match[1]
        const style = match[3]
        if (code === '0') {
            str = str.replace(match[0], '</span>')
        } else if (code === '1K') {
            str = str.replace(match[0], '')
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
        //match = ansiRegex.exec(str)
    }
    return str    
}

const ansiToBootstrapClasses: Record<string, string> = {
    "\x1b[0m": "",             // Reset all styles
    "\x1b[1m": "fw-bold",      // Bold
    "\x1b[2m": "text-muted",   // Faint
    "\x1b[3m": "text-italic",  // Italic
    "\x1b[4m": "text-underline", // Underline
    "\x1b[5m": "text-blink",   // Blink
    "\x1b[7m": "bg-primary",   // Reverse colors (swap background and foreground)
    "\x1b[9m": "text-line-through", // Crossed-out
    "\x1b[30m": "text-dark",   // Black text
    "\x1b[31m": "text-danger", // Red text
    "\x1b[32m": "text-success", // Green text
    "\x1b[33m": "text-warning", // Yellow text
    "\x1b[34m": "text-info",   // Blue text
    "\x1b[35m": "text-purple", // Magenta text
    "\x1b[36m": "text-cyan",   // Cyan text
    "\x1b[37m": "text-light",  // White text
    "\x1b[40m": "bg-dark",     // Black background
    "\x1b[41m": "bg-danger",   // Red background
    "\x1b[42m": "bg-success",  // Green background
    "\x1b[43m": "bg-warning",  // Yellow background
    "\x1b[44m": "bg-info",     // Blue background
    "\x1b[45m": "bg-purple",   // Magenta background
    "\x1b[46m": "bg-cyan",     // Cyan background
    "\x1b[47m": "bg-light",    // White background
    "\x1b[90m": "text-muted",  // Light Black (Gray) text
    "\x1b[91m": "text-danger", // Light Red text
    "\x1b[92m": "text-success", // Light Green text
    "\x1b[93m": "text-warning", // Light Yellow text
    "\x1b[94m": "text-primary", // Light Blue text
    "\x1b[95m": "text-purple", // Light Magenta text
    "\x1b[96m": "text-info",   // Light Cyan text
    "\x1b[97m": "text-white",  // Light White text
};
