# OptionDance

## 标的物
Optiondance的标的物以 ${Currency}-${day}${month}${year}-${strike}-${C|P}的格式命名

> BTC-9MAY21-40000-P: 表示到期日为2021年5月9日，行权价为40000USDT，标的资产为BTC的PUT
> 
> ETH-12NOV21-35000-C: 表示到期日为2021年11月12日，行权价为35000USDT，标的资产为ETH的CALL


## 行权

跟欧式期权类似，OptionDance行权只能在到期日的北京时间结算时刻的前8小时内进行

到期日当天，前一天和两天，会通过mixin messenger给用户推送消息，提醒用户是否需要行权

未撮合的订单在标的到期日15:30时，进行系统自动撤单。


## 期权到期的结算安排

optiondance期权的到期时刻为到期日当天北京时间16：15，到期日当天北京时间16：30会进行结算，会按照如下指派规则进行履约结算

## 期权到期的履约指派

按比例配对：在按比例配对的履约指派方式下，持有多头头寸的清算会员的多单将会按比例分配给所有持有空头头寸的清算会员进行行权配对，

分配比例是按照空头头寸的清算会员持有的空单占该合约所有空单头寸总和的比例决定。