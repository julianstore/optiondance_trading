<template>
    <div class="home-header">
        <img @click="openUserCenter" class="profile-icon" src="@/assets/svg/profile.svg"/>
    </div>
</template>

<script>
export default {
name: "HomeHeader",
    methods:{
        openUserCenter() {
            this.$parent.openUserCenter()
        }
    }
}
</script>

<style scoped lang="scss">

.home-header{
    .profile-icon{
        margin-right: 20px;
    }
    z-index: 1200;
    background-color: #F7F7F7;
    padding: 10px 24px;
    width: 100%;
    display: block;
    top:0;
    box-sizing: border-box;
    position: fixed;
}
</style>
