<template>
  <div class="back-header" v-bind:style="{'background-color':backGroundColor}">
    <img src="../../../assets/svg/back_arrow.svg" @click="back"/>
  </div>
</template>

<script>
export default {
name: "BackHeader",
    props:{
        backTo:String,
        bgColor:String
    },
    computed:{
      backGroundColor() {
        return !this.bgColor ? '#F7F7F7': this.bgColor
      }
    },
  methods:{
        back() {
            if (this.backTo === '-1') {
              if (window.history.length <= 1) {
                this.$router.push('/home')
                return
              }else {
                this.$router.go(-1)
                return
              }
            }
            this.$router.push(this.$props.backTo)
        }
    }
}
</script>

<style scoped>

</style>
