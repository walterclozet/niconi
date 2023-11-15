package handler

import (
	"elichika/config"
	"elichika/userdata"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func FetchEmblem(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, UserID)
	defer session.Close()
	signBody := session.Finalize(GetData("fetchEmblem.json"), "user_model")
	resp := SignResp(ctx, signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ActivateEmblem(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	UserID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, UserID)
	defer session.Close()
	var emblemId int64
	gjson.Parse(reqBody).ForEach(func(key, value gjson.Result) bool {
		if value.Get("emblem_master_id").String() != "" {
			emblemId = value.Get("emblem_master_id").Int()
			session.UserStatus.EmblemID = int(emblemId)
			return false
		}
		return true
	})

	signBody := session.Finalize(GetData("activateEmblem.json"), "user_model")

	resp := SignResp(ctx, signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
