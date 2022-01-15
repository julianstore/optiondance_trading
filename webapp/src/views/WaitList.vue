<template>
  <div class="waitlist">
    <div class="logo-anim mb-30">
      <div id="lottie-logo"></div>
    </div>
    <p class="mb-30">
      现在加入WAITLIST，我们将在开发完成后<br />第一时间向你推送体验版
    </p>
    <button class="waitlist-btn mb-32" @click="joinWaitList">
      <span>WAITLIST</span>
    </button>
    <button class="product-intro-btn mb-16" @click="toProductIntro">
      <span>✨ 了解 OptionDance ✨</span>
    </button>
  </div>
</template>

<script>
import $axios from "@/api";
import { mapActions, mapGetters } from "vuex";
import { getPlatformType, removeRegionCode } from "@/util/utils";
import Loading2 from "@/components/common/loading/Loading2.vue";
import CoinTextRoll from "@/components/common/CoinTextRoll.vue";
import dataJson from "@/assets/lottie/data.json";
import lottie from "lottie-web";
import { PlatformType } from "@/util/constants";

export default {
  name: "WaitList",
  components: { CoinTextRoll, Loading2 },
  mounted() {
    this.loadAnim();
    let client_id = import.meta.env.VITE_MIXIN_CLIENT_ID;
    this.mixinAuthUrl = `https://mixin-www.zeromesh.net/oauth/authorize?client_id=${client_id}&scope=PROFILE:READ+PHONE:READ&response_type=code&state=login&return_to=joinSuccess`;
    let wid = this.$route.query.wid;
    let referer = this.$route.query.referer;
    this.waitlist.type = getPlatformType();
    console.log(this.$route.query);
    if (wid || referer !== "auth") {
      this.waitlist.inviterWid = Number(wid);
    }
    this.setWaitlist(this.waitlist);
  },
  computed: {
    ...mapGetters({
      mixinToken: "user/mixinToken",
      mixinUser: "user/mixinUser",
      waitlist: "user/waitlist",
    }),
  },
  data() {
    return {
      mixinAuthUrl: "",
      gifLoaded: false,
      coinTextStartRoll: false,
    };
  },
  methods: {
    ...mapActions({
      setMixinUser: "user/SetMixinUser",
      setWaitlist: "user/SetWaitlist",
    }),
    async joinWaitList() {
      //mixin env
      if (this.waitlist.type === PlatformType.MIXIN) {
        this.$loading.show();
        let res = await $axios.get("https://mixin-api.zeromesh.net/me", {
          headers: { Authorization: "Bearer " + this.mixinToken },
        });
        if (res.error || !res.data.identity_number) {
          window.location.href = this.mixinAuthUrl;
          return;
        }
        let user = res.data;
        this.setMixinUser(user);
        let data = await $axios.post("/v1/waitlist", {
          mixin_id: user.identity_number,
          mixin_uid: user.user_id,
          mixin_phone: removeRegionCode(user.phone),
          mixin_name: user.full_name,
          inviter_wid: this.waitlist.inviterWid,
          type: this.waitlist.platformType,
        });
        if (data.code === 0) {
          this.$router.push({ name: "join-success" });
          this.$toast.success("成功加入waitlist");
        } else {
          this.$toast.error(data.msg);
        }
        this.$loading.hide();
      } else {
        //Other environment
        await this.$router.push({ name: "subscribe-email" });
      }
    },
    async toProductIntro() {
      await this.$router.push({ name: "product-intro" });
    },
    async gifLoadedHandler() {
      this.gifLoaded = true;
      this.coinTextStartRoll = true;
    },
    loadAnim() {
      lottie.loadAnimation({
        container: document.getElementById("lottie-logo"),
        renderer: "svg",
        loop: true,
        autoplay: true,
        animationData: dataJson,
      });
    },
  },
};
</script>

<style scoped lang="scss">
.waitlist {
  margin-top: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  @media screen and (max-height: 670px) {
    .logo-anim {
      max-width: 250px;
    }
  }
  @media screen and (min-height: 670px) {
    .logo-anim {
      max-width: 300px;
    }
  }
  h1,
  p,
  span {
    font-family: Hiragino Sans GB, Source Han Sans, pingfang SC;
    font-style: normal;
    font-weight: normal;
    text-align: center;
  }
  a {
    text-decoration: none;
    border-bottom: 1px solid;
  }
  h1 {
    font-weight: bold;
    font-size: 35px;
    line-height: 52px;
    letter-spacing: 0.02em;
    color: #35384a;
  }
  p {
    width: 283px;
    height: 42px;
    font-size: 14px;
    line-height: 21px;
    color: #35384a;
  }
  .waitlist-btn {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    padding: 14px 50px;
    width: 209px;
    height: 58px;
    background: #ff6f61;
    box-shadow: 0px 6px 12px rgba(255, 111, 97, 0.25);
    border-radius: 54px;
    span {
      font-weight: bold;
      font-size: 18px;
      line-height: 27px;
      text-align: center;
      color: #ffffff;
    }
  }
  .product-intro-btn {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    padding: 8px 4px;
    width: 183px;
    height: 36px;
    background: #35384a;
    border-radius: 4px;
    span {
      font-size: 14px;
      //line-height: 20px;
      color: #ffffff;
    }
  }
  .coinTextRoll {
    display: inline-block;
    letter-spacing: 0;
  }
}
</style>
