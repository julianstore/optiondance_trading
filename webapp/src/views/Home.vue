<template>
  <div>
    <div class="home-header">
      <img
        @click="openUserCenter"
        class="profile-icon"
        src="@/assets/svg/profile.svg"
      />
      <img
        class="sell-put-icon"
        v-show="option.side === 'ASK' && optionType === 'PUT'"
        @click="swapOptionType"
        src="@/assets/svg/sideType/switch/zhCn/sell_put.svg"
      />
      <img
        class="buy-put-icon"
        v-show="option.side === 'BID' && optionType === 'PUT'"
        @click="swapOptionType"
        src="@/assets/svg/sideType/switch/zhCn/buy_put.svg"
      />
      <img
        class="sell-call-icon"
        v-show="option.side === 'ASK' && optionType === 'CALL'"
        @click="swapOptionType"
        src="@/assets/svg/sideType/switch/zhCn/sell_call.svg"
      />
      <img
        class="buy-call-icon"
        v-show="option.side === 'BID' && optionType === 'CALL'"
        @click="swapOptionType"
        src="@/assets/svg/sideType/switch/zhCn/buy_call.svg"
      />
      <img
        class="test-version"
        v-show="isTestVersion"
        @click="testVersionTipVisible = true"
        src="@/assets/svg/userCenter/version/test_version.svg"
      />
      <img
        class="test-version"
        v-show="isDevVersion"
        @click="testVersionTipVisible = true"
        src="@/assets/svg/userCenter/version/dev_version_tw.svg"
      />
    </div>
    <div v-show="!unavailable">
      <div class="home">
        <div>
          <div class="top-card-group">
            <div class="top-card">
              <div class="title" @click="openStrikePriceModal">
                <span v-show="option.side === 'ASK'">{{
                  strikePriceModalTitle
                }}</span>
                <span v-show="option.side === 'BID'">{{
                  strikePriceModalTitle
                }}</span>
                <img src="@/assets/svg/arrow_left.svg" />
              </div>
              <p
                class="content"
                @click="openStrikePriceModal"
                v-if="option.strikePrice > 0"
              >
                {{ option.strikePrice }}
              </p>
              <p
                class="sold-out"
                @click="openStrikePriceModal"
                v-if="
                  side === 'ASK' &&
                  (option.strikePrice === '0' || option.strikePrice === 0)
                "
              >
                {{ $t("home.soldOut") }}
              </p>
              <p class="footer">
                {{ option.quote_currency }}/{{ option.base_currency }}
              </p>
            </div>
            <div class="top-card">
              <div class="title" @click="openExpirationDateModal">
                <span>{{ $t("home.expirationDate") }}</span>
                <img src="@/assets/svg/arrow_left.svg" />
              </div>
              <p
                class="content"
                @click="openExpirationDateModal"
                v-if="isExpiryDateInit"
              >
                {{ expirationDate }}
              </p>
              <p
                class="to-be-selected"
                @click="openExpirationDateModal"
                v-if="!isExpiryDateInit"
              >
                {{ $t("home.toBeSelected") }}
              </p>
              <p class="footer">{{ settlementDays }}</p>
            </div>
          </div>
          <div class="amount-card">
            <img class="btc-mask" src="../assets/svg/btc_mask.svg" />
            <div class="circle-mask"></div>
            <p class="title" v-show="option.side === 'ASK'">
              {{ sideString }}{{ option.base_currency }}{{ $t("home.size") }}
            </p>
            <p class="title" v-show="option.side === 'BID'">
              {{ sideString }}{{ option.base_currency }}{{ $t("home.size") }}
            </p>
            <input
              type="number"
              placeholder="0.0"
              v-model="option.size"
              v-on:input="checkOptionSize"
              inputmode="decimal"
            />
            <p class="approximate" v-show="!minimumSizeTip && !maximumSizeTip">
              ≈ {{ underlyingAmount }} {{ option.quote_currency }}
            </p>
            <p class="min-amount-tip" v-show="minimumSizeTip">
              {{ $t("home.minimum") }}{{ sideString }} {{ this.minOptionSize }}
            </p>
            <p class="min-amount-tip" v-show="maximumSizeTip">
              {{ $t("home.maximum") }}{{ sideString }} {{ this.maxSize }}
            </p>
          </div>

          <!---Premium Funds-->
          <div class="card reward-card">
            <img
              class="main-icon"
              v-show="option.side === 'ASK'"
              src="../assets/svg/lightning.svg"
            />
            <img
              class="main-icon"
              v-show="option.side === 'BID'"
              src="../assets/svg/buy_fee_icon.svg"
            />
            <span v-show="option.side === 'ASK'" @click="openRewardTip">{{
              $t("home.getItNow")
            }}</span>
            <span v-show="option.side === 'BID'" @click="openRewardTip">{{
              $t("home.paidAmount")
            }}</span>
            <div
              v-show="rewardTip"
              class="reward-tip"
              v-bind:class="{ baomai: optionType === 'PUT' && side === 'BID' }"
              id="reward-tip"
            >
              <div v-show="option.side === 'ASK'">
                <p>{{ $t("home.rewardTipAsk") }}</p>
              </div>
              <div v-show="option.side === 'BID'">
                <p>{{ $t("home.rewardTipBid") }}</p>
              </div>
            </div>
            <img
              class="info-icon"
              @click="openRewardTip"
              src="../assets/svg/info-icon.svg"
            />
            <div class="amount" v-show="option.side === 'ASK'">
              <span class="reward-amount"
                >+{{ option.funds }} {{ option.quote_currency }}</span
              >
            </div>
            <div
              class="amount"
              @click="openEditFeeModal"
              v-show="option.side === 'BID'"
            >
              <span v-show="!isBuyFeeSetup" class="fee-amount-tip">{{
                $t("home.pleaseInputPrice")
              }}</span>
              <span v-show="isBuyFeeSetup" class="fee-amount"
                >{{ optionFundsSign }}{{ option.funds }}
                {{ option.quote_currency }}</span
              >
              <img src="@/assets/svg/buy-fee-arrow.svg" />
            </div>
          </div>

          <!---annualizedRateOfReturn-->
          <div class="card reward-card" v-show="option.side === 'ASK'">
            <img
              class="main-icon"
              v-show="option.side === 'ASK'"
              src="../assets/svg/home/percent.svg"
            />
            <span v-show="option.side === 'ASK'" @click="openRewardTip">{{
              $t("home.expectedAnnualizedRateOfReturn")
            }}</span>
            <div class="amount" v-show="option.side === 'ASK'">
              <span class="reward-amount">{{ annualizedRateOfReturn }} %</span>
            </div>
          </div>

          <!---earn tip card description-->
          <div class="card tip-card">
            <img src="../assets/svg/hourglass.svg" />
            <span v-show="option.side === 'ASK'">{{
              $t("home.earnTipCardDescAsk", {
                expirationDateMonthDay: expirationDateMonthDay,
                sizeString: sizeString,
                baseCurrency: option.base_currency,
              })
            }}</span>
            <span v-show="option.side === 'BID'">{{
              $t("home.earnTipCardDescBid", {
                expirationDateMonthDay: expirationDateMonthDay,
                strikePrice: option.strikePrice,
                quoteCurrency: option.quote_currency,
                sizeString: sizeString,
                baseCurrency: option.base_currency,
              })
            }}</span>
          </div>

          <!---video tutorial-->
          <div class="video">
            <img src="@/assets/svg/home/video.svg" @click="toVideoPage" />
            <span @click="toVideoPage">{{ $t("home.videoTutorial") }}</span>
          </div>
        </div>
      </div>
      <div class="checkout">
        <p class="guard-tip mb-12">{{ $t("home.mtgSecure") }}</p>
        <button class="submit-order-btn" @click="orderConfirm">
          <span>{{ $t("home.submitOrder") }}</span>
        </button>
      </div>
    </div>
    <Unavailable v-show="unavailable" />
    <SelectorMask v-show="selectorMask" />
    <select-modal
      :type="SelectModalType.EXERCISE_PRICE_MODAL"
      :title="strikePriceModalTitle"
      :items="strikePriceList"
      :select-item="String(option.strikePrice)"
      v-show="strikePriceModal"
    />
    <select-modal
      :type="SelectModalType.EXERCISE_DATE_MODAL"
      :title="$i18n.t('home.expirationDate')"
      :items="expirationDateList"
      :select-item="expirationDate"
      v-show="expirationDateModal"
    />
    <user-center :visible="userCenterModal"></user-center>
    <order-success-modal v-show="orderSuccessModal"></order-success-modal>
    <order-confirm-modal
      v-show="orderConfirmModal"
      :size="option.size"
      :base-currency="option.base_currency"
      :strike-price="String(option.strikePrice)"
      :side="option.side"
      :option-type="optionType"
      :expiration-date="option.expirationDate"
    >
    </order-confirm-modal>
    <edit-fee
      :visible="editFeeModal"
      :amount="option.buyFee"
      :guide-price="Number(option.price)"
    >
    </edit-fee>
    <EditExpiryDate :visible="editExpiryDateModal"></EditExpiryDate>
    <EditStrikePrice :visible="editStrikePriceModal"></EditStrikePrice>
    <payment-loading v-show="paymentLoading" />
    <TestVersionTip v-show="testVersionTipVisible" />
  </div>
</template>
<script>
import { defineComponent } from "vue";
import ButtonRepo from "@/components/ButtonRepo.vue";
import $axios from "@/api";
import {
  diffDays,
  getInstrumentName,
  isDAppTokenValid,
  isJWTTokenExpired,
  setHtmlMeta,
  uuid,
} from "@/util/utils";
import SelectModal from "@/components/home/SelectModal.vue";
import UserCenter from "@/components/home/UserCenter.vue";
import { QuoteCurrency, SelectModalType } from "@/util/constants.js";
import OrderSuccessModal from "@/components/home/OrderSuccessModal.vue";
import EditFee from "@/components/home/EditFee.vue";
import { mapActions, mapGetters } from "vuex";
import PaymentLoading from "@/components/common/paymentLoading/PaymentLoading.vue";
import OrderConfirmModal from "@/components/home/OrderConfirmModal.vue";
import TestVersionTip from "@/components/home/TestVersionTip.vue";
import Unavailable from "@/components/home/Unavailable.vue";
import SelectorMask from "@/components/home/SelectorMask.vue";
import {
  getBidExpiryDates,
  leftSettlementPeriod,
  toExpiryDate,
} from "@/util/instrument";
import EditExpiryDate from "@/components/home/EditExpiryDate.vue";
import EditStrikePrice from "@/components/home/EditStrikePrice.vue";
import { optionFundsSign, SideStr, SideString } from "@/util/option";
import { getStrikePrices, listExpiryDatesByPrice } from "@/api/market";

export default defineComponent({
  components: {
    EditStrikePrice,
    EditExpiryDate,
    SelectorMask,
    Unavailable,
    TestVersionTip,
    OrderConfirmModal,
    PaymentLoading,
    EditFee,
    OrderSuccessModal,
    UserCenter,
    SelectModal,
    ButtonRepo,
  },
  beforeMount() {
    setHtmlMeta("theme-color", "#F7F7F7");
    document.body.style.backgroundColor = "#F7F7F7";
  },
  async mounted() {
    let envMode = import.meta.env.MODE;
    if (envMode === "uat") {
      this.maxOptionSize = 0.2;
    }
    this.option.quote_currency = QuoteCurrency;
    this.option.base_currency = "BTC";
    this.SelectModalType = SelectModalType;
    if (!this.optionSide || !this.optionType) {
      this.setOptionType("PUT");
      if (!this.optionSide || this.optionSide.length === 0) {
        this.setOptionSide("ASK");
      }
    }
    if (!this.deliveryType || this.deliveryType.length === 0) {
      this.setDeliveryType("CASH");
    }
    this.option.side = this.optionSide;
    let client_id = import.meta.env.VITE_MIXIN_CLIENT_ID;
    this.mixinAuthUrl = `https://mixin-www.zeromesh.net/oauth/authorize?client_id=${client_id}&scope=PROFILE:READ+PHONE:READ&response_type=code&state=login&return_to=home`;
    await this.initStrikeAndExpiry();
    //check if mixin is logged in
    if (this.mixinToken) {
      let expired = isJWTTokenExpired(this.mixinToken);
      this.isMixinLogIn = !expired;
    } else {
      this.isMixinLogIn = false;
    }
    //check mixinUser
    if (!this.mixinUser.user_id && this.isMixinLogIn) {
      let r = await $axios.get("https://mixin-api.zeromesh.net/me", {
        headers: { Authorization: "Bearer " + this.mixinToken },
      });
      if (r.error || !r.data.identity_number) {
        window.location.href = this.mixinAuthUrl;
        return;
      }
      let user = r.data;
      this.setMixinUser(user);
    }
    //check dapp login state
    if (this.isMixinLogIn) {
      let jwtTokenValid = isDAppTokenValid(this.dappToken);
      if (!jwtTokenValid) {
        let resp = await $axios.post(`/v1/auth`, {
          token: this.mixinToken,
        });
        if (resp.code === 0) {
          this.SetDAppToken(resp.data.token);
          await this.syncSettings();
        } else {
          this.$toast.error(resp.msg);
          window.location.href = this.mixinAuthUrl;
        }
      } else {
        await this.syncSettings();
      }
    }
  },
  computed: {
    ...mapGetters({
      optionType: "option/optionType",
      deliveryType: "option/deliveryType",
      optionSide: "option/optionSide",
      optionSideType: "option/optionSideType",
      mixinToken: "user/mixinToken",
      mixinUser: "user/mixinUser",
      dappToken: "user/token",
      settings: "user/settings",
    }),
    simpleMode() {
      return this.settings.app_mode === 0;
    },
    isTestVersion() {
      return import.meta.env.MODE === "uat";
    },
    isDevVersion() {
      return (
        import.meta.env.MODE === "development" ||
        import.meta.env.MODE === "beta"
      );
    },
    expirationDate() {
      if (this.option.expirationDate) {
        this.isExpiryDateInit = true;
        return this.option.expirationDate.Format("yyyy/MM/dd");
      } else {
        this.isExpiryDateInit = false;
        return this.$t("home.toBeSelected");
      }
    },
    expirationDateMonthDay() {
      if (this.option.expirationDate) {
        return this.option.expirationDate.Format("MM月dd日");
      } else return "x月x日";
    },
    settlementDays() {
      return leftSettlementPeriod(this.option.expirationDate);
    },
    expirationDateList() {
      if (this.optionSide === "ASK") {
        if (this.dateInstrumentList) {
          return this.dateInstrumentList.map((e) =>
            new Date(e.expiration_date).Format("yyyy/MM/dd")
          );
        }
      } else {
        return getBidExpiryDates();
      }
    },
    side() {
      return this.option.side;
    },
    sideString() {
      return SideString(this.option.side, this.optionType, this.$i18n.locale);
    },
    sideStr() {
      return SideStr(this.option.side, this.optionType);
    },
    underlyingAmount() {
      return this.option.size * this.option.strikePrice;
    },
    strikePriceModalTitle() {
      return this.sideString + this.$t("home.date");
    },
    sizeString() {
      if (this.option.size) return Number(this.option.size);
      else return 0;
    },
    optionFundsSign() {
      return optionFundsSign(this.option.side, this.optionType);
    },
    unavailable() {
      if (this.optionType === "CALL") {
        return true;
      }
      if (this.optionSide === "BID") {
        let isUserInWhiteList = false;
        let buyputwhitelist = import.meta.env.VITE_MIXIN_BUY_PUT_WHITE_LIST;
        if (!buyputwhitelist) {
          return false;
        }
        let strings = buyputwhitelist.split(",");
        let userId = this.mixinUser.user_id;
        for (let i in strings) {
          if (strings[i] === userId) {
            isUserInWhiteList = true;
            break;
          }
        }
        return !isUserInWhiteList;
      }
      return false;
    },
    annualizedRateOfReturn() {
      return this.calcAnnualizedRateOfReturn();
    },
  },
  data() {
    return {
      mixinAuthUrl: "",
      SelectModalType: {},
      strikePriceList: [],
      minOptionSize: 0.1,
      maxOptionSize: 99999,
      maxSize: 99999,
      option: {
        side: "",
        strikePrice: "0",
        expirationDate: null,
        size: null,
        buyFee: 500,
        guideBuyFee: 500,
        instrument_name: "",
        margin: "",
        instrument_data: "",
        funds: 0,
        price: "",
        quote_currency: QuoteCurrency,
        base_currency: "BTC",
      },

      modalType: "",
      strikePriceModal: false,
      expirationDateModal: false,
      userCenterModal: false,
      orderSuccessModal: false,
      orderConfirmModal: false,
      editFeeModal: false,
      editExpiryDateModal: false,
      editStrikePriceModal: false,
      rewardTip: false,
      minimumSizeTip: false,
      maximumSizeTip: false,
      isBuyFeeSetup: false,
      payResultQueryTask: {},
      instrumentMarketDataTask: {},
      paymentLoading: false,
      order: {
        traceId: "",
      },
      dateInstrumentList: [],
      baseCurrency: "BTC",
      baseCurrencyUSDPrice: "",
      instrument_name: "",
      bids: [],
      asks: [],
      isMixinLogIn: false,
      isExpiryDateInit: false,
      testVersionTipVisible: false,
      priceList: [
        20000, 22000, 24000, 28000, 30000, 32000, 34000, 36000, 38000, 40000,
      ],
      selectorMask: false,
    };
  },
  watch: {
    instrument_name: function (newName, old) {
      if (this.instrumentMarketDataTask) {
        clearInterval(this.instrumentMarketDataTask);
      }
      this.instrumentMarketDataTask = setInterval(async () => {
        let res = await $axios.get(`/v1/market/instrument/${newName}`);
        this.bids = res.data.bids;
        this.asks = res.data.asks;
      }, 2000);
    },
  },
  methods: {
    ...mapActions({
      SetDAppToken: "user/SetDAppToken",
      setMixinUser: "user/SetMixinUser",
    }),
    ...mapActions({
      setOptionType: "option/SetOptionType",
      setOptionSide: "option/SetOptionSide",
      setDeliveryType: "option/SetDeliveryType",
      syncSettings: "user/SyncSettings",
    }),
    checkOptionSize(e) {
      //Limit to 1 decimal place
      this.maxSize = this.maxOptionSize;
      if (this.side === "BID") {
        this.maxSize = 99999;
      }
      this.option.size = this.option.size.match(/^\d*(\.?\d{0,1})/g)[0] || null;
      if (this.option.size < this.minOptionSize) {
        this.maximumSizeTip = false;
        this.minimumSizeTip = true;
      } else if (this.option.size > this.maxSize) {
        this.option.size = String(this.maxSize);
        this.minimumSizeTip = false;
        this.maximumSizeTip = true;
        // this.option.funds = 0
        // return
      } else {
        this.minimumSizeTip = false;
        this.maximumSizeTip = false;
      }
      if (this.optionSide === "ASK") {
        // calculate margin and funds in ask side
        this.option.margin = (
          this.option.strikePrice * this.option.size
        ).toFixed(2);
        this.option.price = this.getLevel1Price("BID");
        this.option.funds = Number(
          this.option.size * this.option.price
        ).toFixed(2);
      }
      if (this.optionSide === "BID") {
        //calculate funds
        if (this.option.price > 0) {
          this.option.funds = this.calculateOptionFunds();
        }
      }
    },
    async setStrikePrice(v) {
      if (!v) return;
      this.option.strikePrice = v;
      this.option.expirationDate = null;
      this.option.price = 0;
      this.option.funds = 0;
      await this.getDateInstrumentData();
    },
    setExpirationDate(v) {
      if (!v) {
        this.option.expirationDate = null;
        return;
      }
      this.option.expirationDate = v;
      let dateStr = new Date(v).Format("yyyy-MM-dd");
      if (this.optionSide === "ASK") {
        for (let i = 0; i < this.dateInstrumentList.length; i++) {
          if (dateStr === this.dateInstrumentList[i].expiration_date_str) {
            this.option.instrument_name =
              this.dateInstrumentList[i].instrument_name;
            this.option.instrument_data = this.dateInstrumentList[i];
            this.instrument_name = this.option.instrument_name;
            break;
          }
        }
      } else {
      }
      // calculate funds
      if (this.optionSide === "ASK") {
        if (this.option.size > this.maxOptionSize) {
          this.option.funds = 0;
          return;
        }
        let level1Price = this.getLevel1Price("BID");
        this.option.funds = (this.option.size * level1Price).toFixed(2);
        this.option.margin = (
          this.option.strikePrice * this.option.size
        ).toFixed(2);
      }
    },
    getLevel1Price(side) {
      let instrument = this.option.instrument_data;
      this.bids = instrument.bids;
      this.asks = instrument.asks;
      let level1Price = "";
      if (side === "BID") {
        level1Price =
          this.bids && this.bids.length > 0 ? this.bids[0].price : "0";
      }
      if (side === "ASK") {
        level1Price =
          this.asks && this.asks.length > 0 ? this.asks[0].price : 0;
      }
      return level1Price;
    },
    hideSelectModal(e) {
      let selectModals = document.getElementsByClassName("select-modal");
      let length = selectModals.length;
      if (e) {
        for (let i = 0; i < length; i++) {
          let scroll = selectModals[i].getElementsByClassName("scroll-area")[0];
          if (
            selectModals[i].contains(e.target) &&
            !scroll.contains(e.target)
          ) {
            return;
          }
          if (scroll.contains(e.target)) {
            let target = e.target;
            let v = "";
            if (target.getElementsByTagName("span").length) {
              v = target.getElementsByTagName("span")[0].innerText;
            } else {
              v = e.target.innerText;
            }
            if (this.modalType === "price") {
              this.setStrikePrice(Number(v));
            }
            if (this.modalType === "date") {
              let date = toExpiryDate(v);
              this.setExpirationDate(date);
            }
          }
        }
      }
      this.strikePriceModal = false;
      this.expirationDateModal = false;
      document.removeEventListener("click", this.hideSelectModal);
    },
    hideBottomModal(e) {
      let uc = document.getElementsByClassName("bottom-modal-content");
      for (let i = 0; i < uc.length; i++) {
        if (uc[i].contains(e.target)) {
          return;
        }
      }
      this.userCenterModal = false;
      this.editFeeModal = false;
      this.editExpiryDateModal = false;
      this.editStrikePriceModal = false;
      document.removeEventListener("click", this.hideBottomModal);
    },
    async swapOptionType() {
      this.selectorMask = true;
    },
    openStrikePriceModal() {
      if (this.optionSide === "BID") {
        this.editStrikePriceModal = true;
        setTimeout(() => {
          document.addEventListener("click", this.hideBottomModal);
        }, 0);
        return;
      }
      this.strikePriceModal = true;
      this.modalType = "price";
      setTimeout(() => {
        document.addEventListener("click", this.hideSelectModal);
      }, 0);
    },
    setStrikePriceModal(v) {
      this.setStrikePrice(Number(v));
      this.closeEditStrikePriceModal();
    },
    openExpirationDateModal() {
      if (this.optionSide === "BID") {
        this.editExpiryDateModal = true;
        setTimeout(() => {
          document.addEventListener("click", this.hideBottomModal);
        }, 0);
        return;
      }
      this.modalType = "date";
      this.expirationDateModal = true;
      setTimeout(() => {
        document.addEventListener("click", this.hideSelectModal);
      }, 0);
    },
    openUserCenter() {
      this.userCenterModal = true;
      setTimeout(() => {
        document.addEventListener("click", this.hideBottomModal);
      }, 0);
    },
    hideRewardTip(e) {
      let uc = document.getElementById("reward-tip");
      if (!uc.contains(e.target)) {
        this.rewardTip = false;
        document.removeEventListener("click", this.hideRewardTip);
      }
    },
    setExpiryDate(v) {
      let reg = /^[1-9]\d{3}\/(0[1-9]|1[0-2])\/(0[1-9]|[1-2][0-9]|3[0-1])$/;
      let regExp = new RegExp(reg);
      if (!regExp.test(v)) {
        this.$toast.error(this.$t("message.formatError", { value: v }));
        return;
      }
      let date = toExpiryDate(v);
      this.setExpirationDate(date);
      this.closeEditExpiryDateModal();
    },
    closeEditExpiryDateModal() {
      this.editExpiryDateModal = false;
      document.removeEventListener("click", this.hideBottomModal);
    },
    closeEditStrikePriceModal() {
      this.editStrikePriceModal = false;
      document.removeEventListener("click", this.hideBottomModal);
    },
    openRewardTip() {
      document.removeEventListener("click", this.hideRewardTip);
      this.rewardTip = true;
      setTimeout(() => {
        document.addEventListener("click", this.hideRewardTip);
      }, 0);
    },
    closeUserCenter() {
      this.userCenterModal = false;
      document.removeEventListener("click", this.hideBottomModal);
    },
    closeSuccessModal() {
      this.orderSuccessModal = false;
    },
    openEditFeeModal() {
      this.editFeeModal = true;
      setTimeout(() => {
        document.addEventListener("click", this.hideBottomModal);
      }, 0);
    },
    closeEditFeeModal() {
      this.editFeeModal = false;
      document.removeEventListener("click", this.hideBottomModal);
    },
    closeTestVersionTipModal() {
      this.testVersionTipVisible = false;
    },
    setFee(price) {
      if (!price) {
        this.option.price = 0;
      } else {
        this.option.price = price.match(/^\d*(\.?\d{0,2})/g)[0] || 0;
      }
      this.option.funds = this.calculateOptionFunds();
      this.isBuyFeeSetup = true;
      this.closeEditFeeModal();
    },
    calculateOptionFunds() {
      return (this.option.price * this.option.size).toFixed(2).toString();
    },
    closeOrderConfirmModal() {
      this.orderConfirmModal = false;
    },
    calcAnnualizedRateOfReturn() {
      if (this.option.side === "ASK") {
        let level1Price = this.getLevel1Price("BID");
        this.option.price = level1Price;
        let strikePrice = this.option.strikePrice;
        let days = diffDays(this.option.expirationDate, new Date());

        let periodDays = Math.abs(Number(days));
        let rate = (36500 * level1Price) / (strikePrice * periodDays);
        return rate ? rate.toFixed(2) : 0;
      }
    },
    toVideoPage() {
      window.location.href =
        "https://option.dance/file/how-to-buy-bitcoin-at-low-prices-in-optiondance.mp4";
    },
    async orderConfirm() {
      //order check
      let valid = this.orderCheck();
      if (!valid) {
        return;
      }
      //check mixin Login
      if (!this.isMixinLogIn) {
        window.location.href = this.mixinAuthUrl;
        return;
      }
      this.orderConfirmModal = true;
    },
    async confirmOrderPay() {
      this.orderConfirmModal = false;
      await this.checkout();
    },
    async checkout() {
      this.paymentLoading = true;
      let res = {};
      this.order.traceId = uuid();
      // let price = this.calculateOptionFunds()
      let instrumentName = getInstrumentName(
        this.deliveryType,
        this.optionType,
        this.option.quote_currency,
        this.option.base_currency,
        this.option.expirationDate,
        this.option.strikePrice
      );
      let postData = {
        trace_id: this.order.traceId,
        user_id: this.mixinUser.user_id,
        quote_currency: this.option.quote_currency,
        base_currency: this.option.base_currency,
        side: this.side,
        price: Number(this.option.price).toFixed(2),
        amount: this.option.size,
        margin: String(this.option.margin),
        funds: String(this.option.funds),
        type: "L",
        instrument_name: instrumentName,
        option_type: this.optionType,
      };
      try {
        res = await $axios.post("/v1/order-request", postData);
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
          let res = await $axios.get(`/v1/order-trace/${this.order.traceId}`);
          if (res.data.order_id === this.order.traceId) {
            this.cancelPaymentResultQuery();
            this.orderSuccessModal = true;
          }
        }, 2000);
      } else {
        this.$toast.error(res.msg);
        this.cancelPaymentResultQuery();
      }
    },
    cancelPaymentResultQuery() {
      this.paymentLoading = false;
      clearInterval(this.payResultQueryTask);
    },
    async getDateInstrumentData() {
      let side = "BID";
      if (this.optionSide === "BID") {
        side = "ASK";
      }
      let res = await listExpiryDatesByPrice(
        this.option.strikePrice,
        side,
        this.optionType,
        this.deliveryType,
        this.option.quote_currency,
        this.option.base_currency
      );
      this.dateInstrumentList = res.data;
      if (this.dateInstrumentList) {
        let expirationDate = new Date(
          this.dateInstrumentList[0].expiration_date
        );
        this.setExpirationDate(expirationDate);
      }
    },
    async syncBaseCurrencyUSDPrice() {
      let res = await $axios.get(
        "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"
      );
      this.baseCurrencyUSDPrice = res.bitcoin.usd;
    },
    orderCheck() {
      let valid = true;
      //check expiry date
      if (!this.option.expirationDate) {
        this.$toast.error(this.$t("message.expiryDateTip"));
        return false;
      }
      //check option size
      if (
        !this.option.size ||
        this.option.size < this.minOptionSize ||
        this.option.size > this.maxSize
      ) {
        this.$toast.error(this.$t("message.optionSizeTip"));
        return false;
      }
      //check price and funds
      if (!this.option.funds || !this.option.price) {
        this.$toast.error(this.$t("message.optionPriceTip"));
        return false;
      }
      return valid;
    },
    async initStrikeAndExpiry() {
      if (
        this.optionSideType === "ASK_PUT" ||
        this.optionSideType === "ASK_CALL"
      ) {
        let prices = await getStrikePrices(
          this.optionSide,
          this.optionType,
          this.deliveryType,
          this.option.quote_currency,
          this.option.base_currency
        );
        this.strikePriceList = prices.data;
        if (this.strikePriceList) {
          this.option.strikePrice = String(this.strikePriceList[0]);
          this.setExpirationDate(null);
          if (this.option.strikePrice) {
            await this.getDateInstrumentData();
            if (this.dateInstrumentList) {
              this.setExpirationDate(
                new Date(this.dateInstrumentList[0].expiration_date)
              );
            }
          }
        }
      }
      this.option.size = "";
      this.maximumSizeTip = false;
      this.minimumSizeTip = false;
    },
    toOrderDetail() {
      this.$router.push(`/order/${this.order.traceId}`);
    },
    checkOptionPrice(e) {
      this.option.price = e.match(/^\d*(\.?\d{0,2})/g)[0] || 0;
      if (this.option.price < 0) {
        this.option.price = 0;
      }
    },
    async setTypeAndSide(side, type) {
      this.option.side = side;
      this.setOptionSide(this.option.side);
      this.setOptionType(type);
      if (!this.unavailable) {
        await this.initStrikeAndExpiry();
      }
      this.selectorMask = false;
    },
  },
});
</script>
