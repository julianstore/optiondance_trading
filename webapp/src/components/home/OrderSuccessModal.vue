<template>
    <div class="modal-bg">
        <div class="order-success-modal">
            <img src="@/assets/svg/success.svg"/>
            <p>{{ $t("home.orderSuccess") }}</p>
            <div class="button-group">
                <button @click="ok"><span>{{ $t("global.know") }}</span></button>
                <button @click="viewDetail"><span>{{ $t("global.viewDetail") }}</span></button>
            </div>

        </div>
    </div>
</template>

<script>
export default {
    name: "OrderSuccessModal",
    methods:{
        ok() {
            this.$parent.closeSuccessModal();
        },
        viewDetail() {
          this.$parent.toOrderDetail();
        }
    }

}
</script>

<style scoped lang="scss">
.order-success-modal{
    width: 323px;
    height: 274px;
    background: #FFFFFF;
    box-sizing: border-box;
    border-radius: 12px;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    margin: auto;
    padding: 16px;
    img{
        width: 72px;
        height: 72px;
        margin: 28px auto 14px;
        display: block;
    }
    p{
        text-align: center;
        font-style: normal;
        font-weight: bold;
        font-size: 17px;
        line-height: 25px;
        color: #151516;
        margin-bottom: 48px;
    }
    .button-group{
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
