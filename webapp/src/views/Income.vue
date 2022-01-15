<template>
  <div class="income">
    <BackHeader back-to="/home"></BackHeader>
    <div class="content">
      <h1>{{$t('income.incomeAnalysis')}}</h1>
      <span class="title">{{$t('income.sellPut')}}</span>
      <IncomeChart v-if="loaded" :data-set="premiumDataSet" type="premium" :year="year"/>
      <MonthlyProfitCard :profit-list="monthlyPremiumList" :year="year" :total-amount="premiumTotal" type="premium" :asset="premiumAsset"/>

      <span class="title">{{$t('income.sellPutEarned')}}</span>
      <IncomeChart v-if="loaded" :data-set="underlyingDataSet" type="underlying" :year="year"/>
      <MonthlyProfitCard :profit-list="monthlyUnderlyingList" :year="year" :total-amount="underlyingTotal" type="underlying" :asset="underlyingAsset"/>
    </div>
  </div>
</template>

<script>

import BackHeader from "@/components/layout/header/BackHeader.vue";
import IncomeChart from "@/components/chart/IncomeChart.vue";
import {annualSellPutPremium, annualSellPutUnderlying} from "@/api/statistics";
import MonthlyProfitCard from "@/components/income/MonthlyProfitCard.vue";
export default {
name: "income",
  components: {MonthlyProfitCard, IncomeChart, BackHeader},
  async mounted() {
    document.body.style.backgroundColor = "#F7F7F7";
    await this.getPageData()
    this.loaded = true
  },
  data() {
    return {
      loaded:false,
      monthlyPremiumList: [],
      monthlyUnderlyingList: [],
      premiumDataSet :[],
      underlyingDataSet :[],
      premiumTotal:0,
      underlyingTotal:0,
      premiumAsset:'',
      underlyingAsset:'',
      year: '',
    }
  },
  methods:{
    async getPageData() {
      let underlyingRes = await annualSellPutUnderlying();
      let premiumRes = await annualSellPutPremium();

      // let underlyingRes = JSON.parse(`{"code":0,"data":{"asset":"BTC","total_amount":"12.1666","year":2021,"monthly_profit_list":[{"year":2021,"month":"1","amount":"0","asset":"BTC"},{"year":2021,"month":"2","amount":"0","asset":"BTC"},{"year":2021,"month":"3","amount":"0","asset":"BTC"},{"year":2021,"month":"4","amount":"0","asset":"BTC"},{"year":2021,"month":"5","amount":"5.3333","asset":"BTC"},{"year":2021,"month":"6","amount":"5.3333","asset":"BTC"},{"year":2021,"month":"7","amount":"1.5000","asset":"BTC"},{"year":2021,"month":"8","amount":"0","asset":"BTC"},{"year":2021,"month":"9","amount":"0","asset":"BTC"},{"year":2021,"month":"10","amount":"0","asset":"BTC"},{"year":2021,"month":"11","amount":"0","asset":"BTC"},{"year":2021,"month":"12","amount":"0","asset":"BTC"}]},"msg":"success"}`)
      // let premiumRes = JSON.parse(`{"code":0,"data":{"asset":"USDT","total_amount":"9104.84","year":2021,"monthly_profit_list":[{"year":2021,"month":"1","amount":"0","asset":"USDT"},{"year":2021,"month":"2","amount":"0","asset":"USDT"},{"year":2021,"month":"3","amount":"0","asset":"USDT"},{"year":2021,"month":"4","amount":"0","asset":"USDT"},{"year":2021,"month":"5","amount":"1815.2990","asset":"USDT"},{"year":2021,"month":"6","amount":"1603.4730","asset":"USDT"},{"year":2021,"month":"7","amount":"5686.0680","asset":"USDT"},{"year":2021,"month":"8","amount":"0","asset":"USDT"},{"year":2021,"month":"9","amount":"0","asset":"USDT"},{"year":2021,"month":"10","amount":"0","asset":"USDT"},{"year":2021,"month":"11","amount":"0","asset":"USDT"},{"year":2021,"month":"12","amount":"0","asset":"USDT"}]},"msg":"success"}`)

      this.year = underlyingRes.data.year

      this.premiumAsset = premiumRes.data.asset
      this.underlyingAsset = underlyingRes.data.asset

      this.monthlyPremiumList = premiumRes.data.monthly_profit_list
      this.monthlyUnderlyingList = underlyingRes.data.monthly_profit_list

      this.premiumTotal = premiumRes.data.total_amount
      this.underlyingTotal = underlyingRes.data.total_amount

      this.underlyingDataSet = this.monthlyUnderlyingList.map(e=>Number(e.amount));
      this.premiumDataSet = this.monthlyPremiumList.map(e=>Number(e.amount));
    }
  }
}
</script>

<style scoped lang="scss">
.income{
  .content{
    margin-top: 54px;
    padding: 12px 24px;
    h1{
      font-style: normal;
      font-weight: bold;
      font-size: 32px;
      line-height: 48px;
      color: #151516;
      margin-bottom: 20px;
    }
    .title{
      border-radius: 4px;
      font-style: normal;
      font-weight: normal;
      font-size: 16px;
      line-height: 24px;
      color: #151516;
      opacity: 0.8;
      padding: 4px 10px;
      background-color:  rgba(191, 191, 191, 0.2);
      margin-bottom: 16px;
      display: inline-block;
    }
    .graph{
      height: 327px;
      background-color: #02AAB0;
      margin-bottom: 16px;
    }
  }
}

</style>