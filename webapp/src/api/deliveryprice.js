import $axios from "@/api/index";

export function listDeliveryPrices(baseCurrency) {
    return  $axios.get(`/v1/delivery-prices`,{
        params:{
            asset: baseCurrency
        }
    });
}