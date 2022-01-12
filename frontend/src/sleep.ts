export async function sleep(milliseconds: number = 0) {
    return new Promise(resolve => {
        setTimeout(resolve, milliseconds);
    });
}