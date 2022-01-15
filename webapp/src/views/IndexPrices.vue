<template>
<div class="index-price">
  <BackHeader back-to="-1"></BackHeader>
  <div class="content">
    <div style="padding: 0 6px">
      <p class="brand">OptionDance</p>
      <h1>{{ $t("indexPrice.btcIndexPrice") }}</h1>
      <p class="desc">{{ $t("indexPrice.btcIndexPriceDesc") }}</p>
    </div>
    <div class="prices-card">
      <img src="../assets/svg/indexprice/index_price.svg"/>
      <div class="prices">
        <ul>
          <li v-for="(e,index) in priceList">
            <span class="date">{{e.date}}</span>
            <span class="price">${{e.price}}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import {listDeliveryPrices} from "@/api/deliveryprice";
import BackHeader from "@/components/layout/header/BackHeader.vue";

export default {
name: "IndexPrices",
  components: {BackHeader},
  data() {
    return {
      priceList: []
    }
  },
  async mounted() {
    document.body.style.backgroundColor = "#F7F7F7";
    let asset = this.$route.params.asset;
    let promise = await listDeliveryPrices(asset);
     console.log(promise)
    this.priceList = promise.data
  }
}
</script>

<style scoped lang="scss">
.index-price{
  .content{
    padding: 0 24px 0;
    margin-top: 56px;
    p.brand{
      font-size: 14px;
      line-height: 21px;
      color: #88888B;
      margin-bottom: 6px;
    }
    h1{
      font-style: normal;
      font-weight: bold;
      font-size: 32px;
      line-height: 48px;
      letter-spacing: -0.02em;
      color: #151516;
      margin-bottom: 16px;
    }
    p.desc{
      font-style: normal;
      font-weight: normal;
      font-size: 13px;
      line-height: 25px;
      color: #88888B;
      margin-bottom: 36px;
      word-break:break-word;
    }
    .prices-card{
      padding: 28px 24px 41px;
      border-radius: 28px;
      background: #FFFFFF;
      img{
        display: block;
        margin-bottom: 26px;
      }
      .prices{
        li {
          margin-bottom: 16px;
          &:last-child{
            margin-bottom: 0;
          }
        }
        .date{
          font-style: normal;
          font-weight: normal;
          font-size: 13px;
          line-height: 19px;
          color: #88888B;
        }
        .price{
          font-style: normal;
          font-weight: normal;
          font-size: 14px;
          line-height: 18px;
          float: right;
          text-align: right;
          letter-spacing: 0.04em;
          color: #151516;
        }
      }
    }
  }
}
</style>