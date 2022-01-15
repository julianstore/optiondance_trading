import { createApp } from 'vue'
import '@/assets/css/index.scss'
import '@/assets/css/common.scss'
import '@/assets/css/option.scss'
import '@/assets/css/home.scss'
import '@/assets/css/article.scss'
import App from './App.vue'
import router from '@/router/router'
import {store} from '@/store'
import Toast from "@/components/common/toast";
import Loading from "@/components/common/loading/loading";
import PaymentLoading from "@/components/common/paymentLoading";
import TabLoading from "@/components/common/TabLoading";
import VueClipboard from 'vue3-clipboard'
import errHandler from "@/errHandler";
import * as Sentry from "@sentry/vue";
import { Integrations } from "@sentry/tracing";
import {createI18n } from "vue-i18n";
import {messages} from '@/i18n/i18n'

const app = createApp(App)


app.config.errorHandler = errHandler


const i18n = createI18n({
	locale: 'zhCn',
	// locale: 'zhTw',
	messages: messages,
})

app.use(i18n)
app.use(router)
app.use(store)
app.use(Toast)
app.use(Loading)
app.use(PaymentLoading)
app.use(TabLoading)
app.use(VueClipboard,{
	autoSetContainer: true,
		appendToBody: true,
})

let sentryDsnMap = {
	development:'',
	beta: 'https://b2006a48095d4a1a81e8b79712937019@o517640.ingest.sentry.io/5920355',
	prod: 'https://e89d528767e041238f362f8de605914c@o517640.ingest.sentry.io/5730434',
}

let dsn = sentryDsnMap[import.meta.env.MODE]
if (dsn) {
	Sentry.init({
		app,
		dsn: dsn,
		integrations: [
			new Integrations.BrowserTracing({
				routingInstrumentation: Sentry.vueRouterInstrumentation(router),
				// tracingOrigins: ["localhost", "my-site-url.com", /^\//],
			}),
		],
		tracesSampleRate: 1.0,
	});
}

app.mount('#app')