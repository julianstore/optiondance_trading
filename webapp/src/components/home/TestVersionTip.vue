<template>
  <div class="modal-bg">
    <div class="test-version-modal">
      <p class="title">{{ $t("home.testVersionIntro") }}</p>
      <p class="content text-wrapper">
        {{$t("home.testVersionDesc",{quoteCurrency:quoteCurrency} )}}
      </p>
      <div class="button-group">
        <button @click="ok"><span>{{ $t("global.know") }}</span></button>
      </div>
    </div>
  </div>
</template>

<script>
import {QuoteCurrency} from "@/util/constants";

export default {
name: "TestVersionTip",
  mounted() {
    this.quoteCurrency = QuoteCurrency
  },
  data() {
  return{
    quoteCurrency : QuoteCurrency
  }
  },
  methods:{
    ok() {
      this.$parent.closeTestVersionTipModal();
    },
  }
}
</script>


<style scoped lang="scss">
.test-version-modal {
  width: 323px;
  height: 238px;
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
    padding: 0 17px;
    font-size: 14px;
    line-height: 149%;
    text-align: center;
    font-feature-settings: 'halt' on;
    color: #151516;
    margin-bottom: 40px;
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
      width: 100%;
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