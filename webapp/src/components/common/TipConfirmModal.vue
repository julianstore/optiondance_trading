<template>
    <div class="modal-bg">
        <div :style="{height: height+'px'}" class="tip-confirm-modal">
            <p class="title">{{ $t("global.operationConfirm") }}</p>
            <p>{{ tip }}</p>
            <div class="button-group">
                <button @click="cancel"><span>{{ $t("global.cancel") }}</span></button>
                <button @click="confirm"><span>{{ $t("global.confirm") }}</span></button>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "TipConfirmModal",
    props:{
      tip:String,
      height:Number
    },
    methods:{
      cancel() {
            this.$parent.TipModalCancel();
        },
       confirm() {
         this.$parent.TipModalConfirm();
      }
    }

}
</script>

<style scoped lang="scss">
.tip-confirm-modal{
    width: 323px;
    background: #FFFFFF;
    box-sizing: border-box;
    border-radius: 12px;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    margin: auto;
    padding: 24px 16px 16px;
    .title{
        text-align: center;
        font-style: normal;
        font-weight: bold;
        font-size: 17px;
        line-height: 25px;
        color: #151516;
        margin-bottom: 16px;
    }
      p{
        font-style: normal;
        font-weight: normal;
        font-size: 16px;
        line-height: 24px;
        text-align: center;
        color: #151516;
      }
    .button-group{
        position: absolute;
        bottom: 16px;
        display: flex;
        flex-direction: row;
        justify-content: center;
        button{
            display: flex;
            justify-content: center;
            align-items: center;
            width: 140px;
            height: 54px;
            box-sizing: border-box;
            border: 1px solid #151516;
            border-radius: 8px;
            span{
                font-style: normal;
                font-weight: 600;
                font-size: 17px;
                line-height: 25px;
            }
            &:nth-child(1){
                background: #FFFFFF;
                margin-right: 12px;
                span{
                    color: #151516;
                }
            }
            &:nth-child(2){
                background: #111111;
                color: #FFFFFF;
            }
        }
    }

}
</style>
