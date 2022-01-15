<template>
    <div class="bottom-modal">
        <div class="modal-bg" v-show="visible"></div>
        <transition name="normal">
            <div v-show="visible" class="bottom-modal-content" >
                <div class="title">
                    <span >{{ $t("home.customPrice") }}</span>
                    <img @click="close" src="@/assets/svg/close.svg"/>
                </div>
                <div class="input-area">
                        <input type="number"
                               :placeholder="$t('home.editPricePlaceholder')"
                               v-model="optionPrice"
                               v-on:input="checkOptionPrice"
                            inputmode="decimal"/>
                    <span>{{ quoteCurrency }}</span>
                </div>
                <button @click="setFee"><span>{{ $t("global.save") }}</span></button>
            </div>
        </transition>
    </div>
</template>

<script>
import {QuoteCurrency} from "@/util/constants";

export default {
name: "EditFee",
    props:{
        visible:Boolean,
        amount: Number,
        guidePrice: Number
    },
    mounted() {
      this.optionPrice = this.guidePrice
      this.quoteCurrency = QuoteCurrency
    },
  data() {
      return {
        optionPrice: 0,
        quoteCurrency :QuoteCurrency
      }
    },
    methods: {
        close() {
            this.$parent.closeEditFeeModal()
        },
        setFee() {
            this.$parent.setFee(this.optionPrice)
        },
        checkOptionPrice() {
          this.optionPrice= (this.optionPrice.match(/^\d*(\.?\d{0,2})/g)[0]) || null
        }
    }
}
</script>

<style scoped lang="scss">
.bottom-modal-content{
    height: 308px;
    .input-area{
        position: relative;
        input{
            height: 56px;
            background: #F7F7F7;
            border-radius: 12px;
            box-sizing: border-box;
            padding: 17px 12px;
            border: none;
            display: block;
            width: 100%;
            margin: 0 auto 48px;
        }
        span{
            position: absolute;
            right: 12px;
            top: 17px;
            font-style: normal;
            font-weight: normal;
            font-size: 14px;
            line-height: 21px;
            color: #151516;
        }
    }
    button{
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;
        background: #151516;
        border-radius: 8px;
        height: 56px;
        width: 157px;
        margin: 0 auto 48px;
        span{
            font-style: normal;
            font-weight: bold;
            font-size: 17px;
            line-height: 25px;
            text-align: center;
            color: #FFFFFF;
        }
    }
}

</style>
