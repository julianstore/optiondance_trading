<template>
<div>
  <BackHeader back-to="/home"/>
  <div class="settings">
    <h1>{{ $t("settings.setting") }}</h1>
    <div class="setting-item">
      <div>
        <h2>{{ $t("settings.simpleMode") }}</h2>
        <p>{{ $t("settings.simpleModeDesc") }}</p>
      </div>
      <div class="switch" @click="switchVersion">
        <SwitchOn v-if="simpleMode"/>
        <SwitchOff v-if="!simpleMode"/>
      </div>
    </div>

    <div class="setting-item">
      <div>
        <h2>{{ $t("settings.cashDelivery") }}</h2>
        <p>{{ $t("settings.cashDeliveryDesc") }}</p>
      </div>
      <div class="switch" @click="switchDeliveryType">
        <SwitchOn v-if="cashDelivery"/>
        <SwitchOff v-if="!cashDelivery"/>
      </div>
    </div>
  </div>

</div>
</template>

<script>
import BackHeader from "@/components/layout/header/BackHeader.vue";
import {mapActions, mapGetters} from "vuex";
import $axios from "@/api";
import SwitchOn from "@/components/svg/SwitchOn.vue";
import SwitchOff from "@/components/svg/SwitchOff.vue";
export default {
name: "settings",
  components: {SwitchOff, SwitchOn, BackHeader},
   async mounted() {
    document.body.style.backgroundColor = "#F7F7F7";
    await this.syncSettings()
  },
  computed: {
    ...mapGetters({
      settings: 'user/settings'
    }),
    simpleMode() {
      return this.settings.app_mode === 0
    },
    cashDelivery() {
      return this.settings.delivery_type === 0
    }
  },
  methods:{
    ...mapActions({
      'syncSettings': 'user/SyncSettings',
    }),
    async switchVersion() {
      let appMode = this.settings.app_mode === 0 ? 1 : 0
      await $axios.post(`/v1/user-settings`,{
        app_mode: appMode,
        delivery_type: this.settings.delivery_type
      })
      await this.syncSettings()
    },
    async switchDeliveryType() {
      let deliveryType = this.settings.delivery_type === 0? 1 :0
      await $axios.post(`/v1/user-settings`,{
        app_mode: this.settings.app_mode,
        delivery_type: deliveryType
      })
      await this.syncSettings()
    }
  }
}
;</script>

<style scoped lang="scss">
.settings {
  padding: 24px;
  margin-top: 42px;
  h1{
    font-style: normal;
    font-weight: bold;
    font-size: 32px;
    line-height: 48px;
    margin-bottom: 18px;
    color:  #151516;;
  }
  .setting-item{
    height: 84px;
    box-sizing: border-box;
    padding: 20px 24px;
    display: flex;
    border-radius: 12px;
    align-items: center;
    background-color: #FFFFFF;
    margin-bottom: 18px;
    h2{
      font-style: normal;
      font-weight: bold;
      font-size: 18px;
      line-height: 27px;
      height: 27px;
      color: #151516;
    }

    p{
      height: 17px;
      font-style: normal;
      font-weight: 300;
      font-size: 11px;
      line-height: 16px;
      color: #88888B;
    }

    .switch{
      flex: 1;
      display: flex;
      justify-content: flex-end;
    }
  }
}

</style>