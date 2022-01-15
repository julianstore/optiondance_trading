<template>
  <div class="modal-bg">
    <div class="order-confirm-modal">
      <p class="title">{{ $t("home.orderConfirm") }}</p>
      <p class="content text-wrapper">
        {{ $t("home.orderConfirmText",{
        expiry:expiry,
        strikePriceString : strikePriceString,
        sideString : sideString,
        sizeString:sizeString,
        baseCurrency :baseCurrency
      })}}
      </p>
      <div class="button-group">
        <button @click="cancel"><span>{{$t("global.cancel")}}</span></button>
        <button @click="pay"><span>{{$t("global.confirm")}}</span></button>
      </div>
    </div>
  </div>
</template>

<script>
import {toMoney} from "@/util/utils";
import {SideString} from "@/util/option";

export default {
  name: "OrderConfirmModal",
  computed:{
    expiry() {
      return new Date(this.expirationDate).Format('yyyy年MM月dd日hh:mm')
    },
    strikePriceString() {
      return toMoney(this.strikePrice)
    },
    sideString() { return SideString(this.side,this.optionType,this.$i18n.locale) },
    sizeString() {return Number(this.size) }
  },
  props:{
    expirationDate: Date,
    strikePrice:String,
    size:String,
    baseCurrency:String,
    side:String,
    optionType:String
  },
  methods: {
    cancel() {
      this.$parent.closeOrderConfirmModal();
    },
    pay() {
      this.$parent.confirmOrderPay();
    }
  }
}
</script>

<style scoped lang="scss">
.order-confirm-modal {
  width: 323px;
  height: 216px;
  background: #FFFFFF;
  box-sizing: border-box;
  border-radius: 12px;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  margin: auto;
  padding: 24px 16px 16px;

  img {
    width: 72px;
    height: 72px;
    margin: 28px auto 14px;
    display: block;
  }
  p{
    text-align: center;
    font-style: normal;
    color: #151516;
  }
  p.title {
    text-align: center;
    font-style: normal;
    font-weight: bold;
    font-size: 19px;
    line-height: 28px;
    margin-bottom: 12px;
  }
  p.content{
    font-style: normal;
    font-weight: normal;
    font-size: 16px;
    line-height: 24px;
    text-align: center;
    margin-bottom: 33px;
  }

  .button-group {
    display: flex;
    flex-direction: row;
    justify-content: center;
    position: relative;
    bottom: 0;

    button {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 140px;
      height: 54px;
      box-sizing: border-box;
      border: 1px solid #151516;
      border-radius: 8px;

      span {
        font-style: normal;
        font-weight: 600;
        font-size: 17px;
        line-height: 25px;
      }

      &:nth-child(1) {
        background: #FFFFFF;
        margin-right: 12px;

        span {
          color: #151516;
        }
      }

      &:nth-child(2) {
        background: #111111;
        color: #FFFFFF;
      }
    }
  }

}
</style>