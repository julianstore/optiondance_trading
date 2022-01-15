<template>
<div class="product-intro">
  <BackHeader back-to="/home" bg-color="#FFFFFF"/>
  <div class="faq-content">
    <h1>{{$t("faq.commonProblem")}}</h1>
    <a href="/faq/when-to-use-sell-put">{{$t("faq.whenToUseSellPut")}}</a>
    <a href="/faq/what-is-option-dance">{{$t("faq.whatIsOptionDance")}}</a>
  </div>
</div>
</template>

<script>
import BackHeader from "@/components/layout/header/BackHeader.vue";
import {setHtmlMeta} from "@/util/utils";

export default {
name: "ProductIntro",
  components: {BackHeader},
  beforeMount() {
    setHtmlMeta("theme-color","#FFFFFF")
    document.body.style.backgroundColor = "#FFFFFF";
  },
}
</script>

<style scoped lang="scss">
.faq-content{
  margin-top: 56px;
  padding: 16px 30px;
  h1{
    font-style: normal;
    font-weight: bold;
    font-size: 32px;
    line-height: 48px;
    letter-spacing: 0.03em;
    color: #151516;
    margin-bottom: 48px;
  }
  a{
    font-style: normal;
    font-weight: normal;
    font-size: 16px;
    line-height: 26px;
    letter-spacing: 0.03em;
    color: #2F80ED;
    display: block;
    margin-bottom: 18px;
    &:hover{
      text-decoration: none;
    }
  }
}

</style>
