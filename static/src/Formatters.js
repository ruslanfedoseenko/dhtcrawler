
export default {
  formatBytes(bytes, decimals) {
    if (bytes === 0) return '0 Byte'
    const k = 1024 // or 1024 for binary
    let dm = decimals || 3
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
    let i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
  },
  nFormatter(num, digits) {
    let si = [
      {value: 1, symbol: ''},
      {value: 1E3, symbol: 'k'},
      {value: 1E6, symbol: 'M'},
      {value: 1E9, symbol: 'G'},
      {value: 1E12, symbol: 'T'},
      {value: 1E15, symbol: 'P'},
      {value: 1E18, symbol: 'E'}
    ]
    let rx = /\.0+$|(\.[0-9]*[1-9])0+$/
    let i
    for (i = si.length - 1; i > 0; i--) {
      if (num >= si[i].value) {
        break
      }
    }
    return (num / si[i].value).toFixed(digits).replace(rx, '$1') + ' ' + si[i].symbol
  }
}
