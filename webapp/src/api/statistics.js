import $axios from "@/api/index";


export function annualSellPutPremium() {
    return $axios.get(`/v1/statistics/annual-sell-put-premium`)
}

export function annualSellPutUnderlying() {
    return $axios.get(`/v1/statistics/annual-sell-put-underlying`)
}