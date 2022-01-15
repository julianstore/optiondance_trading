<template>
  <div class="option-card" @click="toPositionDetail">
    <div class="top">
      <div class="instrument">
        <span class="name">{{title}}</span>
        <img v-show=" side==='BID' && optionType === 'PUT' " src="@/assets/svg/sideType/card/zhCn/buy_put.svg"/>
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
          <span class="key">{{ $t("card.positionStatus") }}</span>
          <span class="value">{{ positionStatusString }}</span>
        </li>
        <li>
          <span class="key">{{ $t("card.exercisedStatus") }}</span>
          <span class="value">{{ exerciseStatusString }}</span>
        </li>
        <li>
          <span class="key">{{ $t("card.totalAndExercised") }}</span>
          <span class="value">{{totalSize}}/{{exercisedSize}} BTC</span>
        </li>
        <!--buy put-->
        <li v-if="side === 'BID'">
          <span class="key">{{ $t("card.funds") }}</span>
          <span class="value">{{ totalFunds }} {{ quoteCurrency }}</span>
        </li>
        <!--sell put-->
        <li v-if="side === 'ASK' && optionType === 'PUT' ">
          <span class="key">{{ $t("card.margin") }}</span>
          <span class="value">{{ margin }} {{ quoteCurrency }}</span>
        </li>
        <li v-if="side === 'ASK'">
          <span class="key">{{ $t("card.extraFunds") }}</span>
          <span class="value">{{ totalFunds }} {{ quoteCurrency }}</span>
        </li>
        <li>
          <span class="key">{{ $t("card.deliveryType") }}</span>
          <span class="value" v-if="deliveryType === 'CASH' ">{{ quoteCurrency }}</span>
          <span class="value" v-if="deliveryType === 'PHYSICAL' ">{{ baseCurrency }}</span>
        </li>
        <li v-if="displayId">
          <span class="key">{{ $t("card.positionNum") }}</span>
          <span class="value">{{ num }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import {instrumentCardTitle} from "@/util/instrument";

export default {
  name: "PositionCard",
  computed:{
    expirationDateString() {
      return new Date(this.expirationDate).Format('yyyy-MM-dd hh:mm')
    },
    exerciseStatusString() {
      return this.statusMap[this.status]
    },
    positionStatusString() {
      let now = new Date().getTime()
      let expirationDate = new Date(this.expirationDate).getTime();
      return now > expirationDate ? this.$t("status.endExpiry") : this.$t("status.open")
    },
    totalSize() {
      return Math.abs(Number(this.size))
    },
    margin() {
      if (this.side === "BID") {

      } else {
        if (this.instrument) {
          return this.totalSize * this.instrument.strikePrice;
        }
      }
    },
    title() {
      return instrumentCardTitle(this.instrument,this.side,this.$i18n.locale)
    },
    deliveryType() {
      if (!this.instrument) return ""
      return this.instrument.deliveryType
    },
    totalFunds() {
      return Math.abs(this.positionFunds)
    }
  },
  props:{
    id:String,
    num: String,
    side: String,
    price: String,
    size: Number,
    optionType:String,
    status: Number,
    expirationDate: String,
    instrumentName:String,
    displayId:Boolean,
    clickable:Boolean,
    instrument:Object,
    funds:Number,
    exercisedSize:Number,
    positionFunds: Number,
    quoteCurrency:String,
    baseCurrency:String,
  },
  data() {
    return {
      statusMap: {
        10: this.$t("status.notExercised"),
        20: this.$t("status.waitSettlement"),
        30:this.$t("status.exercised"),
      },
    }
  },
  methods:{
    toPositionDetail() {
      if (this.clickable) {
        this.$router.push('/position/'+this.id)
      }
    }
  }
}
</script>