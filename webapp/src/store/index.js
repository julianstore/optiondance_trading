import { createStore } from 'vuex';
import {user} from "@/store/module/user";
import {option} from "@/store/module/option";
import createPersistedState from "vuex-persistedstate";
import getters from "@/store/getters";

const vuexLocal = createPersistedState()

export const store = createStore({
	modules: {
		user,option
	},
	plugins: [vuexLocal],
	getters
});