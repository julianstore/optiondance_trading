<template>
<div class="wait-rank">
  <BackHeader/>
  <div class="content">
    <h4>你现在的位置是</h4>
    <h1 class="mb-40">#{{ rank }}</h1>
    <h5>分享给好友可提升你的排名</h5>
    <h5 class="mb-40">前10名将获得 500 USDT 体验金</h5>
    <button class="share-btn" @click="share"><span>分享给 Mixin 好友</span></button>
    <p class="mb-8">你推荐的每个人都会使你前进排名</p>
    <p>（ 你已邀请 {{ inviteCount }} ）</p>
  </div>

</div>
</template>

<script>
import {mapGetters} from "vuex";
import $axios from "@/api";
import BackHeader from "@/components/layout/header/BackHeader.vue";

export default {
name: "WaitListRank",
  components: {BackHeader},
  computed: {
    ...mapGetters({
      // map `this.doneCount` to `this.$store.getters.doneTodosCount`
      mixinUser: 'user/mixinUser'
    })
  },
  async mounted() {
    let res = await $axios.get(`/v1/waitlist/rank-info`,{
      params:{
        mid: this.mixinUser.identity_number
      }
    });
    this.inviteCount = res.data.invite_count
    this.rank = res.data.rank
  },
  data() {
    return {
      inviteCount:0,
      rank:1
    }
  },
  methods:{
    share() {
      const data = `${import.meta.env.VITE_APP_URL}/waitlist?mid=${this.mixinUser.identity_number}`
      let base64Data = window.btoa(data)
      let openUrl = `mixin://send?category=text&data=${encodeURIComponent(base64Data)}`
      window.open(openUrl)
    }
  }
}
</script>

<style scoped lang="scss">
.wait-rank{
  .header{
    padding: 16px 24px;
  }
  .content{
    margin-top: 56px;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    h1 ,h4 ,h5 ,p{
      font-style: normal;
      font-weight: normal;
      text-align: center;
      color: #35384A;
    }
    h4{
      font-size: 18px;
      line-height: 27px;
    }
    h5{
      font-size: 14px;
      line-height: 21px;
    }
    h1{
      font-weight: bold;
      font-size: 60px;
      line-height: 90px;
      letter-spacing: 0.02em;
    }
    p{
      font-size: 14px;
      line-height: 21px;
      color: #959597;
    }
    .share-btn{
      display: flex;
      flex-direction: row;
      justify-content: center;
      align-items: center;
      padding: 14px 50px;
      position: static;
      width: 253px;
      height: 58px;
      background: #4CA1EE;
      box-shadow: 0px 6px 12px rgba(76, 161, 238, 0.25);
      border-radius: 54px;
      margin-bottom: 88px;
      span{
        font-style: normal;
        font-weight: bold;
        font-size: 18px;
        line-height: 27px;
        text-align: center;
        color: #FFFFFF;
      }
    }
  }
}
</style>
