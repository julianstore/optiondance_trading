<template>
<div class="selector-mask">
  <ul>
    <li @click="setTypeAndSide('ASK','PUT')">
      <img src="@/assets/svg/selector_mask/sell_put.svg">
      <div class="right">
        <h5>{{ $t("options.sellPut") }}</h5>
        <p>{{ $t("options.sellPutDesc") }}</p>
      </div>
    </li>
    <li v-if="!simpleMode" @click="setTypeAndSide('BID','PUT')">
      <img src="@/assets/svg/selector_mask/buy_put.svg">
      <div class="right">
        <h5>{{ $t("options.buyPut") }}</h5>
        <p>{{ $t("options.buyPutDesc") }}</p>
      </div>
    </li>
    <li @click="setTypeAndSide('ASK','CALL')">
      <img src="@/assets/svg/selector_mask/sell_call.svg">
      <div class="right">
        <h5>{{ $t("options.sellCall") }}</h5>
        <p>{{ $t("options.sellCallDesc") }}</p>
      </div>
    </li>
    <li v-if="!simpleMode" @click="setTypeAndSide('BID','CALL')">
      <img src="@/assets/svg/selector_mask/buy_call.svg">
      <div class="right">
        <h5>{{ $t("options.buyCall") }}</h5>
        <p>{{ $t("options.buyCallDesc") }}</p>
      </div>
    </li>
  </ul>
  <div class="toFaq" @click="$router.push('/faq')">
    <img src="@/assets/svg/faq/faq.svg">
    <span>{{ $t("faq.commonProblem") }}</span>
  </div>
</div>
</template>

<script>
import {mapGetters} from "vuex";

export default {
name: "SelectorMask",
  computed: {
    ...mapGetters({
      settings: 'user/settings'
    }),
    simpleMode() {
      return this.settings.app_mode === 0
    }
  },
  methods:{
    setTypeAndSide(side,type) {
      this.$parent.setTypeAndSide(side,type);
    }
  }
}
</script>

<style scoped lang="scss">
.selector-mask{
  position: absolute;
  top:60px;
  background-color: #F7F7F7;
  height: calc(100vh - 60px);
  z-index: 300;
  width: 100%;
  ul{
    padding: 8px 24px 0;
  }
  li{
    height: 84px;
    border-radius: 12px;
    background-color: #FFFFFF;
    margin-bottom: 16px;
    padding: 21px;
    box-sizing: border-box;
    display: flex;
    align-items: center;
    &:last-child{
      margin-bottom: 0;
    }
    .right{
      margin-left: 18px;
      display: inline-block;
      h5{
        font-weight: bold;
        font-size: 16px;
        line-height: 24px;
        color: #151516;
      }
      p{
        font-style: normal;
        font-weight: 300;
        font-size: 11px;
        line-height: 16px;
        color: #88888B;
      }
    }
  }
  .toFaq{
    display: flex;
    align-items: center;
    justify-content: center;
    img{
      display: inline-block;
      margin-right: 4px;
    }
    margin:32px auto 0;
    span{
      font-style: normal;
      font-weight: 300;
      font-size: 14px;
      line-height: 21px;
      color: #88888B;
    }
  }
}
</style>