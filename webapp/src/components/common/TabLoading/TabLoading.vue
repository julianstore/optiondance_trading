<template>
  <div class="tab-loading">
    <div class="center">
      <div id="tab-loader"></div>
    </div>
  </div>
</template>

<script>
import lottie from "lottie-web";
import TabLoadingJson from "@/assets/lottie/tab_loading.json";

export default {
  name: "TabLoading",
  mounted() {
    lottie.loadAnimation({
      container: document.getElementById('tab-loader'),
      renderer: 'svg',
      loop: true,
      autoplay: true,
      animationData: TabLoadingJson
    })
  },
}
</script>

<style scoped lang="scss">
.tab-loading {
  position: fixed;
  //background: rgba(0, 0, 0, 0.7);
  margin: 0;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 2000;

  .center {
    width: 200px;
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    margin: auto;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;

    #tab-loader {
      width: 32px;
      height: 32px;
    }
  }
}
</style>