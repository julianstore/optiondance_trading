import $axios from "@/api/index";

export function getStrikePrices(side,optionType,deliveryType,quoteCurrency,baseCurrency) {
    return  $axios.get('/v1/market/strike-prices',{
        params:{
            side:side,
            optionType:optionType,
            deliveryType: deliveryType,
            quoteCurrency:quoteCurrency,
            baseCurrency:baseCurrency
        }
    });
}



export function listExpiryDatesByPrice(strikePrice,side,optionType,deliveryType,quoteCurrency,baseCurrency) {
    return  $axios.get(`/v1/market/price/${strikePrice}`,{
        params:{
            side:side,
            optionType:optionType,
            deliveryType: deliveryType,
            quoteCurrency:quoteCurrency,
            baseCurrency:baseCurrency
        }
    });
}