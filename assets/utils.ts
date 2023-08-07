const OKGRAY   = /\[90m/g    
const OKRED    = /\[91m/g
const OKGREEN  = /\[92m/g
const OKORANGE = /\[93m/g
const OKBLUE   = /\[94m/g
const OKPURPLE = /\[95m/g
const OKCYAN   = /\[96m/g
const OKWHITE  = /\[97m/g
const RESET    = /\[0m/g

export const cleanStringColors = (str: string) => {
    return str.replace(OKGRAY, '')
            .replace(OKRED, '')
            .replace(OKGREEN, '')
            .replace(OKORANGE, '')
            .replace(OKBLUE, '')
            .replace(OKPURPLE, '')
            .replace(OKCYAN, '')
            .replace(OKWHITE, '')
            .replace(RESET, '')
}

export const translateColors = (str: string) => {
    return str.replace(OKGRAY, '<span class="text-secondary">')
            .replace(OKRED, '<span class="text-danger">')
            .replace(OKGREEN, '<span class="text-success">')
            .replace(OKORANGE, '<span class="text-warning">')
            .replace(OKBLUE, '<span class="text-primary">')
            .replace(OKPURPLE, '<span class="text-primary-emphasis">')
            .replace(OKCYAN, '<span class="text-primary">')
            .replace(OKWHITE, '<span class="text-white">')
            .replace(RESET, '</span>')
}