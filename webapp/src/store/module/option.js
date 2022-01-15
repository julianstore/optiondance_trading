export const option = {
	namespaced: true,
	state: {
		optionType: '',
		optionSide: '',
		deliveryType: '',
	},
	mutations: {
		setOptionType(state,w) {
			state.optionType = w
		},
		setOptionSide(state,w) {
			state.optionSide = w
		},
		setDeliveryType(state,w) {
			state.deliveryType = w
		}
	},
	actions: {
		async SetOptionType({commit}, w) {
			commit('setOptionType', w)
		},
		async SetOptionSide({commit}, w) {
			commit('setOptionSide', w)
		},
		async SetDeliveryType({commit}, w) {
			commit('setDeliveryType', w)
		},
	},
	getters: {
		optionType(state) {
			return state.optionType
		},
		optionSide(state) {
			return state.optionSide
		},
		optionSideType(state) {
			return state.optionSide + '_' + state.optionType
		},
		deliveryType(state) {
			return state.deliveryType
		},
	}
}
