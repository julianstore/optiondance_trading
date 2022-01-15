<template>
  <div data-v-e22b9bb2="" class="main-view">
<!--    <div data-v-e22b9bb2="" class="g-body">-->
<!--      参考走势-->
<!--    </div>-->
<!--    <div data-v-e22b9bb2="" class="g-footer" style="">-->
<!--      <div data-v-e22b9bb2="" class="footer-container row">&lt;!&ndash;&ndash;&gt;-->
<!--        <div data-v-e22b9bb2="" class="footer-item row-center">-->
<!--          <div data-v-e22b9bb2="" class="btn-buy row-center">买入</div>-->
<!--        </div>-->
<!--      </div>-->
<!--    </div>-->
<!--    <TabLoading/>-->
    <div id="bm"></div>
  </div>
</template>

<script>

import dataJson from '@/assets/lottie/data.json';
import currencyDataJson from '@/assets/lottie/currency.json';
import TabLoading from "@/components/common/TabLoading/TabLoading.vue";
import lottie from "lottie-web";
export default {
  name: 'app',
  components: {TabLoading},
  mounted() {
    // this.$loading.show()
    // this.load()
    // setTimeout(()=>{
    //   this.$tabLoading.hide()
    // },5000)
    this.load()
  },
  data() {
    return {
      content:''
    }
  },
  methods: {
    load() {
      lottie.loadAnimation({
        container: document.getElementById('bm'),
        renderer: 'svg',
        loop: true,
        autoplay: true,
        animationData: currencyDataJson
      })
    }
  }
}
</script>
<style lang="scss">

</style>

