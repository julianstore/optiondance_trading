import { createRouter, createWebHistory } from "vue-router";
import NotFound from "@/views/NotFound.vue";
import Test from "@/views/Test.vue";
import WaitList from "@/views/WaitList.vue";
import Auth from "@/views/Auth.vue";
import WaitListEmail from "@/views/WaitListEmail.vue";
import JoinSuccess from "@/views/JoinSuccess.vue";
import Layout from "@/components/layout/Layout.vue";
import Faq from "@/views/Faq.vue";
import Home from "@/views/Home.vue";
import Options from "@/views/option/Options.vue";
import Order from "@/views/option/Order.vue";
import Position from "@/views/option/Position.vue";
import Settings from "@/views/Settings.vue";
import OptionDanceIntro from "@/views/faq/OptionDanceIntro.vue";
import SellPutIntro from "@/views/faq/SellPutIntro.vue";
import Income from "@/views/Income.vue";
import IndexPrices from "@/views/IndexPrices.vue";
import ShareCard from "@/views/SharePage.vue";

const routes = [
  { path: "/", component: WaitList, meta: { title: "OptionDance" } },
  {
    path: "/test",
    component: Test,
    meta: { title: "test" },
  },
  {
    name: "share",
    path: "/share",
    component: ShareCard,
    meta: { title: "test" },
  },
  {
    path: "/home",
    component: Home,
    meta: { title: "OptionDance - 首页" },
  },
  {
    path: "/waitlist",
    component: WaitList,
    meta: { title: "OptionDance - WaitList" },
  },
  {
    path: "/faq",
    component: Faq,
    meta: { title: "OptionDance - 常见问题" },
  },
  {
    name: "income",
    path: "/income",
    component: Income,
    meta: { title: "OptionDance - 收益分析" },
  },
  {
    path: "/faq",
    component: Layout,
    meta: { title: "OptionDance - 常见问题" },
    children: [
      {
        name: "when-to-use-sell-put",
        path: "when-to-use-sell-put",
        component: SellPutIntro,
        meta: { title: "OptionDance - 什么时候使用优买？" },
      },
      {
        name: "product-intro",
        path: "what-is-option-dance",
        component: OptionDanceIntro,
        meta: { title: "OptionDance - 产品介绍" },
      },
    ],
  },
  {
    name: "subscribe-email",
    path: "/subscribe-email",
    component: WaitListEmail,
    meta: { title: "OptionDance - 常见问题" },
  },
  {
    name: "join-success",
    path: "/join-success",
    component: JoinSuccess,
    meta: { title: "OptionDance - 加入成功" },
  },
  {
    path: "/auth",
    component: Auth,
    meta: { title: "Auth" },
  },

  {
    name: "options",
    path: "/options",
    component: Options,
    meta: { title: "OptionDance - 交易详情" },
  },
  {
    path: "/order/:id",
    component: Order,
    meta: { title: "OptionDance - 订单详情" },
  },
  {
    path: "/index-prices/:asset",
    component: IndexPrices,
    meta: { title: "OptionDance - 报价日志" },
  },
  {
    path: "/settings",
    component: Settings,
    meta: { title: "OptionDance - 设置" },
  },
  {
    path: "/position/:id",
    component: Position,
    meta: { title: "OptionDance - 持仓详情" },
  },
  { path: "/:path(.*)", component: NotFound },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  /* The routing changes modify the page title */
  if (to.meta.title) {
    document.title = to.meta.title;
  }
  next();
});

export default router;
