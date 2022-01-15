import $axios from "@/api/index";

export function userMe() {
    return  $axios.get('/v1/user/me');
}

export function userSettings() {
    return  $axios.get('/v1/user-settings');
}