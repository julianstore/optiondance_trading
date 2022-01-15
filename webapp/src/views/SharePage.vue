<template>
  <div>

    <share-card-bg id="share-card" style="display: none"
                :instrument-name="instrumentName"
                :option-type="optionType"
                :side="side"
                :profit-rate="profitRate"
    />
    <img :src="imgData">
  </div>
</template>

<script>
import ShareCardBg from "@/components/share/ShareCardBg.vue";
import { encode } from 'js-base64';
export default {
name: "SharePage",
  components: {ShareCardBg},
  mounted() {


    this.instrumentName = this.$route.query.instrumentName
    this.optionType = this.$route.query.optionType
    this.side = this.$route.query.side
    this.profitRate = this.$route.query.profitRate
    let canvasImgData = ''

    this.$nextTick(function(){
      let svg = document.getElementById("share-card");
      const s = new XMLSerializer().serializeToString(svg);
      let src = `data:image/svg+xml;base64,${encode(s)}`;
      const img = new Image();
      img.src = src;
      img.onload = () => {
        const canvas = document.createElement('canvas');
        canvas.width = img.width *3;
        canvas.height = img.height *3 ;
        const context = canvas.getContext('2d');
        context.drawImage(img, 0, 0,canvas.width,canvas.height );
        const ImgBase64 = canvas.toDataURL('image/png');
        this.imgData = ImgBase64
      }
    })
  },
  data() {
    return {
      instrumentName: '',
      optionType:'',
      side:'',
      profitRate:'',
      imgData:'',
    }
  }
}
</script>

<style scoped lang="scss">
img{
  width: 100%;
}
</style>