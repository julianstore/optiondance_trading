<template>
  <div class="product-intro">
    <BackHeader back-to="-1" bg-color="#FFFFFF"/>
    <div class="article content">
      <h1>{{  $t("faq.whatIsOptionDanceTitle") }}</h1>
      <p>
        {{ $t("faq.odDesc") }}
      </p>
      <img class="compare" src="@/assets/img/compare.png"/>
      <div class="border-bottom"></div>
    </div>
    <join-group></join-group>
  </div>
</template>

<script>
import BackHeader from "@/components/layout/header/BackHeader.vue";
import JoinGroup from "@/components/faq/JoinGroup.vue";
import {setHtmlMeta} from "@/util/utils";

export default {
  name: "OptionDanceIntro",
  components: {JoinGroup, BackHeader},

  beforeMount() {
    setHtmlMeta("theme-color","#FFFFFF")
    document.body.style.backgroundColor = "#FFFFFF";
  }
}
</script>

<style scoped lang="scss">
.content {
  h1{
    width: 252px;
  }
  padding: 0 30px;
  margin: 70px 0 32px;

  .border-bottom {
    height: 48px;
    border-bottom: 0.5px solid rgba(136, 136, 139, 0.3);
  }

}

.compare {
  filter: drop-shadow(0px 4px 10px rgba(0, 0, 0, 0.12));
  margin-bottom: 0;
}
</style>