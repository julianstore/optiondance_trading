<template>
  <div>
    <back-header back-to="-1"></back-header>
    <div ref="wrapper" class="wrapper">
      <div class="position-detail">
        <h1>{{ title }}</h1>
        <div class="tip" v-if="showTip">
          <img src="@/assets/svg/tip.svg" />
          {{ $t("position.exerciseTip") }}
        </div>
        <position-card
          :display-id="true"
          :size="position.size"
          :id="position.position_id"
          :num="position.position_num_string"
          :expiration-date="position.expiration_date"
          :option-type="position.type"
          :side="position.side"
          :funds="position.funds"
          :instrument-name="position.instrument_name"
          :status="position.status"
          :price="position.price"
          :exercised-size="position.exercised_size"
          :instrument="position.instrument"
          :position-funds="position.initial_funds"
          :quote-currency="position.quote_currency"
          :base-currency="position.base_currency"
          class="position-card"
        >
        </position-card>
        <trade-card class="trade-card" :trade-list="tradeList"></trade-card>
        <div
          class="settlement-info"
          v-if="
            position.side === 'ASK' &&
            position.delivery_type === 'CASH' &&
            settlementInfo.underlying_price > 0
          "
        >
          <h1>{{ $t("position.settlementInfo") }}</h1>
          <li v-if="settlementInfo.underlying_price > 0">
            <span
              >{{ position.base_currency }}
              {{ $t("position.indexPrice") }}</span
            >
            <img
              class="info-icon"
              @click="$router.push('/index-prices/btc')"
              src="../../assets/svg/info-icon.svg"
            />
            <span class="right-item"
              >{{ settlementInfo.underlying_price }}
              {{ position.quote_currency }}</span
            >
          </li>
          <li>
            <span>{{ $t("position.refundMargin") }}</span>
            <span class="right-item"
              >{{ settlementInfo.refund_margin }}
              {{ position.quote_currency }}</span
            >
          </li>
          <li v-if="settlementInfo.size > 0">
            <span>{{ $t("position.refundSize") }}</span>
            <span class="right-item"
              >{{ settlementInfo.size }} {{ position.base_currency }}</span
            >
          </li>
        </div>
        <div class="save-profit-image" v-if="position.side === 'ASK'">
          <img src="../../assets/svg/share/download.svg" />
          <span @click="saveProfitImage">{{
            $t("position.saveProfitImage")
          }}</span>
        </div>
      </div>

      <div class="scroll-loading" v-show="scrollLoadingVisible">
        <div id="scroll-loader"></div>
      </div>
    </div>
    <div class="checkout" v-if="position.btn_status !== 'hidden'">
      <button
        v-if="position.btn_status === 'open_exercise'"
        class="exercise-btn btn-active"
        @click="TipModalDisplay"
      >
        <span>{{ $t("position.exercise") }}</span>
      </button>
      <button
        v-if="position.btn_status === 'wait_exercise'"
        class="exercise-btn btn-disable"
      >
        <span>{{ $t("position.exercise") }}</span>
        <span class="small">{{ $t("position.exerciseDesc") }}</span>
      </button>
      <button
        v-if="position.btn_status === 'wait_settlement'"
        class="exercise-btn btn-disable"
      >
        <span>{{ $t("position.waitSettlement") }}</span>
      </button>
    </div>
    <tip-confirm-modal
      v-if="exerciseModal"
      :tip="$t('position.exerciseConfirm')"
      :height="192"
    >
    </tip-confirm-modal>
    <payment-loading v-show="paymentLoading" />
  </div>
</template>

<script>
import BackHeader from "@/components/layout/header/BackHeader.vue";
import OrderCard from "@/components/options/OrderCard.vue";
import { diffDays, setHtmlMeta } from "@/util/utils";
import { parseInstrumentName } from "@/util/instrument";
import TradeCard from "@/components/options/TradeCard.vue";
import TipConfirmModal from "@/components/common/TipConfirmModal.vue";
import PositionCard from "@/components/options/PositionCard.vue";
import PaymentLoading from "@/components/common/paymentLoading/PaymentLoading.vue";
import $axios from "@/api";
import lottie from "lottie-web";
import TabLoadingJson from "@/assets/lottie/tab_loading.json";
import BScroll from "better-scroll";
import { cashExercise, physicalExercise, positionDetail } from "@/api/position";
export default {
  name: "order",
  components: {
    PaymentLoading,
    PositionCard,
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
    showTip() {
      return this.position.size > 0 && !this.isExpiry;
    },
    isExpiry() {
      return this.position.settlement === 1;
    },
    title() {
      return this.position.settlement === 1
        ? this.$i18n.t("status.end")
        : this.$i18n.t("status.open");
    },
  },
  data() {
    return {
      exerciseModal: false,
      position: {},
      settlementInfo: {},
      tradeList: [],
      paymentLoading: false,
      payResultQueryTask: {},
      scrollLoadingVisible: false,
      scroll: {},
    };
  },
  methods: {
    //animType: tab/scroll
    async getPageData(animType) {
      this.startAnim(animType);
      let axiosResponse = await positionDetail(this.$route.params.id);
      this.position = axiosResponse.data;
      this.settlementInfo = this.position.settlement_info;
      this.position.instrument = parseInstrumentName(
        this.position.instrument_name
      );
      let i = this.position.instrument;
      this.tradeList = this.position.trade_list.map((e) => {
        let funds = "";
        let side = "";
        //funds
        let number = e.amount * e.price;
        funds = number.toString();
        //side
        if (i.optionType === "PUT" && e.side === "ASK") {
          side = this.$t("options.sellPut");
        }
        if (i.optionType === "PUT" && e.side === "BID") {
          side = this.$t("options.buyPut");
        }
        funds = e.side === "ASK" ? "+" + funds : (-funds).toString();
        return {
          createdAt: new Date(e.created_at).Format("yyyy/MM/dd"),
          amount: e.amount,
          funds: funds,
          side: side,
        };
      });
      this.endAnim(animType);
    },
    async exerciseRequest() {
      let postData = {
        size: String(this.position.size),
        instrument_name: this.position.instrument_name,
        position_id: this.position.position_id,
      };
      if (this.position.delivery_type === "CASH") {
        let res = await cashExercise(postData);
        if (res.code === 0) {
          this.$toast.success(this.$i18n.t("position.exerciseSuccess"));
        }
        return;
      }
      //order check
      this.paymentLoading = true;
      let res = {};
      try {
        res = await physicalExercise(postData);
      } catch (e) {
        if (e) {
          this.cancelPaymentResultQuery();
          throw e;
        }
      }
      if (res.code === 0) {
        let codeId = res.data.code_id;
        window.open(`mixin://codes/${codeId}`);
        this.payResultQueryTask = setInterval(async () => {
          let res = await $axios.get(
            `/v1/position-status/${this.position.position_id}`
          );
          if (res.data === 20) {
            this.cancelPaymentResultQuery();
            await this.getPageData("tab");
            this.$toast.success(this.$i18n.t("position.exerciseSuccess"));
          }
        }, 2000);
      } else {
        this.cancelPaymentResultQuery();
        this.$toast.error(res.msg);
      }
    },
    cancelPaymentResultQuery() {
      this.paymentLoading = false;
      clearInterval(this.payResultQueryTask);
    },
    TipModalDisplay() {
      this.exerciseModal = true;
    },
    TipModalCancel() {
      this.exerciseModal = false;
    },
    async TipModalConfirm() {
      this.exerciseModal = false;
      await this.exerciseRequest();
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
      this.$tabLoading.hide();
      this.scrollLoadingVisible = false;
    },
    saveProfitImage() {
      let days = diffDays(
        this.position.expiration_date,
        this.position.created_at
      );
      let periodDays = Math.abs(Number(days));
      let profitRate = Number(
        (this.position.initial_funds * 36500) / this.position.margin
      ).toFixed(2);
      let instrumentName = this.position.instrument_name;
      instrumentName = instrumentName.slice(2);
      this.$router.push({
        name: "share",
        query: {
          instrumentName: instrumentName,
          side: this.position.side,
          profitRate: String(profitRate),
          optionType: this.position.optionType,
        },
      });
    },
  },
};
</script>

<style scoped lang="scss">
.position-detail {
  padding: 12px 24px;
  //margin-top: 54px;
  h1 {
    font-style: normal;
    font-weight: bold;
    font-size: 32px;
    line-height: 48px;
    color: #151516;
    margin-bottom: 20px;
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
  .position-card {
    margin-bottom: 14px;
  }
  .trade-card {
    margin-bottom: 14px;
  }
  .settlement-info {
    padding: 26px;
    background: #ffffff;
    border-radius: 12px;
    margin-bottom: 14px;
    h1 {
      font-size: 18px;
      line-height: 27px;
      color: #151516;
    }
    li {
      margin-bottom: 16px;
      position: relative;
      &:last-child {
        margin-bottom: 0;
      }
      img {
        top: 1px;
        position: absolute;
        margin-left: 6px;
      }
      span:nth-child(1) {
        font-style: normal;
        font-weight: normal;
        font-size: 13px;
        line-height: 19px;
        color: #88888b;
      }
      span.right-item {
        font-style: normal;
        font-weight: normal;
        font-size: 14px;
        line-height: 18px;
        float: right;
        letter-spacing: 0.04em;
        color: #000000;
      }
    }
  }
}

.wrapper {
  position: fixed;
  left: 0;
  top: 54px;
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
      margin-top: 20px;
      width: 20px;
      height: 20px;
    }
  }
}

.save-profit-image {
  margin-bottom: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  img {
    margin-right: 6px;
  }
  span {
    font-size: 14px;
    line-height: 21px;
    display: flex;
    align-items: center;
    letter-spacing: 0.04em;
    color: #88888b;
  }
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
  .exercise-btn {
    display: block;
    height: 58px;
    margin: 0 auto;
    width: calc(100% - 48px);
    background: rgba(21, 21, 22, 0.5);
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
  .btn-active {
    background: #151516;
    span {
      color: #ffffff;
    }
  }
  .btn-disable {
    span.small {
      display: block;
      font-style: normal;
      font-weight: normal;
      font-size: 11px;
      line-height: 16px;
      text-align: center;
      color: #ffffff;
    }
  }
}
</style>
