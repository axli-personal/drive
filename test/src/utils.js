const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

export function randomString(length) {
    let result = '';
    while (length--) {
        result += charset[Math.floor(charset.length * Math.random())];
    }
    return result;
}

export function statusCheck(response) {
    return response.status >= 200 && response.status < 300;
}
