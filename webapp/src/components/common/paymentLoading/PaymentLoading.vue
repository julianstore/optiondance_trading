<template>
<div class="loading">
  <div class="center">
    <div id="payment-loader"></div>
    <span class="text">正在检查支付结果...</span>
    <button @click="close()"><span>取消</span></button>
  </div>
</div>
</template>

<script>
import lottie from "lottie-web";
import LoadingJson from "@/assets/lottie/loading.json";

export default {
name: "PaymentLoading",
  mounted() {
      lottie.loadAnimation({
        container: document.getElementById('payment-loader'),
        renderer: 'svg',
        loop: true,
        autoplay: true,
        animationData: LoadingJson
      })
  },
  methods:{
    close() {
      this.$parent.cancelPaymentResultQuery()
    }
  }
}
</script>

<style scoped lang="scss">
.loading{
  position: fixed;
  background: rgba(0, 0, 0, 0.7);
  margin: 0;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 2000;
  .center{
    width: 200px;
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    margin: auto;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    #payment-loader{
      width: 60px;
      height: 67px;
      margin-bottom: 14px;
    }
    .text{
      display: block;
      font-style: normal;
      font-weight: normal;
      font-size: 14px;
      line-height: 21px;
      color: #FFFFFF;
      margin-bottom: 48px;
    }
    button{
      width: 112px;
      height: 38px;
      border: 1px solid #FFFFFF;
      box-sizing: border-box;
      border-radius: 38px;
      background-color: rgba(0, 0, 0, 0);
      vertical-align: center;
      span{
        font-size: 16px;
        line-height: 24px;
        color: #FFFFFF;
      }
    }
  }
}
</style>