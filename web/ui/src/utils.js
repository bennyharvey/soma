const appendLeadingZero = (n) => {
    if(n <= 9){
        return '0' + n;
    }
    return n
}

export const formatDateTime = (dt) => {
    return dt.getFullYear()
        + '-' + appendLeadingZero(dt.getMonth() + 1)
        + '-' + appendLeadingZero(dt.getDate())
        + ' ' + appendLeadingZero(dt.getHours())
        + ':' + appendLeadingZero(dt.getMinutes())
        + ':' + appendLeadingZero(dt.getSeconds())
}