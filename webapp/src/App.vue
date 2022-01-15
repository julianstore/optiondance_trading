<template>
  <main>
    <router-view />
    <PreLoad />
  </main>
</template>

<script>
import PreLoad from "@/components/common/PreLoad.vue";
export default {
  name: "App",
  components: { PreLoad },
  created() {
    let loading = document.getElementById("loading");
    if (loading !== null) {
      setTimeout(() => {
        document.body.removeChild(loading);
      }, 300); //Delay 500ms to ensure smooth animation
    }
  },
};
</script>
