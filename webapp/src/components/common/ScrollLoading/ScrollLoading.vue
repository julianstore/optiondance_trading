<template>
  <div ref="wrapper" class="wrapper" >
    <div class="content">
      <slot></slot>
    </div>
    <div class="scroll-loading" v-show="scrollLoadingVisible">
      <div id="scroll-loader"></div>
    </div>
  </div>
</template>

<script>
import lottie from "lottie-web";
import TabLoadingJson from "@/assets/lottie/tab_loading.json";
import BScroll from "better-scroll";

export default {
name: "ScrollLoading",
  props:{
    scrollLoadingVisible:false
  },
  async mounted() {
    console.log(11)
    lottie.loadAnimation({
      container: document.getElementById('scroll-loader'),
      renderer: 'svg',
      loop: true,
      autoplay: true,
      animationData: TabLoadingJson
    })
    await this.$nextTick(() => {
      this.$parent.initScroll(this.$refs.wrapper)
    });
    // setInterval(()=>{
    //   console.log(this.scroll)
    // },1000)
  },
  data() {
    return {
      scroll:{}
    }
  },
  methods:{
    async getPageData() {
      await this.$parent.getPageData()
    },
  }
}
</script>

<style scoped>
#scroll-loader{
  margin-top: 10px;
  width: 20px;
  height: 20px;
}
</style>