<template>
  <div class="trade-info">
    <back-header back-to="-1"></back-header>
    <div ref="wrapper" class="wrapper">
      <div class="order-detail">
        <h1>{{ orderStatusTitle }}</h1>
        <div class="tip" v-if="showTip">
          <img src="@/assets/svg/tip.svg" />
          {{ $t("order.autoCancelTip") }}
        </div>
        <order-card
          :id="order.order_id"
          :order-num="order.order_num_string"
          :display-id="true"
          :remaining-amount="order.remaining_amount"
          :filled-amount="order.filled_amount"
          :remaining-funds="order.remaining_funds"
          :filled-funds="order.filled_funds"
          :expiration-date="order.expiration_date"
          :option-type="order.option_type"
          :side="order.side"
          :instrument-name="order.instrument_name"
          :order-status="order.order_status"
          :price="order.price"
          :strike-price="order.strike_price"
          :base-currency="order.base_currency"
          :quote-currency="order.quote_currency"
          :margin="order.margin"
          class="order-card"
        >
        </order-card>
        <order-trade-card
          :trade-list="tradeList"
          class="trade-card"
        ></order-trade-card>
      </div>
      <div class="scroll-loading" v-show="scrollLoadingVisible">
        <div id="scroll-loader"></div>
      </div>
    </div>
    <div v-if="canCancel" class="checkout">
      <button class="cancel-order-btn" @click="cancelOrder">
        <span>{{ $t("order.cancelOrder") }}</span>
      </button>
    </div>
    <tip-confirm-modal
      v-if="cancelOrderModal"
      :tip="$t('order.cancelOrderConfirm')"
      :height="192"
    >
    </tip-confirm-modal>
  </div>
</template>

<script>
import BackHeader from "@/components/layout/header/BackHeader.vue";
import OrderCard from "@/components/options/OrderCard.vue";
import { setHtmlMeta } from "@/util/utils";
import TradeCard from "@/components/options/TradeCard.vue";
import TipConfirmModal from "@/components/common/TipConfirmModal.vue";
import $axios from "@/api";
import OrderTradeCard from "@/components/options/OrderTradeCard.vue";
import BScroll from "better-scroll";
import lottie from "lottie-web";
import TabLoadingJson from "@/assets/lottie/tab_loading.json";

export default {
  name: "order",
  components: {
    OrderTradeCard,
    TipConfirmModal,
    TradeCard,
    OrderCard,
    BackHeader,
  },
  async mounted() {
    document.body.style.backgroundColor = "#F7F7F7";
    setHtmlMeta("theme-color", "#F7F7F7");
    await this.getPageData("tab");
    lottie.loadAnimation({
      container: document.getElementById("scroll-loader"),
      renderer: "svg",
      loop: true,
      autoplay: true,
      animationData: TabLoadingJson,
    });
    await this.$nextTick(() => {
      this.initScroll();
    });
  },
  computed: {
    orderStatusTitle() {
      return this.statusMap[this.order.order_status];
    },
    canCancel() {
      return this.order.order_status < 30;
    },
    showTip() {
      return this.order.order_status === 10;
    },
  },
  data() {
    return {
      cancelOrderModal: false,
      scrollLoadingVisible: false,
      scroll: {},
      order: {},
      tradeList: [],
      statusMap: {
        10: "待撮合",
        20: "已撮和",
        30: "已撮和",
        40: "已结束",
        50: "已结束",
      },
    };
  },
  methods: {
    async getPageData(animType) {
      this.startAnim(animType);
      let res = await $axios.get(`/v1/order/${this.$route.params.id}`);
      this.order = res.data;
      this.tradeList = this.order.trade_list.map((e) => {
        let funds = Number(e.amount * e.price).toFixed(2);
        funds = e.side === "ASK" ? "+" + funds : (-funds).toFixed(2).toString();
        return {
          createdAt: new Date(e.created_at).Format("yyyy/MM/dd"),
          amount: e.amount,
          funds: funds,
        };
      });
      this.scrollLoadingVisible = true;
      this.endAnim(animType);
    },
    async cancelOrder() {
      this.cancelOrderModal = true;
    },
    TipModalCancel() {
      this.cancelOrderModal = false;
    },
    async TipModalConfirm() {
      this.cancelOrderModal = false;
      this.$tabLoading.show();
      let res = await $axios.post(`/v1/order/${this.order.order_id}/cancel`);
      this.$tabLoading.hide();
      if (res.code === 0) {
        await this.getPageData();
        this.$toast.success("撤销成功");
      } else {
        this.$toast.error(res.msg);
      }
    },
    async initScroll() {
      this.scroll = new BScroll(this.$refs.wrapper, {
        // Pull up loading
        pullUpLoad: {
          // The pullingUp event is triggered when the pull-up distance exceeds 30px
          threshold: -30,
        },
        // Pull down to refresh
        pullDownRefresh: {
          // Pulling down more than 30px triggers the pullingDown event
          threshold: 30,
          // The rebound stays 20px from the top
          stop: 40,
        },
        mouseWheel: true,
        click: true,
      });
      this.scroll.on("pullingDown", async () => {
        this.getPageData("scroll").then(() => {
          this.$nextTick(() => {
            this.scroll.refresh(); // After the DOM structure changes, reinitialize BScroll
          });
          setTimeout(() => {
            // When things are done, you need to call this method to tell better-scroll that the data has been loaded,
            // otherwise the drop-down event will only be executed once
            this.scroll.finishPullDown();
            // setTimeout(()=>this.scroll = this.initScroll(),0)
          }, 100);
        });
      });
    },
    startAnim(animType) {
      if (animType === "tab") {
        this.$tabLoading.show();
      }
      if (animType === "scroll") {
        this.scrollLoadingVisible = true;
      }
    },
    endAnim(animType) {
      if (animType === "tab") {
        this.$tabLoading.hide();
        this.scrollLoadingVisible = false;
      }
      if (animType === "scroll") {
        this.$tabLoading.hide();
        this.scrollLoadingVisible = false;
      }
    },
  },
};
</script>

<style scoped lang="scss">
.empty-tip {
  margin-top: 216px;
  text-align: center;
  span {
    //width: 108px;
    height: 29px;
    background: rgba(0, 0, 0, 0.02);
    border-radius: 27px;
    padding: 8px 26px;
    font-size: 14px;
    line-height: 21px;
    text-align: center;
    color: #88888b;
    box-sizing: border-box;
    margin: 0 auto;
  }
}
.order-detail {
  padding: 12px 24px;
  margin-top: 54px;
  h1 {
    font-style: normal;
    font-weight: bold;
    font-size: 32px;
    line-height: 48px;
    color: #151516;
    margin-bottom: 20px;
  }
  .order-card {
    margin-bottom: 14px;
  }
  .trade-card {
    margin-bottom: 120px;
  }
}
.tip {
  img {
    margin-right: 14px;
  }
  display: flex;
  align-items: center;
  padding: 14px 16px;
  background: #ffffff;
  border-radius: 12px;
  font-style: normal;
  font-weight: normal;
  font-size: 12px;
  line-height: 18px;
  color: #151516;
  margin-bottom: 14px;
}
.checkout {
  padding: 10px 0 30px;
  background-color: #f7f7f7;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  margin: 0 auto 0;
  .cancel-order-btn {
    display: block;
    height: 58px;
    margin: 0 auto;
    width: calc(100% - 48px);
    background: #151516;
    box-shadow: 0px 7px 21px rgba(0, 0, 0, 0.16);
    border-radius: 8px;
    span {
      font-weight: 600;
      font-size: 17px;
      line-height: 25px;
      text-align: center;
      color: #ffffff;
    }
  }
}
.wrapper {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  right: 0;
  overflow: hidden;
  .content {
    padding-top: 10px;
  }
  .scroll-loading {
    position: absolute;
    width: 100%;
    top: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    transition: all;
    #scroll-loader {
      margin-top: 65px;
      width: 20px;
      height: 20px;
    }
  }
}
</style>
