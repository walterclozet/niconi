package live

import (
	"elichika/config"
	"elichika/generic"
	"elichika/handler"
	"elichika/klab"
	"elichika/model"
	"elichika/userdata"
	"elichika/utils"

	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"xorm.io/xorm"
)

type DropWithTypes struct {
	StandardDrops   [0]model.RewardDrop `json:"standard_drops"`
	AdditionalDrops [0]model.RewardDrop `json:"additional_drops"`
	GimmickDrops    [0]model.RewardDrop `json:"gimmick_drops"`
}

type SkipLiveResult struct {
	LiveDifficultyMasterID        int                                              `json:"live_difficulty_master_id"`
	LiveDeckID                    int                                              `json:"live_deck_id"`
	Drops                         []DropWithTypes                                  `json:"drops"`
	MemberLoveStatuses            generic.ObjectByObjectIDWrite[*MemberLoveStatus] `json:"member_love_statuses"`
	GainUserExp                   int                                              `json:"gain_user_exp"`
	IsRewardAccessoryInPresentBox bool                                             `json:"is_reward_accessory_in_present_box"`
	ActiveEventResult             *int                                             `json:"active_event_result"`
	LiveResultMemberGuild         *int                                             `json:"live_result_member_guild"`
}

func LiveSkip(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type LiveSkipReq struct {
		LiveDifficultyMasterID int `json:"live_difficulty_master_id"`
		DeckID                 int `json:"deck_id"`
		TicketUseCount         int `json:"ticket_use_count"`
	}
	req := LiveSkipReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	utils.CheckErr(err)

	userID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, userID)
	defer session.Close()
	session.UserStatus.LastLiveDifficultyID = req.LiveDifficultyMasterID
	db := ctx.MustGet("masterdata.db").(*xorm.Engine)
	info := LiveFinishLiveDifficultyInfo{}

	exists, err := db.Table("m_live_difficulty").Where("live_difficulty_id = ?", req.LiveDifficultyMasterID).Get(&info)
	utils.CheckErrMustExist(err, exists)
	liveMemberMappingID := 0
	db.Table("m_live").Where("live_id = ?", info.LiveID).Cols("live_member_mapping_id").Get(&liveMemberMappingID)
	centerPositions := []int{}
	err = db.Table("m_live_member_mapping").Where("mapping_id = ? AND is_center = 1", liveMemberMappingID).
		Cols("position").Find(&centerPositions)
	utils.CheckErr(err)
	info.RewardCenterLovePoint = klab.CenterBondGainBasedOnBondGain(info.RewardBaseLovePoint) / len(centerPositions)

	skipLiveResult := SkipLiveResult{
		LiveDifficultyMasterID: req.LiveDifficultyMasterID,
		LiveDeckID:             req.DeckID,
		GainUserExp:            info.RewardUserExp * req.TicketUseCount}

	for i := 1; i <= req.TicketUseCount; i++ {
		skipLiveResult.Drops = append(skipLiveResult.Drops, DropWithTypes{})
	}
	session.UserStatus.Exp += skipLiveResult.GainUserExp
	deck := session.GetUserLiveDeck(req.DeckID)
	deckJsonByte, _ := json.Marshal(deck)
	cardMasterIDs := []int{}
	gjson.Parse(string(deckJsonByte)).ForEach(func(key, value gjson.Result) bool {
		if strings.Contains(key.String(), "card_master_id") {
			cardMasterIDs = append(cardMasterIDs, int(value.Int()))
		}
		return true
	})

	bondCardPosition := make(map[int]int)
	for i, cardMasterId := range cardMasterIDs {
		userCard := session.GetUserCard(cardMasterId)
		userCard.LiveJoinCount += req.TicketUseCount // count skip clear in pfp
		session.UpdateUserCard(userCard)
		// update member love point
		isCenter := (i+1 == centerPositions[0])
		isCenter = isCenter || ((len(centerPositions) > 1) && (i+1 == centerPositions[1]))
		addedBond := info.RewardBaseLovePoint
		if isCenter {
			addedBond += info.RewardCenterLovePoint
		}
		memberMasterID := klab.MemberMasterIDFromCardMasterID(cardMasterId)

		pos, exists := bondCardPosition[memberMasterID]
		// only use 1 card master id or an idol might be shown multiple times
		if !exists {
			memberLoveStatus := skipLiveResult.MemberLoveStatuses.AppendNew()
			memberLoveStatus.RewardLovePoint = addedBond
			memberLoveStatus.CardMasterID = cardMasterId
			bondCardPosition[memberMasterID] = skipLiveResult.MemberLoveStatuses.Length - 1
		} else {
			(*skipLiveResult.MemberLoveStatuses.Objects[pos]).RewardLovePoint += addedBond
		}
	}
	for memberMasterID, pos := range bondCardPosition {
		addedBond := session.AddLovePoint(memberMasterID,
			req.TicketUseCount*(*skipLiveResult.MemberLoveStatuses.Objects[pos]).RewardLovePoint)
		(*skipLiveResult.MemberLoveStatuses.Objects[pos]).RewardLovePoint = addedBond
	}

	if info.IsCountTarget { // counted toward target and profiles
		liveStats := session.GetUserLiveStats()
		idx := klab.LiveDifficultyTypeIndexFromLiveDifficultyID(req.LiveDifficultyMasterID)
		liveStats.LivePlayCount[idx] += req.TicketUseCount
		session.UpdateUserLiveStats(liveStats)
	}

	signBody := session.Finalize(handler.GetData("userModelDiff.json"), "user_model_diff")
	signBody, _ = sjson.Set(signBody, "skip_live_result", skipLiveResult)

	resp := handler.SignResp(ctx, signBody, config.SessionKey)
	// fmt.Println(resp)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
