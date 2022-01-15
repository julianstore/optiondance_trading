import {userSettings} from "@/api/user";

export const user = {
	namespaced: true,
	state: {
		userInfo: {
			id:"",
			uuid: "",
			nickName: "",
			headerImg: "",
			authority: "",
		},
		waitlist: {
			email: '',
			inviterWid: 0,
			type: 0,
			wid:0,
		},
		token: "",
		expiresAt: "",
		loggedIn:false,
		mixinToken:'',
		mixinUser:{},
		registerInfo:{},
		inviterWid: 0,
		platformType: 0,
		settings:{}
	},
	mutations: {
		setWaitlist(state,w) {
			state.waitlist = w
		},
		setSettings(state,s) {
			state.settings = s
		},
		setPlatformType(state, type) {
			state.platformType = type
		},
		setInviterMid(state, inviterMid) {
			state.inviterMid = inviterMid
		},
		setUserInfo(state, userInfo) {
			state.userInfo = userInfo
		},
		setToken(state, token) {
			state.token = token
		},
		setExpiresAt(state, expiresAt) {
			state.expiresAt = expiresAt
		},
		setLoggedIn(state, loggedIn) {
			state.loggedIn = loggedIn
		},
		LoginOut(state,type) {
			state.userInfo = {}
			state.token = ""
			state.loggedIn = false
			state.expiresAt = ""
			let path = type === 'bithub' ? '/login' : '/login'
			router.push(path)
			sessionStorage.clear()
			// window.location.reload()
		},
		unauthorized(state) {
			state.userInfo = {}
			state.token = ""
			state.loggedIn = false
		},
		ResetUserInfo(state, userInfo = {}) {
			state.userInfo = {...state.userInfo,
				...userInfo
			}
		},
		setMixinToken(state, mixinToken = {}) {
			state.mixinToken = mixinToken
		},
		saveRegisterInfo(state, info = {}) {
			state.registerInfo = info
		},
		setInstaller(state, info = []) {
			state.installer = info
		},
		setMixinUser(state,mixinUser) {
			state.mixinUser = mixinUser
		}
	},
	actions: {
		async LoginOut({ commit },info) {
			commit("LoginOut",info.type)
		},
		async SetSettings({ commit },settings) {
			commit("setSettings",settings)
		},
		async SyncSettings({ commit,dispatch }) {
			let res = await userSettings();
			commit("setSettings",res.data)
			dispatch('option/SetDeliveryType', res.data.delivery_type ===1 ? 'PHYSICAL':'CASH',{root:true})
		},
		async MixinLogin({ commit },info) {
			commit('setMixinToken',info.token)
		},
		async SetMixinUser({ commit },user) {
			commit('setMixinUser',user)
		},
		async SaveRegisterInfo({commit},info){
			commit('saveRegisterInfo',info)
		},
		async SetInviterWid({commit},Wid) {
			commit('setInviterWid',Wid)
		},
		async SetPlatformType({commit},mid) {
			commit('setPlatformType',mid)
		},
		async SetWaitlist({commit},w) {
			commit('setWaitlist',w)
		},
		async SetDAppToken({commit},w) {
			commit('setToken',w)
		},
	},
	getters: {
		userInfo(state) {
			return state.userInfo
		},
		token(state) {
			return state.token
		},
		loggedIn(state) {
			return state.loggedIn
		},
		expiresAt(state) {
			return state.expiresAt
		},
		mixinToken(state) {
			return state.mixinToken
		},
		registerInfo(state) {
			return state.registerInfo
		},
		installer(state){
			return state.installer
		},
		mixinUser(state){
			return state.mixinUser
		},
		inviterWid(state){
			return state.inviterWid
		},
		platformType(state){
			return state.platformType
		},
		waitlist(state){
			return state.waitlist
		},
		settings(state){
			return state.settings
		}
	}
}
