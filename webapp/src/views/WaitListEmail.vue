<template>
  <div>
    <BackHeader/>
    <div class="waitlist-email">
      <h4>订阅电子邮件</h4>
      <p>我们将在开发完成后，第一时间向你推送体验版</p>
      <input v-model="email" placeholder="输入你的电子邮箱地址" v-on:input="checkEmail">
      <button @click="joinWaitlist" :disabled="!isEmailComplete" v-bind:class="{ 'btn-disable': !isEmailComplete , 'btn-primary': isEmailComplete }" ><span>确认</span></button>
    </div>
  </div>
</template>

<script>
import BackHeader from "@/components/layout/header/BackHeader.vue";
import $axios from "@/api";
import {mapActions, mapGetters} from "vuex";
export default {
name: "WaitListEmail",
  components: {BackHeader},
  computed: {
    ...mapGetters({
      waitlist: 'user/waitlist'
    })
  },
  data() {
    return{
      email:'',
      isEmailComplete:false
    }
  },
  methods:{
    ...mapActions({
      'setMixinUser': 'user/SetMixinUser',
      'setWaitlist':'user/SetWaitlist'
    }),
    checkEmail() {
      let regExp = new RegExp(/^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+\.([a-zA-Z0-9_-]+)+$/);
      this.isEmailComplete= regExp.test(this.email);
    },
    async joinWaitlist() {
      let waitlist = this.waitlist;
      let data =  await $axios.post('/v1/waitlist',{
        inviter_wid: this.waitlist.inviterWid,
        type: this.waitlist.type,
        email: this.email
      });
      if (data.code===0){
        this.waitlist.email = this.email
        this.setWaitlist(this.waitlist)
        this.$toast.success('成功加入waitlist')
        this.$router.push({name:'join-success'})
      }else {
        console.log(data.msg)
      }
    }
  }
}
</script>

<style scoped lang="scss">
.waitlist-email{
  padding: 0 30px;
  text-align: left;
  h4{
    font-style: normal;
    font-weight: bold;
    font-size: 24px;
    line-height: 36px;
    color: #35384A;
    margin: 118px 0 4px;
  }
  p{
    font-style: normal;
    font-weight: normal;
    font-size: 14px;
    line-height: 21px;
    color: #959597;
    margin-bottom: 54px;
  }
  input{
    padding: 16px 18px;
    width: 100%;
    height: 56px;
    top: 115px;
    background: #F2F2F2;
    box-sizing: border-box;
    border: 1px solid #F2F2F2;
    border-radius: 8px;
    font-style: normal;
    font-weight: normal;
    font-size: 16px;
    line-height: 24px;
    letter-spacing: 0.04em;
    &:hover{
      border: 1px solid #000000;
      background: #FFFFFF;
    }
  }
  button{
    margin-top: 72px;
    display: block;
    width: 100%;
    height: 56px;
    border-radius: 8px;
    span{
      font-style: normal;
      font-weight: bold;
      font-size: 17px;
      line-height: 25px;
      text-align: center;
      color: #FFFFFF;
    }
  }
  .btn-disable{
    opacity: 0.2;
    background: #000000;
  }
  .btn-primary{
    background: #000000;
    box-shadow: 0px 7px 21px rgba(0, 0, 0, 0.16);
  }
}

</style>
