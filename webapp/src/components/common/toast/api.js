import mount from "@/util/mount-component";
import Toast from "@/components/common/toast/Toast.vue";

const ToastApi = (globalOptions = {}) => {
	return {
		show(message, options = {}) {
			let localOptions = { message, ...options }
			const c = mount(Toast, {
				props: { ...globalOptions, ...localOptions }
			})
			return c
		},
		// clear() {
		// 	eventBus.$emit('toast-clear')
		// },
		success(message, options = {}) {
			options.type = 'success'
			return this.show(message, options)
		},
		error(message, options = {}) {
			options.type = 'error'
			return this.show(message, options)
		},
		info(message, options = {}) {
			options.type = 'info'
			return this.show(message, options)
		},
		warning(message, options = {}) {
			options.type = 'warning'
			return this.show(message, options)
		}
	}
}

export default ToastApi
