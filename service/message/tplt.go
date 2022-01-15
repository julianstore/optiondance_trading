package message

import (
	"bytes"
	"github.com/dustin/go-humanize"
	"strconv"
	"text/template"
	"time"
)

const (
	TpltGreet        = "你好，欢迎使用 OptionDance 让你的资金更有效率，首次使用可先看视频教程 :)"
	TpltVideoBtnText = "视频演示如何使用OptionDance低价买入BTC"
	VideoLink        = "https://option.dance/file/how-to-buy-bitcoin-at-low-prices-in-optiondance.mp4"
	BtnColor         = "#0066FF"
	ContractId       = "a757520f-cfd3-4a4b-a1c8-85647e908ee0"

	TpltExpiryNotify = `【到期提醒】
🎯 保卖：%s $%d 卖出 %s
💎 数量：%s %s
✨ 该标的即将到期，（BTC结算用户）如需操作行权请进入持仓进行，如不需要行权可忽略此条消息。`

	TitleTradeConfirmed           = "撮合成功"
	TitleOrderCancelled           = "撤单成功"
	TitleOrderFilled              = "下单成功"
	TitleOrderInvalid             = "订单无效"
	TitlePositionExercised        = "到期结束"
	TitlePositionExercisedRefund  = "到期结束"
	TitlePositionExercisedInvalid = "行权无效"
	TitlePositionClosedRefund     = "退回购币金"
	TitlePositionCashDeliveryEarn = "获得购币金"

	TpltCreateOrderGroupNtf = `【来新订单啦】

{{ sideTypeCn .Side .OptionType }}：{{ .ExpirationDate | timeDate }} $ {{ .StrikePrice | commaStyle  }} {{ sideCn .Side .OptionType }} {{.BaseCurrency}}

{{.Side}}Size： {{.RemainingAmount}} {{.BaseCurrency}}

{{.Side}}Price({{.QuoteCurrency}})：{{.Price}}`
)

var funcMap template.FuncMap = map[string]interface{}{
	"timeDate": func(input time.Time) string {
		return input.Format("2006年01月02日")
	},
	"sideTypeCn": func(side, opType string) string {
		return getTypeCn(side, opType)
	},
	"commaStyle": func(price string) string {
		sp, _ := strconv.Atoi(price)
		return humanize.Comma(int64(sp))
	},
	"sideCn": func(side, opType string) string {
		return getSideString(side, opType)
	},
}

func GetTplt(tplt string, data interface{}) (string, error) {
	t := template.New("tplt").Funcs(funcMap)
	parse, _ := t.Parse(tplt)
	var buf bytes.Buffer
	err := parse.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
