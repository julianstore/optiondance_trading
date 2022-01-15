package core

import (
	"github.com/MixinNetwork/go-number"
	"testing"
)

var pack = MsgPack{
	T:  "day0",
	AC: "EXPIRING_NOTIFY",
}

func TestMsgPack_Pack(t *testing.T) {

	memo, err := pack.Pack()
	if err != nil {
		t.Error(err)
	}
	t.Logf("memo: %s", memo)
}

func TestDecryptAction(t *testing.T) {
	m1 := "jqFBojIwokFDrENSRUFURV9PUkRFUqFCoKJCQ6NCVEOhSbNCVEMtMTFBVUcyMS0zODAwMC1QoU2goU+goVCmMTAwLjAwoVGgolFDpFVTRFShU6FCoVShTKJUQ9oAJDQ0NzViZTZjLTc4MmItNGMxNS1hZDAyLTYxOWEwOTVkYTNkN6FVoA=="
	m2 := "jqFBoKJBQ69FWFBJUklOR19OT1RJRlmhQqCiQkOgoUmgoU2goU+goVCgoVGgolFDoKFToKFUpGRheTCiVEOgoVWg"
	action1, _ := DecryptAction(m1)
	t.Log(action1.AC)
	t.Log(action1.T)
	remainingFunds := number.FromString("2000").Integer(8)
	t.Log(remainingFunds.Persist())
	price := number.FromString("10.39").Integer(24)
	t.Log(price.Persist())
	action2, _ := DecryptAction(m2)
	t.Log(action2.AC)
	t.Log(action2.T)

}

func TestGoNumber(t *testing.T) {
	i2 := number.FromString("99999.213").Integer(4)
	i1 := number.FromString("99999.12").Integer(4)
	mul := i2.Mul(i1)
	t.Log(mul.Persist())
	number.FromString("2000").Integer(8)
}

func TestGoNumber1(t *testing.T) {
	fromString := number.FromString("")
	t.Log(fromString.String())
}
