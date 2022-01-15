<template>
<div class="join-group">
  <h4>{{ $t("faq.joinGroup") }}</h4>
  <div class="card" @click="joinGroup">
    <img src="@/assets/svg/faq/od_group.svg">
    <div>
      <h3>{{ $t("faq.odGroup") }}</h3>
      <p>{{ $t("faq.groupDesc") }}</p>
    </div>
  </div>
</div>
</template>

<script>
export default {
name: "JoinGroup",
  methods:{
    joinGroup() {
      window.open('mixin://apps/a757520f-cfd3-4a4b-a1c8-85647e908ee0')
    }
  }
}
</script>

<style scoped lang="scss">

.join-group{
  padding: 0 30px;
  h4{
    margin: 48px 0 14px;
    font-style: normal;
    font-weight: bold;
    font-size: 20px;
    line-height: 30px;
    letter-spacing: 0.03em;
    color: #151516;
  }
  .card{
    display: flex;
    flex-direction: row;
    align-items: flex-start;
    padding: 22.5px 25px;
    box-sizing: border-box;
    height: 88px;
    background: #F7F7F7;
    border-radius: 8px;
    margin-bottom: 84px;
    img{
      margin-right: 18px;
    }
    h3{
      font-style: normal;
      font-weight: bold;
      font-size: 18px;
      line-height: 27px;
      letter-spacing: 0.03em;
      color: #151516;
    }

    p{
      font-style: normal;
      font-weight: 300;
      font-size: 11px;
      line-height: 16px;
      letter-spacing: 0.03em;
      color: #88888B;
    }
  }
}
</style>