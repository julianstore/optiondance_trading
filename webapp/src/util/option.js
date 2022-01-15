const zhCnSellString = "卖出"
const zhCnSellStr = "卖"
const zhCnBuyString = "买入"
const zhCnBuyStr = "买"

const zhTwSellString = "賣出"
const zhTwSellStr = "賣"
const zhTwBuyString = "買入"
const zhTwBuyStr = "買"


export function SideString(side,optionType,locale) {
    let sideStr = ''
    if (locale === 'zhCn') {
        if (side === 'BID' && optionType === 'PUT') {
            sideStr = '卖出'
        }else if (side === 'ASK' && optionType === 'PUT') {
            sideStr = '买入'
        }else if (side === 'BID' && optionType === 'CALL') {
            sideStr = '买入'
        }else if (side === 'ASK' && optionType === 'CALL') {
            sideStr = '卖出'
        }
    }else if (locale === 'zhTw') {
        if (side === 'BID' && optionType === 'PUT') {
            sideStr = zhTwSellString
        }else if (side === 'ASK' && optionType === 'PUT') {
            sideStr = zhTwBuyString
        }else if (side === 'BID' && optionType === 'CALL') {
            sideStr = zhTwBuyString
        }else if (side === 'ASK' && optionType === 'CALL') {
            sideStr = zhTwSellString
        }
    }
    return sideStr
}

export function SideStr(side,optionType,locale) {
    let sideStr = ''
    if (locale === 'zhCn') {
        if (side === 'BID' && optionType === 'PUT') {
            sideStr = '卖'
        }else if (side === 'ASK' && optionType === 'PUT') {
            sideStr = '买'
        }else if (side === 'BID' && optionType === 'CALL') {
            sideStr = '买'
        }else if (side === 'ASK' && optionType === 'CALL') {
            sideStr = '卖'
        }
    }else if(locale === 'zhTw'){
        if (side === 'BID' && optionType === 'PUT') {
            sideStr = zhTwSellString
        }else if (side === 'ASK' && optionType === 'PUT') {
            sideStr = zhTwBuyString
        }else if (side === 'BID' && optionType === 'CALL') {
            sideStr = zhTwBuyString
        }else if (side === 'ASK' && optionType === 'CALL') {
            sideStr = zhTwSellString
        }
    }
    return sideStr
}

export function optionFundsSign (side,optionType) {
    if ((optionType==='PUT' && side === 'ASK') || (optionType==='CALL' && side === 'ASK') ) {
        return '+'
    }
    if ( (optionType==='PUT' && side === 'BID') || (optionType==='CALL' && side === 'BID')  ) {
        return ''
    }
}