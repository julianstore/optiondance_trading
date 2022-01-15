<template>
  <div class="coinTextRoll" >
    <TextRoll :style="{color:`${currentColor}` }" :roll-word="currentCoin" type="letter" :dur="270" :cell-height="37"/>
  </div>
</template>

<script>
import TextRoll from "@/components/common/TextRoll.vue";
export default {
name: "CoinTextRoll",
  components: {TextRoll},
  props: {
      type:String,
      startRoll:Boolean
    },
    data() {
    return{
      index:1,
      currentCoin: 'BTC',
      currentColor: '#FFA726',
      coinList: [
        {
          name:'BTC',
          color:'#FFA726'
        },
        {
          name:'XIN',
          color:'#29B6F6'
        },
        {
          name:'ETH',
          color:'#8B89F5'
        },
      ]
    }
  },
  mounted() {
    switch (this.$props.type) {
      case 'static':
        this.currentCoin = this.coinList[0].name
        this.currentColor = this.coinList[0].color
        break
      case 'roll':
          setInterval(() => {
            this.currentCoin = this.coinList[this.index].name
            this.currentColor = this.coinList[this.index].color
            this.index++
            if (this.index === 3) {
              this.index = 0
            }
          }, 1660)
            break
    }
  },
  watch: {
    startRoll(value, oldVal) {
      if(value){
        setInterval(() => {
          this.currentCoin = this.coinList[this.index].name
          this.currentColor = this.coinList[this.index].color
          this.index++
          if (this.index === 3) {
            this.index = 0
          }
        }, 1660)
      }
    },
  }
}
</script>

<style scoped>

</style>
