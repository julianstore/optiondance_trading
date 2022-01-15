<template>
<div class="option-card" @click="toOrderDetail">
  <div class="top">
    <div class="instrument">
      <span class="name">{{ title }}</span>
      <img v-show=" side==='BID' && optionType === 'PUT' "  src="@/assets/svg/sideType/card/zhCn/buy_put.svg"/>
      <img v-show=" side==='ASK' && optionType === 'PUT' " src="@/assets/svg/sideType/card/zhCn/sell_put.svg"/>
      <img v-show=" side==='BID' && optionType === 'CALL' " src="@/assets/svg/sideType/card/zhCn/buy_call.svg"/>
      <img v-show=" side==='ASK' && optionType === 'CALL' " src="@/assets/svg/sideType/card/zhCn/sell_call.svg"/>
    </div>
    <p>{{ $t("card.strikePrice") }}</p>
    <ul>
      <li>
        <span class="key">{{ $t("card.expiryDate") }}</span>
        <span class="value">{{ expirationDateString }}</span>
      </li>
      <li>
        <span class="key">{{ $t("card.orderStatus") }}</span>
        <span class="value">{{ orderStatusString }}</span>
      </li>
      <li>
        <span class="key">{{ $t("card.totalAndClosed") }}</span>
        <span class="value">{{totalAmount}}/{{filledAmountString}} BTC</span>
      </li>
      <!-- sell put-->
      <li v-if="side==='ASK' && optionType === 'PUT'">
        <span class="key">{{ $t("card.margin") }}</span>
        <span class="value">{{ marginString }} {{ quoteCurrency }}</span>
      </li>
      <li v-if="side==='ASK'">
        <span class="key">{{ $t("card.extraFunds") }}</span>
        <span class="value">{{ filledFunds }} {{ quoteCurrency }} </span>
      </li>
      <!-- buy put && buy call-->
      <li v-if="side==='BID' ">
        <span class="key">{{ $t("card.funds") }}</span>
        <span class="value">{{ totalFunds }} {{ quoteCurrency }}</span>
      </li>
      <li>
        <span class="key">{{ $t("card.deliveryType") }}</span>
        <span class="value" v-if="deliveryType === 'CASH' ">{{ quoteCurrency }}</span>
        <span class="value" v-if="deliveryType === 'PHYSICAL' ">{{ baseCurrency }}</span>
      </li>
      <li v-if="displayId">
        <span class="key">{{ $t("card.orderNum") }}</span>
        <span class="value">{{ orderNum }}</span>
      </li>
    </ul>
  </div>
</div>
</template>

<script>
import {instrumentCardTitle,parseInstrumentName} from "@/util/instrument";
import {QuoteCurrency} from "@/util/constants";

export default {
  name: "OrderCard",
  mounted() {
    this.quoteCurrency = QuoteCurrency
  },
  computed:{
    expirationDateString() {
      return new Date(this.expirationDate).Format('yyyy-MM-dd hh:mm')
    },
    orderStatusString() {
      return this.statusMap[this.orderStatus]
    },
    totalAmount() {
      let totalAmount = 0
      if (this.side==="BID") {
        totalAmount = (Number(this.remainingFunds) + Number(this.filledFunds))/ Number(this.price)
      } else {
        totalAmount = Number(this.remainingAmount) + Number(this.filledAmount)
      }
      return totalAmount
    },
    totalFunds() {
      return (Number(this.remainingFunds) + Number(this.filledFunds))
    },
    filledAmountString() {
      // let filledAmount = "0"
      // if (this.side==="BID") {
      //   let totalFunds =  (Number(this.remainingFunds) + Number(this.filledFunds))
      //   let filledRatio =  Number(Number(this.filledFunds) / totalFunds)
      //   filledAmount = ((totalFunds * filledRatio)/this.price).toFixed(2)
      // }else {
      //   filledAmount = this.filledAmount
      // }
      return this.filledAmount
    },
    marginString() {
      return this.margin
    },
    instrument() {
      if(!this.instrumentName) return {}
      return parseInstrumentName(this.instrumentName);
    },
    title() {
      return instrumentCardTitle(this.instrument,this.side,this.$i18n.locale)
    },
    deliveryType() {
      return this.instrument.deliveryType
    }
  },
  props:{
    id:String,
    orderNum:String,
    side: String,
    price: String,
    optionType:String,
    orderStatus:Number,
    expirationDate: String,
    instrumentName:String,
    displayId:Boolean,
    clickable:Boolean,
    remainingAmount:String,
    filledAmount:String,
    remainingFunds:String,
    filledFunds:String,
    margin:String,
    strikePrice:String,
    quoteCurrency:String,
    baseCurrency:String
  },
  data() {
    return {
      statusMap:{
        10: this.$t("status.matching"),
        20: this.$t("status.matched"),
        30: this.$t("status.matched"),
        35: this.$t("status.endOrderCancel"),
        40: this.$t("status.endOrderCancel"),
        50: this.$t("status.endOrderMatched"),
      },
      quoteCurrency: QuoteCurrency
    }
  },
  methods:{
    toOrderDetail() {
      if(this.clickable) {
        this.$router.push('/order/'+this.id)
      }
    }
  }
}
</script>