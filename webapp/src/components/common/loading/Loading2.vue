<template>
  <div id="loading">
    <div class="lds-ring"><div></div><div></div><div></div><div></div></div>
  </div>
</template>

<script>
export default {
  name: "Loading2",

}
</script>

<style scoped>
#loading{
  position: fixed;
  background: rgba(255, 255, 255, 0.8);
  -webkit-backdrop-filter: blur(60px);
  backdrop-filter: blur(60px);
  margin: 0;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  transition: opacity .3s;
  z-index: 2000;
}
.lds-ring {
  position: fixed;
  margin: auto;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  width: 46px;
  height: 46px;
}
.lds-ring div {
  position: fixed;
  margin: 0;
  width: 41px;
  height: 41px;
  border: 5px solid;
  border-radius: 50%;
  animation: lds-ring 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
  border-color: #0366d6 transparent transparent transparent;
}
.lds-ring div:nth-child(1) {
  animation-delay: -0.3s;
}
.lds-ring div:nth-child(2) {
  animation-delay: -0.2s;
}
.lds-ring div:nth-child(3) {
  animation-delay: -0.10s;
}
@keyframes lds-ring {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
