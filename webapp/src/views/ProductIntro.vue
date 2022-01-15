<template>
<div class="product-intro">
  <BackHeader back-to="/waitlist"/>
  <ProductIntroduce :title="this.$t('faq.productIntro')"/>
</div>
</template>

<script>
import BackHeader from "@/components/layout/header/BackHeader.vue";
import ProductIntroduce from "@/components/ProductIntroduce.vue";

export default {
name: "ProductIntro",
  components: {ProductIntroduce, BackHeader},
}
</script>

<style scoped lang="scss">
</style>
