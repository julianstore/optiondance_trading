package mixin

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"option-dance/cmd/config"
	"option-dance/pkg/util"
	"testing"

	"github.com/MixinNetwork/mixin/common"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
)

var client *mixin.Client
var err error

func TestMain(m *testing.M) {
	config.InitTestConfig("dev")
	client, err = Client()
	if err != nil {
		log.Panicln(err)
	}
	m.Run()
}

//cnb transfer link: https://mixin.one/pay?amount=1e-7&asset=965e5c6e-434c-3fa9-b780-c50f43cd955c&recipient=443bf7c0-308f-461e-ac0f-88ef19b9162b&trace=94a6fb89-2bea-4982-a146-77357f6d788b&memo=donate:label
func TestSignAuthenticationToken(t *testing.T) {
	//Public transfer records
	//sid := "aebbc767-8622-4560-b7fd-068e0efc0af8"
	sid1 := "4f8c73cd-7113-4ebf-a93d-97286d28c314"
	res, err := MixinRequest("/network/snapshots/"+sid1, "GET", nil)
	println(res, err.Error())
}

func TestPayResultByTraceIdUseRobot(t *testing.T) {
	robot := PayResultByTraceIdUseRobot("bae3cdc6-e274-48fb-8b7b-ae1b4ed3e356")
	//robot := PayResultByTraceIdUseRobot("3bb4aefd-5a0b-4472-8f0f-754cbacd3303")
	println(robot)
}

func TestSendPlainTextMessage(t *testing.T) {
	//service.SendPlainTextMessage("019308a5-e1a9-427c-af2a-e05093beedaa", "hello")
}

var rawErr = `77770002b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d80001aa836a315472e227c0e7562230da94ea2f17d167ed684703535ccc11ff890914000100000000000000020000000416f2a24000013184c42f7d39707cc2fed9ac247c7d2f264d5754de6f2ab3a807810575281a8c5ed6bfbe8f82170631116dfdb7ab403287b867e3895ac21bcbc11f2c3949400e0003fffe0100000000000501ae0c29c0000323eedefe29b473e5931531cc54e0f2a9593eb10501aaca0b746c9ec983e41549df2f171810fa999ee9c57e09dfd190c1cc10a4a1cf545dfce3521e24205ab484bf65b4dd6007fd216734e774de10f1887dcf2d2535ae7752b8b14edc255318368f5d63279925ad33bc947605393124cb92b799ca20ced1cde8cdb408123aa57c0003fffe020000005c684b46427341414141414141414141414141414141414141414143685172414141414141414141414141414141414141414141416f552b77757678746165644852304f495774666f6c30436352614654706b4e42546b4e4654413d3d0001000200016afe5a59102d62b5c6f37f9680df9d83eae717146bd97ab6ac279137bd10eea65e738d253a5b8743ddcd26f3aa7b3509f7a38a83232bc990b8a824d28a472f0a00029fce2ec2770d38d0a64437e9b466871800a66b4d67edd6d52031181388a94411838de20b3808a80daf07c1a439e5c5f7bc83133ac18bbf8748c7a65bc25ae301`

func TestGetRawTransaction(t *testing.T) {
	raw := `77770002b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d8000169be5064df266feb5a743624834d184c5ac401638e715581efc8ede99011f1fb0001000000000000000200000005905af07a4000014b80fc7be44f5ae1d24ac476a5e2625e6bb46998aeb2ce9ff6c7702350e9fbb51f06e4abcdc14ecacc877def36ab97b3a6fecb34c3c1d9417268709936ce01520003fffe0100000000000585e09c73c00003a3400acc5bebdb5cb1226af328f6e9f8177c23acc5151bfb00fbe4ac4e9678f9290a51fb34cbf0fdb8de35e30fcc3dcccfcda3cbc9a8aaec924fce5db424df0a2f318979a9ec6359ea598fd733affdedb5e74ce596b43d4e7d6161ebae99a8329ee9d00578b53e0da55979bae0cb41baa4c0b31cf1268255c88f2c56c92156820003fffe020000005c684b46426f4b46436f4b46503267416b597a637a595745315a6a6b745a574d785a53307a4f546b324c546b314e4463744d6a67304e444a6b5a6d4e6d5a6d566a6f564f79554539545356524a5430356652566846556b4e4a55305645ffffff01e49f31822a3a44f45bf82ffefdcceed5c7cc77fb04dd7735713a29420eb1c04c248dd1d49b2b8be30e373ce676d484c6e45468cfe0d2017538a94a56abf2720500000105`
	var txHash mixin.Hash
	if tx, err := mixin.TransactionFromRaw(raw); err == nil {
		txHash, _ = tx.TransactionHash()
	}

	transaction, err := client.GetRawTransaction(context.Background(), txHash)
	if err != nil {
		t.Error(err)
	}
	id := mixinRawTransactionTraceID(txHash.String(), 0)
	println(id)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(transaction)
	util.JsonPrintS(string(marshal))
	log.Println(transaction)
}

func TestUnlockMutilsig(t *testing.T) {
	raw := "77770002b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d800021d2cd6fb695fa464b59388e9b4526a888282aab7156149d14ecc884d081040560000000000000000f9f725fd96d7bd20a97a224d6921358a5f92e8bcb95c2e2bf1462ff7276dd035000000000000000000020000000603fbcc48d5c000014a8bcc4a590a8a57c7e1b476734170069018efd9449973e437f469270db93dd521a8dc51d728dfd9f8b9ae208ff908d48b39e320949ea16d77565145a68982360003fffe0100000000000592aefc5e400003aa7ea71d64d0e9abc01864052bee92d78e1e752a7e718121a561e98c96502cb8cbc3f5a6bb9cd5b987a12262debc5aee3848f636cf44bf32756488bb6fb9790fb357cf6b9332ee903b9193c0aa4bf70d0124afb49f8f45250469a9e1df34cc55c51a3c8c7b978ce494143cb4f5eda287662472a3d7047c64b4defc25cc4e1d830003fffe020000005c684b46426f4b46436f4b46503267416b595449354f57526c595451744f5442694d43307a4d6a67794c546b79595451745a6a6b314e5459325a474d774e6a51316f564f79554539545356524a5430356652566846556b4e4a55305645ffffff01214600b046d69883ecef8c5909a455006e783338f54ac967aa0fdb609d7979d82f53d727d838e39c17e77aaa9d9d6f9f680cbe0530ee2de9569396181587290100000136"
	multisig, err := client.CreateMultisig(context.Background(), mixin.MultisigActionUnlock, raw)
	if err != nil {
		log.Println(err)
	}
	err = client.UnlockMultisig(context.Background(), multisig.RequestID, config.Cfg.DApp.Pin)
	println(err)
	if err != nil {
		log.Println(err)
	}
}

func mixinRawTransactionTraceID(hash string, index uint8) string {
	h := md5.New()
	_, _ = io.WriteString(h, hash)
	b := new(big.Int).SetInt64(int64(index))
	h.Write(b.Bytes())
	s := h.Sum(nil)
	s[6] = (s[6] & 0x0f) | 0x30
	s[8] = (s[8] & 0x3f) | 0x80
	sid, err := uuid.FromBytes(s)
	if err != nil {
		panic(err)
	}

	return sid.String()
}

func TestCreateMutilsig(t *testing.T) {
	//core.Init()
	//client, err := Client()
	//if err != nil {
	//	log.Panicln(err)
	//}
	//client.CreateMultisig(context.Background(),mixin.MultisigActionUnlock)
}

func TestSendTx(t *testing.T) {
	client, err := Client()
	if err != nil {
		log.Panicln(err)
	}
	raw := `77770002b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d800021d2cd6fb695fa464b59388e9b4526a888282aab7156149d14ecc884d081040560000000000000000f9f725fd96d7bd20a97a224d6921358a5f92e8bcb95c2e2bf1462ff7276dd035000000000000000000020000000603fbcc48d5c000014a8bcc4a590a8a57c7e1b476734170069018efd9449973e437f469270db93dd521a8dc51d728dfd9f8b9ae208ff908d48b39e320949ea16d77565145a68982360003fffe0100000000000592aefc5e400003aa7ea71d64d0e9abc01864052bee92d78e1e752a7e718121a561e98c96502cb8cbc3f5a6bb9cd5b987a12262debc5aee3848f636cf44bf32756488bb6fb9790fb357cf6b9332ee903b9193c0aa4bf70d0124afb49f8f45250469a9e1df34cc55c51a3c8c7b978ce494143cb4f5eda287662472a3d7047c64b4defc25cc4e1d830003fffe020000005c684b46426f4b46436f4b46503267416b595449354f57526c595451744f5442694d43307a4d6a67794c546b79595451745a6a6b314e5459325a474d774e6a51316f564f79554539545356524a5430356652566846556b4e4a55305645ffffff01214600b046d69883ecef8c5909a455006e783338f54ac967aa0fdb609d7979d82f53d727d838e39c17e77aaa9d9d6f9f680cbe0530ee2de9569396181587290100000136`

	transaction, err := client.SendRawTransaction(context.Background(), raw)
	if err != nil {
		println(err)
	}
	println(transaction)
}

func printRawJsonTx(r string) {
	raw, err := mixin.TransactionFromRaw(r)
	if err != nil {
		println(err)
	}
	marshal, err := json.Marshal(&raw)
	if err != nil {
		println(err)
	}
	log.Print(marshal)
}

func TestDumpTx(t *testing.T) {
	printRawJsonTx(`77770002b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d80001f12d7d02c69103767e2109155353e24f298439d4008886f67bca736d0a0e04df000000000000000000020000000411e1a30000017dbed3b798fb089027aebf6629c1c9da8475588ec4675afd692b3d018232fbf5b6a80231de055053f1180701b50abd7137f6747284646b92c6f0702c2f20326b0003fffe010000000000046553f1000003e8f42b3937dd09cae3b83b2022a92d927d1abff835288f339d75f08347333485721d19450b8a39262d2f149d4bc2a7082cb52d357d3041516c71fc28a9c8816f0142f40a6016c5eb56e960d4fa03cb0cabb5fc4e15e958b0a397ce06645859e69620ec690cdc6abaa6950bc9ea0f62eb908e53ab41320ea6114f1fc0270fb1e70003fffe0200000058684b46427342797762484b35443037456f354f35636a713737432b685172416c3333716b443264504c377231463979726a3978536f552b774141414141414141414141414141414141414141414b46547055314256454e490000`)
	printRawJsonTx(`77770002b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d800015b75793c43b8877a660707de1ca48a07b6c2f729d3f7316fe42c0e17dac5931d00000000000000000002000000040bebc20000012c4d228de4d8bf65de0529f5b68092cc148c3614bbf97a02c9d70d57fa9255ee05ca6adb68d304bcbfade21986528572f13af4fdd4dd4acbc6043c3c6e6adca80003fffe0100000000000405f5e1000003160681313e7daf7885e7d0bb90a09da11264138fa35defa672c7eb4f443b917a3777228579dd775a0d1424b6ec8dd1e6c70dd0953e502eff317620824a0cae5067998e146ee58f39abc3133c1595dd3a2d8bc244fcc88f947983488badd8636519eb2603e44d886a82c51624261a8197179e474e7232f17066e43917672fac470003fffe0200000058684b464273496d78452b326351554e6b694164394742384c6269696851724442387879764f6b524f355a59536b3839422b574b506f552b774141414141414141414141414141414141414141414b46547055314256454e490000`)
}

func TestCreateSig(t *testing.T) {
	ctx := context.Background()
	rawTx := "77770002b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d800021d2cd6fb695fa464b59388e9b4526a888282aab7156149d14ecc884d081040560000000000000000f9f725fd96d7bd20a97a224d6921358a5f92e8bcb95c2e2bf1462ff7276dd035000000000000000000020000000603fbcc48d5c000014a8bcc4a590a8a57c7e1b476734170069018efd9449973e437f469270db93dd521a8dc51d728dfd9f8b9ae208ff908d48b39e320949ea16d77565145a68982360003fffe0100000000000592aefc5e400003aa7ea71d64d0e9abc01864052bee92d78e1e752a7e718121a561e98c96502cb8cbc3f5a6bb9cd5b987a12262debc5aee3848f636cf44bf32756488bb6fb9790fb357cf6b9332ee903b9193c0aa4bf70d0124afb49f8f45250469a9e1df34cc55c51a3c8c7b978ce494143cb4f5eda287662472a3d7047c64b4defc25cc4e1d830003fffe020000005c684b46426f4b46436f4b46503267416b595449354f57526c595451744f5442694d43307a4d6a67794c546b79595451745a6a6b314e5459325a474d774e6a51316f564f79554539545356524a5430356652566846556b4e4a55305645ffffff01214600b046d69883ecef8c5909a455006e783338f54ac967aa0fdb609d7979d82f53d727d838e39c17e77aaa9d9d6f9f680cbe0530ee2de9569396181587290100000136"
	multisig, err := client.CreateMultisig(ctx, mixin.MultisigActionSign, rawTx)
	if err != nil {
		log.Println(err)
	}
	signMultisig, err := client.SignMultisig(ctx, multisig.RequestID, config.Cfg.DApp.Pin)
	println(signMultisig)
	transaction, err := client.SendRawTransaction(context.Background(), signMultisig.RawTransaction)
	if err != nil {
		println(err)
	}
	println(transaction)
}

func TestParseTx(t *testing.T) {
	data, err := hex.DecodeString(rawErr)
	if err != nil {
		println(err)
	}
	ver, err := common.UnmarshalVersionedTransaction(data)
	marshal := ver.Marshal()

	_, err = hex.DecodeString(string(marshal))
	if err != nil {
		println(err)
	}
	if len(ver.SignaturesMap) > 0 && len(ver.SignaturesMap[0]) < 1 {
		println(err)
	}
}
