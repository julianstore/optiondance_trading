<template>
  <div class="trade-card">
    <table>
      <tr class="t-header">
        <th>{{ $t("card.matchedDate") }}</th>
        <th>{{ $t("card.size") }}</th>
        <th>{{ $t("card.premium") }}</th>
      </tr>
      <tr v-for="(e,i) in tradeList">
        <td>{{ e.createdAt }}</td>
        <td>{{ e.amount }}</td>
        <td>{{ e.funds }}</td>
      </tr>
    </table>
  </div>
</template>

<script>
export default {
  name: "OrderTradeCard",
  props:{
    tradeList:Array
  },
  computed:{

  }
}
</script>

<style scoped lang="scss">
.trade-card{
  background: #FFFFFF;
  border-radius: 12px;
  padding: 12px 26px;
  table{
    width: 100%;
    tr{
      margin-bottom: 14px;
    }
    th,td{
      width: 33%;
      text-align: center;
      &:nth-child(1){
        text-align: left;
      }
      &:last-child{
        text-align: right;
      }
    }
    th{
      font-style: normal;
      font-weight: normal;
      font-size: 13px;
      line-height: 19px;
      color: #88888B;
    }
    td{
      font-style: normal;
      font-weight: normal;
      font-size: 14px;
      line-height: 18px;
      letter-spacing: 0.04em;
      color: #000000;
    }
    border-collapse:separate;
    border-spacing:0px 14px;
  }
}
</style>