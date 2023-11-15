package handler

import (
	"elichika/encrypt"
	"elichika/locale"
	"elichika/utils"

	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	presetDataPath = "assets/preset/"
	userDataPath   = "assets/userdata/"
)

func init() {
	os.Mkdir(userDataPath, 0755)
}

func SignResp(ctx *gin.Context, body, key string) (resp string) {
	ep := ctx.MustGet("ep").(string)
	masterVersion := ctx.MustGet("locale").(*locale.Locale).MasterVersion
	signBody := fmt.Sprintf("%d,\"%s\",0,%s", time.Now().UnixMilli(), masterVersion, body)
	sign := encrypt.HMAC_SHA1_Encrypt([]byte(ep+" "+signBody), []byte(key))

	resp = fmt.Sprintf("[%s,\"%s\"]", signBody, sign)
	return
}

func GetData(fileName string) string {
	presetDataFile := presetDataPath + fileName
	if !utils.PathExists(presetDataFile) {
		panic("File not exists")
	}

	return utils.ReadAllText(presetDataFile)
}
