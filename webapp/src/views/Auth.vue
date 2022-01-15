<template>
<div></div>
</template>

<script>
import {mapActions} from "vuex";
import axios from "@/api";
export default {
name: "Auth",
  async mounted() {
    let code = this.$route.query.code
    let action = this.$route.query.action
    let client_id = import.meta.env.VITE_MIXIN_CLIENT_ID
    let client_secret = import.meta.env.VITE_MIXIN_CLIENT_SECRET
    if (code) {
      let res = await axios.post('https://mixin-api.zeromesh.net/oauth/token', {
        "client_id": client_id,
        "code": code,
        "client_secret": client_secret
      });
      console.log(res)
      if (!res.error) {
        await this.MixinLogin({token: res.data.access_token})
        let returnTo = this.$route.query.return_to;
        if (returnTo) {
          await this.$router.push({
            path:returnTo,
            query:{referer:'auth'}
          })
        }
        }
      }
      else if(action === 'auth') {
      let client_id = import.meta.env.VITE_MIXIN_CLIENT_ID
      window.location.href = `https://mixin-www.zeromesh.net/oauth/authorize?client_id=${client_id}&scope=PROFILE:READ+PHONE:READ&response_type=code&state=login&return_to=home`
      }
    },
  methods: {
    ...mapActions({
      'MixinLogin': 'user/MixinLogin',
    }),
  }
}
</script>

<style scoped>

</style>
