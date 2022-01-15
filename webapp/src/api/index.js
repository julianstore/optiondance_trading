import axios from "axios";
import { store } from "@/store";
import router from "@/router/router";
import { useTabLoading } from "@/components/common/TabLoading/index";

const $axios = axios.create({
  // Set timeout
  timeout: 300000,
  // The base url, will automatically add a pre-link in the request url
  baseURL: import.meta.env.VITE_BASE_API + "/api",
});
// Vue.prototype.$http = axios // Concurrent request
// Add request status in the global request and response interceptor
let loading = null;

// Request interceptor
let NoLoadingUris = [
  "https://mixin-api.zeromesh.net/me",
  "/v1/mtg/payments",
  "/v1/order-request",
  "v1/order-trace",
  "https://api.coingecko.com",
  "/v1/market/instrument/",
  "/v1/exercise-request/",
  "/v1/exercise-cash/",
  "/v1/position-status/",
  "/v1/market/price/",
  "/v1/orders",
  "/v1/positions",
  "/v1/order",
  "/v1/position",
  "/v2/position",
];

let NoCommonHeaderUris = ["https://mixin-api.zeromesh.net/me"];

$axios.interceptors.request.use(
  (config) => {
    let isLoading = true;
    let addCommonAuthHeader = true;
    for (let i in NoLoadingUris) {
      if (config.url.indexOf(NoLoadingUris[i]) !== -1) {
        isLoading = false;
        break;
      }
    }
    for (let i in NoCommonHeaderUris) {
      if (config.url.indexOf(NoCommonHeaderUris[i]) !== -1) {
        addCommonAuthHeader = false;
        break;
      }
    }
    if (isLoading) {
      loading = useTabLoading();
      loading.show();
    }
    const token = store.getters.token;
    if (token && addCommonAuthHeader) {
      config.headers.Authorization = token;
    }
    return config;
  },
  (error) => {
    loading.hide();
    return Promise.reject(error);
  }
);
// Response interceptor
$axios.interceptors.response.use(
  (response) => {
    if (loading) {
      loading.hide();
    }
    const code = response.status;
    if ((code >= 200 && code < 300) || code === 304) {
      return Promise.resolve(response.data);
    } else {
      return Promise.reject(response);
    }
  },
  (error) => {
    if (loading) {
      loading.hide();
    }
    console.log(error);
    if (error.response) {
      switch (error.response.status) {
        case 401:
          store.commit("user/unauthorized");
          router.push({ path: "/auth", query: { action: "auth" } });
          break;
        case 404:
          // Message.error('Network request does not exist')
          break;
        default:
        // Message.error(error.response.data.toast)
      }
    } else {
      // Request timed out or network problem
      if (error.message.includes("timeout")) {
        // Message.error('Request timed out! Please check if the network is normal')
      } else {
        // Message.error('The request failed, please check if the network is connected')
      }
    }
    return Promise.reject(error);
  }
);

// get, post request method
// export default {
// post(url, data) {
// 	return $axios({
// 		method: 'post',
// 		url,
// 		data: JSON.stringify(data),
// 		headers: {
// 			'Content-Type': 'application/json'
// 		}
// 	})
// },
//
// postRaw(url, data) {
// 	return $axios({
// 		method: 'post',
// 		url,
// 		data: data,
// 		headers: {
// 			'Content-Type': 'application/json'
// 		}
// 	})
// },
//
// postForm(url,data) {
// 	return $axios({
// 		method: 'post',
// 		url,
// 		data: data,
// 		transformRequest: [function (data) {
// 			// Do whatever you want to transform the data
// 			let ret = ''
// 			for (let it in data) {
// 				ret += encodeURIComponent(it) + '=' + encodeURIComponent(data[it]) + '&'
// 			}
// 			if (ret.charAt(ret.length - 1) === '&') {
// 				ret = ret.substr(0, ret.length - 1);
// 			}
// 			return ret;
// 		}],
// 		headers: {
// 			'Content-Type': 'application/x-www-form-urlencoded'
// 		}
// 	})
// },
// postMultipartFormData(url,data) {
// 	return $axios({
// 		method: 'post',
// 		url,
// 		data: data,
// 		headers: {
// 			'Content-Type': 'multipart/form-data'
// 		}
// 	})
// },
// get(url, params) {
// 	return $axios({
// 		method: 'get',
// 		url,
// 		params
// 	})
// },
// delete(url, params) {
// 	return $axios({
// 		method: 'delete',
// 		url,
// 		params
// 	})
// },
//$axios
// }

export default $axios;
