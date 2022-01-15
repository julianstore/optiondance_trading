<template>
    <div class="content">
        <h1>OptionDance<br/>{{ title }}</h1>
        <h4>OptionDance 是什么？</h4>
        <p class="mb-18">OptionDance 是Mixin生态内的期权交易市场，也是一个被重新设计的期权交易工具，解决了传统期权交易工具复杂难用的问题，普通人也可以轻松上手。</p>
        <img class="compare" src="../assets/img/compare.png"/>
        <h4>为什么可以打折买入比特币？</h4>
        <p class="mb-36">OptionDance上线的第一个版本帮助用户卖出现金担保看跌期权，卖出成功后用户将立即获得权利金收入。期权到期后如果对方交付比特币，资产的实际买入成本便是（行权价-权利金）。</p>
        <h4>实际案例？</h4>
        <p class="mb-18">1993年4月此时可口可乐的股价在40多美元，巴菲特以1.5美元的价格卖出500万股当年12月到期、行权价35美元的看跌期权，共收取权利金750万美元。</p>
        <img src="../assets/img/buffett.png"/>
        <h4>看不懂!</h4>
        <p class="pb-36 border-bottom">没关系，本来像资本家一样思考和投资就是件极难的事情，产品上线后你可以使用 OptionDance 来完成以前你无法学会的投资策略。</p>
        <div class="links">
            <h4>了解更多</h4>
            <a href="https://www.wyattresearch.com/article/warren-buffett-approach-to-selling-puts">Warren Buffett’s Approach to Selling Puts</a>
            <a href="https://www.nasdaq.com/articles/how-buffett-used-this-simple-strategy-to-boost-returns-2019-11-06">How Buffett Used This Simple Strategy To Boost Returns</a>
        </div>
    </div>
</template>

<script>
export default {
name: "ProductIntroduce",
    props:{
        title:String
    }
}
</script>

<style scoped lang="scss">
.content{
    padding: 0 24px 0;
    margin-top: 56px;
    h1 ,h4 ,h5 ,p{
        font-style: normal;
        font-weight: normal;
        text-align: left;
        color: #35384A;
    }
    h4{
        font-style: normal;
        font-weight: bold;
        font-size: 20px;
        line-height: 30px;
        letter-spacing: 0.03em;
        margin-bottom: 18px;
    }
    h5{
        font-size: 14px;
        line-height: 21px;
    }
    h1{
        font-style: normal;
        font-weight: bold;
        font-size: 32px;
        line-height: 48px;
        width: 252px;
        letter-spacing: 0.03em;
        margin-top: 16px;
        margin-bottom: 48px;
    }
    img{
        margin: 0  0 36px;
        display: block;
        width: 100%;
    }
    .compare{
        filter: drop-shadow(0px 4px 10px rgba(0, 0, 0, 0.12));
    }
    p{
        font-weight: 300;
        font-size: 16px;
        line-height: 26px;
        letter-spacing: 0.03em;
    }
    .links{
        margin-top: 36px;
        margin-bottom: 60px;
        h4{
            font-size: 20px;
            line-height: 30px;
            font-style: normal;
            font-weight: bold;
            margin-bottom: 14px;
        }
        a{
            font-size: 14px;
            line-height: 21px;
            background: url("../assets/svg/link.svg") no-repeat;
            padding-left: 19px;
            display: block;
            text-decoration: underline;
            color: #2F80ED;
            margin-bottom: 14px;
        }
    }
    .joined-success{
        margin-top: 24px;
        padding-bottom: 40px;
        margin-bottom: 40px;
        h1{
            font-style: normal;
            font-weight: bold;
            font-size: 42px;
            line-height: 63px;
            letter-spacing: 0.01em;
            color: #35384A;
        }
        a{
            font-style: normal;
            font-weight: bold;
            font-size: 14px;
            line-height: 21px;
            height: 26px;
            text-decoration-line: underline;
            color: #519AEC;
            display: block;
            padding: 2.5px 0;
        }
        .mixin-share{
            background: url("../assets/img/mixin.png") no-repeat;
            background-size: 26px;
            padding-left: 34px;
            margin-bottom: 25px;
        }
        .copy-share-link{
            background: url("../assets/img/link_big.png") no-repeat;
            background-size: 26px;
            padding-left: 34px;
        }
    }
}
.border-bottom{
    border-bottom: 0.5px solid rgba(149, 149, 151, 0.3);
}

</style>
