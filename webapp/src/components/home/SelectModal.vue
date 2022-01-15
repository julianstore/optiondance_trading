<template>
    <div class="modal-bg">
        <div class="select-modal" :id="type">
            <div class="title-item">
                <img v-show="type===selectModalType.EXERCISE_DATE_MODAL" src="@/assets/svg/date-icon.svg"/>
                <img v-show="type===selectModalType.EXERCISE_PRICE_MODAL" src="@/assets/svg/stock-icon.svg"/>
                <span>{{ title }}</span>
            </div>
            <div class="scroll-area">
                <div tabindex="0"  v-for="(v,index) in items" class="select-item"
                     v-bind:class= "{'selected': v === selectItem }">
                    <span>{{ v }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import {SelectModalType} from "@/util/constants";

export default {
name: "SelectModal",
    props:{
        title:String,
        items:Array,
        type: String,
        selectItem: String
    },
    mounted() {
        this.selectModalType = SelectModalType
    },
    data() {
        return {
            selectModalType:{}
        }
    },
}
</script>

<style scoped lang="scss">
    .select-modal{
        position: absolute;
        width: 327px;
        height: 374px;
        left: 0;
        right: 0;
        top: 0;
        bottom: 0;
        margin: auto;
        background: #FFFFFF;
        border-radius: 12px;
        padding: 16px 16px 32px;
        box-sizing: border-box;
        .title-item{
            width: 295px;
            height: 62px;
            background: #F7F7F7;
            border-radius: 12px;
            display: flex;
            justify-content: left;
            align-items: center;
            padding: 18px 16px;
            box-sizing: border-box;
            img{
                width: 28px;
                height: 28px;
                margin-right: 8px;
            }
            span{
                font-style: normal;
                font-weight: bold;
                font-size: 17px;
                line-height: 25px;
                color: #151516;
            }
        }
        .scroll-area{
            width: 295px;
            height: 260px;
            margin-top: 18px;
            background: #FFFFFF;
            overflow:auto;
            .select-item{
                outline: 0;
                height: 54px;
                border-radius: 12px;
                display: flex;
                justify-content: left;
                align-items: center;
                padding: 18px 16px;
                box-sizing: border-box;
                span{
                    font-style: normal;
                    font-weight: 300;
                    font-size: 15px;
                    line-height: 22px;
                    letter-spacing: 0.02em;
                    color: #151516;
                }
                &:focus,&:hover{
                    background: #F7F7F7;
                }
            }
        }
    }
    .selected{
        background: #F7F7F7;
    }
</style>
