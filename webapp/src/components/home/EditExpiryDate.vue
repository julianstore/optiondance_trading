<template>
    <div class="bottom-modal">
        <div class="modal-bg" v-show="visible"></div>
        <transition name="normal">
            <div v-show="visible" class="bottom-modal-content" >
                <div class="title">
                    <span >{{ $t("home.customExpiryDate") }}</span>
                    <img @click="close" src="@/assets/svg/close.svg"/>
                </div>
                <div class="input-area">
                        <input
                               placeholder="yyyy/MM/dd"
                               v-model="expiryDate"/>
<!--                               v-on:input="checkOptionPrice"-->

<!--                    <span>{{ quoteCurrency }}</span>-->
                </div>
                <button @click="setExpiryDate"><span>{{ $t("global.save") }}</span></button>
            </div>
        </transition>
    </div>
</template>

<script>

export default {
name: "EditExpiryDate",
    props:{
        visible:Boolean,
        amount: Number,
        guidePrice: Number
    },
  data() {
      return {
        expiryDate: '',
      }
    },
    methods: {
        close() {
            this.$parent.closeEditExpiryDateModal()
        },
        setExpiryDate() {
              this.$parent.setExpiryDate(this.expiryDate)
        },
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
