package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MixinNetwork/go-number"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Interface2String(inter interface{}) (res string) {
	switch inter.(type) {
	case string:
		res = inter.(string)
		break
	}
	return
}

//string to int64
func String2Int64(s string) (res int64) {
	i, _ := strconv.Atoi(s)
	return int64(i)
}

func Interface2Int64(inter interface{}) int64 {
	var res int64
	switch inter.(type) {
	case int64:
		res = inter.(int64)
		break
	}
	return res
}

func Interface2Float64(inter interface{}) float64 {
	var res float64
	switch inter.(type) {
	case float64:
		res = inter.(float64)
		break
	}
	return res
}

func GetFileExNameWithDot(filename string) string {
	lastPointIndex := strings.LastIndex(filename, ".")
	return filename[lastPointIndex:]
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func GetTpltString(tpltPath string, params map[string]interface{}) (result string) {
	file, err := os.Open(path.Join(AbsPath(), "../../template/") + "/" + tpltPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	tpltStr, err := ioutil.ReadAll(file)
	tmpl, _ := template.New("tmpl1").Parse(string(tpltStr))
	buf := new(bytes.Buffer)
	_ = tmpl.Execute(buf, params)
	result = buf.String()
	return
}

func AbsPath() string {
	var absPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		absPath = path.Dir(filename)
	}
	return absPath
}

func GetPages(total int64, size int64) (pages int64) {
	if total%size == 0 {
		pages = total / size
	} else {
		pages = total/size + 1
	}
	return
}

func JsonPrint(v interface{}) {
	body, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", body)
}

func JsonPrintS(uglyBody string) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(uglyBody), "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", out.String())
}

func UuidNewV4() uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		zap.L().Panic("error", zap.Error(err))
	}
	return id
}

func LogGoNumberIntegers(integer ...number.Integer) {
	var s string
	for _, e := range integer {
		s += e.Persist() + " "
	}
	log.Println(s)
}
