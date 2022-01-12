const f_replaceAll = require("replaceall");

export function replaceAll(a: string, b: string, text: string): string {
    return f_replaceAll(a, b, text);
}

export function splitAll(text: string, splitters: string[]): string[] {
    if (splitters?.length) {
        const parts = text.split(splitters[0]);
        if (splitters.length > 1) {
            splitters = splitters.slice(1);
            const subParts: string[] = [];
            parts.forEach(part => subParts.push(...splitAll(part, splitters)));
            return subParts;
        } else {
            return parts;
        }
    } else {
        return [text];
    }
}