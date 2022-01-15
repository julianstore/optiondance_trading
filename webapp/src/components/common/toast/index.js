import Toast from "@/components/common/toast/Toast.vue";
import ToastApi from './api'


const Plugin = (app, options = {}) => {
	let methods = ToastApi(options)
	app.$toast = methods
	app.config.globalProperties.$toast = methods
}

Toast.install = Plugin

export default Toast
