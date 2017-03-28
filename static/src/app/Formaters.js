/**
 * Created by Ruslan on 09-Feb-17.
 */



export default class Formatter{
    static formatBytes(bytes, decimals) {
        if (bytes === 0) return '0 Byte';
        const k = 1024; // or 1024 for binary
        let dm = decimals || 3;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
        let i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    }
}