<template>
<div class="bottom-modal">
    <div class="modal-bg" v-show="visible"></div>
    <transition name="normal">
        <div v-show="visible" class="bottom-modal-content" >
            <div class="title">
                <span >{{ $t("userCenter.userCenter") }}</span>
                <img @click="close" src="@/assets/svg/close.svg"/>
            </div>
            <div class="option-position" @click="toPosition">
                <img src="@/assets/svg/avatar.svg"/>
                <span>{{ $t("userCenter.position") }}</span>
                <img class="left-arrow-icon" src="@/assets/svg/arrow_left_big.svg"/>
            </div>
            <div>
              <div class="menu-item" @click="toIncome">
                <div class="item-content">
                  <img src="@/assets/svg/userCenter/income.svg"/>
                  <span>{{ $t("userCenter.income") }}</span>
                  <img class="left-arrow-icon" src="@/assets/svg/arrow_left_big.svg"/>
                </div>
              </div>
              <div class="menu-item" @click="toSettings">
                <div class="item-content">
                <img src="@/assets/svg/userCenter/settings.svg"/>
                <span>{{ $t("userCenter.setting") }}</span>
                <img class="left-arrow-icon" src="@/assets/svg/arrow_left_big.svg"/>
                </div>
              </div>
                <div class="menu-item" @click="toFaq">
                  <div class="item-content">
                    <img src="@/assets/svg/faq.svg"/>
                    <span>{{ $t("userCenter.commonProblem") }}</span>
                    <img class="left-arrow-icon" src="@/assets/svg/arrow_left_big.svg"/>
                  </div>
                </div>
                <div class="menu-item" @click="joinGroup">
                  <div class="item-content">
                    <img src="@/assets/svg/group.svg"/>
                    <span>{{ $t("userCenter.joinGroup") }}</span>
                    <img class="left-arrow-icon" src="@/assets/svg/arrow_left_big.svg"/>
                  </div>
                </div>
            </div>
        </div>
    </transition>
</div>
</template>

<script>
export default {
name: "UserCenter",
    props:{
        visible:Boolean,
    },
    methods:{
        close() {
            this.$parent.closeUserCenter()
        },
        toFaq() {
            this.close()
            this.$router.push('/faq')
        },
        toSettings() {
          this.close()
          this.$router.push('/settings')
        },
        toIncome() {
          this.close()
          this.$router.push('/income')
        },
        joinGroup() {
            this.close()
            window.open('mixin://apps/a757520f-cfd3-4a4b-a1c8-85647e908ee0')
        },
      toPosition() {
        this.close()
        this.$router.push('/options')
      }
    }
}
</script>

<style scoped lang="scss">
.bottom-modal-content{
  .title{
    padding: 0 12px;
  }
    .option-position{
        width: 100%;
        height: 80px;
        background: #FFFFFF;
        box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.08);
        border-radius: 12px;
        display: flex;
        align-items: center;
        padding: 18px 20px;
        box-sizing: border-box;
        img{
            margin-right: 18px;
        }
        span{
            font-style: normal;
            font-weight: bold;
            font-size: 18px;
            line-height: 25px;
            color: #151516;
        }
    }
    .menu-item{
        padding:0 16px;
        &:first-child{
          margin-top: 16px;
        }
        &:last-child .item-content{
            border-bottom: none;
        }
        .item-content{
          height: 76px;
          display: flex;
          align-items: center;
          box-sizing: border-box;
          border-bottom: 0.5px solid #EAEAEA;
          img{
            margin-right: 18px;
          }
          span{
            font-style: normal;
            font-weight: bold;
            font-size: 18px;
            line-height: 25px;
            color: #151516;
          }
        }
    }
}
</style>
