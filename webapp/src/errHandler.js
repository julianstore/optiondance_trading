import ToastApi from "@/components/common/toast/api";
import { captureException } from '@sentry/browser';

export default function errHandler (err, vm, info) {
	let msg = err.message
	if (msg === 'Network Error') {
		let config = err.config
		let requestUrl = config.url.indexOf("http") === -1 ? config.baseURL + config.url : config.url
		ToastApi().error(`${msg}:${requestUrl}`);
	}else if (msg.indexOf('401')>-1){

	}else {
		ToastApi().error(`${msg}`);
	}
	console.log(import.meta.env.MODE)
	if (import.meta.env.MODE === 'prod' || import.meta.env.MODE === 'uat'){
		console.log("capture error")
		captureException(err,{
			tags:{
				env:import.meta.env.MODE
			}
		})
	}else {
		console.log(err)
	}
}
