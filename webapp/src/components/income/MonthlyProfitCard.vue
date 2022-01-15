<template>
  <div class="detail">
    <div class="detail-header">
      <span class="title" v-if="type==='premium'">{{$t('income.totalProfit')}}</span>
      <span class="title" v-if="type==='underlying'">{{$t('income.totalEarned')}}</span>
      <span class="year">/ {{ year }}年</span>
      <span class="total" v-bind:class="{'color-yellow':type==='underlying'}">{{ totalAmount }} {{ asset }}</span>
    </div>
    <ul>
      <li v-for="(e,index) in dataList">
        <span class="month">{{ year }}年{{ e.month }}月</span>
        <span class="income">{{ e.amount }} {{ asset }}</span>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: "MonthlyProfitCard",
  props: {
    year: Number,
    profitList: Array,
    asset: String,
    totalAmount: Number,
    type: String,
  },
  computed: {
    dataList() {
      this.$props.profitList.splice(new Date().getMonth() + 1)
      return this.$props.profitList
    }
  },

}
</script>

<style scoped lang="scss">
.detail {
  background-color: #FFFFFF;
  //height: 305px;
  padding: 26px;
  margin: 16px 0 16px;
  border-radius: 12px;

  .detail-header {
    padding-bottom: 18px;
    border-bottom: 0.5px solid #DBDBDC;

    span {
      &.title {
        font-style: normal;
        font-weight: bold;
        font-size: 18px;
        line-height: 27px;
        color: #151516;
        margin-right: 12px;
      }

      &.year {
        font-style: normal;
        font-weight: normal;
        font-size: 13px;
        line-height: 19px;
        color: #88888B;
      }

      &.total {
        font-style: normal;
        font-weight: bold;
        font-size: 18px;
        line-height: 27px;
        align-items: center;
        text-align: right;
        float: right;
        color: #FF6F61;
      }
      &.color-yellow{
        color: rgb(255, 146, 45) !important;
      }
    }
  }

  ul {
    margin-top: 18px;

    li {
      margin-bottom: 14px;

      .month {
        font-style: normal;
        font-weight: normal;
        font-size: 13px;
        line-height: 19px;
        color: #88888B;
      }

      .income {
        font-style: normal;
        font-weight: normal;
        font-size: 13px;
        line-height: 19px;
        float: right;
        color: #151516;
      }
      &:last-child {
        margin-bottom: 0px;
      }
    }
  }
}
</style>