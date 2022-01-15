<template>
  <teleport to="body">
    <div v-if="visible" class="toast">
      <div class="content">
        <img v-show="type==='success'" src="../../../assets/svg/msg_success.svg"/>
        <img v-show="type==='error'" src="../../../assets/svg/msg_error.svg"/>
        <span class="message">{{message}}</span>
      </div>
    </div>
  </teleport>
</template>

<script>

export default {
  name: "Toast",
  props:{
    message:String,
    type:String
  },
  data() {
    return {
      visible: false
    }
  },
  mounted() {
    let toast = document.getElementsByClassName('toast');
    if (toast.length === 0) {
      this.visible = true;
      setTimeout(()=>{
        this.visible = false
      },3000)
    }
  },
}
</script>

<style scoped lang="scss">
.toast{
  display: flex;
  justify-content: center;
    .content{
      z-index: 1600;
      padding: 0 27px;
      display: flex;
      margin: 0 auto;
      justify-content: left;
      align-items: center;
      position: absolute;
      min-width: 100px;
      height: 48px;
      top: 60px;
      background: #FFFFFF;
      box-shadow: 0px 6px 10px rgba(149, 149, 151, 0.2);
      border-radius: 24px;
      .message{
        margin-left: 9px;
      }
    }
}
</style>
