<template>
  <div class="sell-put-intro">
    <BackHeader back-to="-1" bg-color="#FFFFFF"/>
    <div class="article content">
      <h1>什么时候使用优买？</h1>
      <h3>1. 想低价买入 BTC </h3>
      <p>2021年6月15日，BTC报价$40,422.84。</p>
      <img src="@/assets/img/faq/CoinMarketCap.png">
      <p class="img-desc">数据来源：CoinMarketCap</p>
      <p>用户翠花曾在6月8号以 $32,000 的价格买入过 BTC，看到今天的价格她觉得有些高了。她心中盘算著：“也许月底币价会降到 $28,000 左右”，于是她决定去优买下单。</p>
      <img src="@/assets/img/faq/SellPut.png">
      <p>看到图片中的“立即获得”后你可能产生了疑问，为什么翠花买入 BTC 还会收到一笔钱？</p>
      <p>这笔钱是权利金，具体是这样产生的，市场中有个比特币大户 <b>尼古拉斯·赵四</b> ，天天担心币价归零，于是便想约定一个保准能成交的卖出价 $28,000。</p>
      <p>这价格正好和翠花的买入心理价位匹配，翠花当然不会白白提供这样的好处给赵四，于是向赵四收取了一定量的权利金。</p>
      <p>由于这是翠花和赵四的一个约定，所以叫 <b>合约(contract)</b> 。合约达成后权利金将立即转入翠花的账户，同时不论未来赵四是否使用这张合约，权利金都不退还。</p>
      <p>最后，多亏了 <b>Mixin Network</b> ，翠花和赵四不用在病情肆虐时期冒著被感染的风险出门签订纸质合约，直接在 <b>Option Dance</b> 里就可以安全地完成全部流程。</p>
      <p>到2021年06月30日，如果赵四决定卖出 BTC，翠花便可收到 1 BTC ，赵四则收到 $28,000。如果赵四不想卖了，那么翠花便收到系统退回的购币金$28,000。</p>
      <img src="@/assets/img/faq/Expiry.png">
<!--      <h3>2. 预期币价会上涨，赚取权利金 </h3>-->
<!--      <p>内容创作中...</p>-->
<!--      <img src="@/assets/img/faq/joke.png" class="mb-0">-->
      <div class="border-bottom "></div>
    </div>
    <join-group></join-group>
  </div>
</template>

<script>
import BackHeader from "@/components/layout/header/BackHeader.vue";
import JoinGroup from "@/components/faq/JoinGroup.vue";
import {setHtmlMeta} from "@/util/utils";
export default {
name: "SellPutIntro",
  components: {JoinGroup, BackHeader},
  mounted() {
    setHtmlMeta("theme-color","#FFFFFF")
    document.body.style.backgroundColor = "#FFFFFF";
  }
}
</script>

<style scoped lang="scss">
.content{
  padding: 0 30px;
  margin: 70px 0 32px;
  .border-bottom{
    height:  48px;
    border-bottom: 0.5px solid rgba(136, 136, 139, 0.3);
  }
}
</style>