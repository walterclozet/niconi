package handler

import (
	"elichika/config"
	"elichika/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func FetchCommunicationMemberDetail(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	var memberId int64
	gjson.Parse(reqBody).ForEach(func(key, value gjson.Result) bool {
		if value.Get("member_id").String() != "" {
			memberId = value.Get("member_id").Int()
			return false
		}
		return true
	})

	lovePanelCellIds := []int{}
	err := MainEng.Table("m_member_love_panel_cell").
		Join("LEFT", "m_member_love_panel", "m_member_love_panel_cell.member_love_panel_master_id = m_member_love_panel.id").
		Cols("m_member_love_panel_cell.id").Where("m_member_love_panel.member_master_id = ?", memberId).
		OrderBy("m_member_love_panel_cell.id ASC").Find(&lovePanelCellIds)
	CheckErr(err)

	now := time.Now()
	year, month, day := now.Year(), now.Month(), now.Day()
	tomorrow := time.Date(year, month, day+1, 0, 0, 0, 0, now.Location()).Unix()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	signBody := GetData("fetchCommunicationMemberDetail.json")
	signBody, _ = sjson.Set(signBody, "member_love_panels.0.member_id", memberId)
	signBody, _ = sjson.Set(signBody, "member_love_panels.0.member_love_panel_cell_ids", lovePanelCellIds)
	signBody, _ = sjson.Set(signBody, "weekday_state.weekday", weekday)
	signBody, _ = sjson.Set(signBody, "weekday_state.next_weekday_at", tomorrow)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func UpdateUserCommunicationMemberDetailBadge(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	var memberMasterId int64
	gjson.Parse(reqBody).ForEach(func(key, value gjson.Result) bool {
		if value.Get("member_master_id").String() != "" {
			memberMasterId = value.Get("member_master_id").Int()
			return false
		}
		return true
	})

	userDetail := []any{}
	userDetail = append(userDetail, memberMasterId)
	userDetail = append(userDetail, model.UserCommunicationMemberDetailBadgeByID{
		MemberMasterID: int(memberMasterId),
	})

	signBody := GetData("updateUserCommunicationMemberDetailBadge.json")
	signBody, _ = sjson.Set(signBody, "user_model.user_status", GetUserStatus())
	signBody, _ = sjson.Set(signBody, "user_model.user_communication_member_detail_badge_by_id", userDetail)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func UpdateUserLiveDifficultyNewFlag(ctx *gin.Context) {
	signBody, _ := sjson.Set(GetData("updateUserLiveDifficultyNewFlag.json"),
		"user_model.user_status", GetUserStatus())
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FinishUserStorySide(ctx *gin.Context) {
	signBody, _ := sjson.Set(GetData("finishUserStorySide.json"),
		"user_model.user_status", GetUserStatus())
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FinishUserStoryMember(ctx *gin.Context) {
	signBody, _ := sjson.Set(GetData("finishUserStoryMember.json"),
		"user_model.user_status", GetUserStatus())
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SetTheme(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	// fmt.Println(reqBody)

	var memberMasterId, suitMasterId, backgroundMasterId int64
	gjson.Parse(reqBody).ForEach(func(key, value gjson.Result) bool {
		if value.Get("member_master_id").String() != "" {
			memberMasterId = value.Get("member_master_id").Int()
			suitMasterId = value.Get("suit_master_id").Int()
			backgroundMasterId = value.Get("custom_background_master_id").Int()

			gjson.Parse(GetUserData("memberSettings.json")).Get("user_member_by_member_id").
				ForEach(func(kk, vv gjson.Result) bool {
					if vv.IsObject() {
						if vv.Get("member_master_id").Int() == memberMasterId {
							SetUserData("memberSettings.json", "user_member_by_member_id."+
								kk.String()+".custom_background_master_id", backgroundMasterId)
							SetUserData("memberSettings.json", "user_member_by_member_id."+
								kk.String()+".suit_master_id", suitMasterId)
							return false
						}
					}
					return true
				})
			return false
		}
		return true
	})

	userMemberRes := []any{}
	userMemberRes = append(userMemberRes, memberMasterId)
	userMemberRes = append(userMemberRes, model.UserMemberInfo{
		MemberMasterID:           int(memberMasterId),
		CustomBackgroundMasterID: int(backgroundMasterId),
		SuitMasterID:             int(suitMasterId),
		LovePoint:                13181880,
		LovePointLimit:           13181880,
		LoveLevel:                500,
		ViewStatus:               1,
		IsNew:                    false,
	})

	userSuitRes := []any{}
	userSuitRes = append(userSuitRes, suitMasterId)
	userSuitRes = append(userSuitRes, model.SuitInfo{
		SuitMasterID: int(suitMasterId),
		IsNew:        false,
	})

	signBody := GetData("setTheme.json")
	signBody, _ = sjson.Set(signBody, "user_model.user_status", GetUserStatus())
	signBody, _ = sjson.Set(signBody, "user_model.user_member_by_member_id", userMemberRes)
	signBody, _ = sjson.Set(signBody, "user_model.user_suit_by_suit_id", userSuitRes)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SetFavoriteMember(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")

	SetUserData("userStatus.json", "favorite_member_id",
		gjson.Parse(reqBody).Array()[0].Get("member_master_id").Int())
	signBody, _ := sjson.Set(GetData("setFavoriteMember.json"),
		"user_model.user_status", GetUserStatus())
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
