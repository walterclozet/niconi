package live

import (
	"elichika/config"
	"elichika/enum"
	"elichika/generic"
	"elichika/handler"
	"elichika/klab"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"xorm.io/xorm"
)

type LiveFinishCard struct {
	CardMasterID        int `json:"-"`
	GotVoltage          int `json:"got_voltage"`
	SkillTriggeredCount int `json:"skill_triggered_count"`
	AppealCount         int `json:"appeal_count"`
}

func (obj *LiveFinishCard) SetID(id int64) {
	obj.CardMasterID = int(id)
}

type MemberLoveStatus struct {
	CardMasterID    int `json:"-"`
	RewardLovePoint int `json:"reward_love_point"`
}

func (obj *MemberLoveStatus) ID() int64 {
	return int64(obj.CardMasterID)
}

type LiveResultAchievement struct {
	Position            int  `json:"position"`
	IsAlreadyAchieved   bool `json:"is_already_achieved"`
	IsCurrentlyAchieved bool `json:"is_currently_achieved"`
}

func (obj *LiveResultAchievement) ID() int64 {
	return int64(obj.Position)
}

type LiveFinishReq struct {
	LiveID           int64 `json:"live_id"`
	LiveFinishStatus int   `json:"live_finish_status"`
	LiveScore        struct {
		StartInfo                  any                                           `json:"start_info"`
		FinishInfo                 any                                           `json:"finish_info"`
		ResultDict                 []any                                         `json:"result_dict"`
		WaveStatDict               []any                                         `json:"wave_stat_dict"`
		TurnStatDict               []any                                         `json:"turn_stat_dict"`
		CardStatDict               generic.ObjectByObjectIDRead[*LiveFinishCard] `json:"card_stat_dict"`
		TargetScore                int                                           `json:"target_score"`
		CurrentScore               int                                           `json:"current_score"`
		ComboCount                 int                                           `json:"combo_count"`
		ChangeSquadCount           int                                           `json:"change_squad_count"`
		HighestComboCount          int                                           `json:"highest_combo_count"`
		RemainingStamina           int                                           `json:"remaining_stamina"`
		IsPerfectLive              bool                                          `json:"is_perfect_live"`
		IsPerfectFullCombo         bool                                          `json:"is_perfect_full_combo"`
		UseVoltageActiveSkillCount int                                           `json:"use_voltage_active_skill_count"`
		UseHealActiveSkillCount    int                                           `json:"use_heal_active_skill_count"`
		UseDebufActiveSkillCount   int                                           `json:"use_debuf_active_skill_count"`
		UseBufActiveSkillCount     int                                           `json:"use_buf_active_skill_count"`
		UseSpSkillCount            int                                           `json:"use_sp_skill_count"`
		CompleteAppealChanceCount  int                                           `json:"complete_appeal_chance_count"`
		TriggerCriticalCount       int                                           `json:"triggered_critical_count"`
		LivePower                  int                                           `json:"live_power"`
		SpSkillScoreList           []int                                         `json:"sp_skill_score_list"`
	} `json:"live_score"`
	ResumeFinishInfo any `json:"resume_finish_info"`
	RoomID           int `json:"room_id"`
}

type LiveFinishLiveResult struct {
	LiveDifficultyMasterID int                                              `json:"live_difficulty_master_id"`
	LiveDeckID             int                                              `json:"live_deck_id"`
	StandardDrops          []any                                            `json:"standard_drops"`
	AdditionalDrops        []any                                            `json:"additional_drops"`
	GimmickDrops           []any                                            `json:"gimmick_drops"`
	MemberLoveStatuses     generic.ObjectByObjectIDWrite[*MemberLoveStatus] `json:"member_love_statuses"`
	MVP                    struct {
		CardMasterID        int `json:"card_master_id"`
		GetVoltage          int `json:"get_voltage"`
		SkillTriggeredCount int `json:"skill_triggered_count"`
		AppealCount         int `json:"appeal_count"`
	} `json:"mvp"`
	Partner                     *model.UserBasicInfo                                  `json:"partner"`
	LiveResultAchievements      generic.ObjectByObjectIDWrite[*LiveResultAchievement] `json:"live_result_achievements"`
	LiveResultAchievementStatus struct {
		ClearCount       int `json:"clear_count"`
		GotVoltage       int `json:"got_voltage"`
		RemainingStamina int `json:"remaining_stamina"`
	} `json:"live_result_achievement_status"`
	Voltage                       int  `json:"voltage"`
	LastBestVoltage               int  `json:"last_best_voltage"`
	BeforeUserExp                 int  `json:"before_user_exp"`
	GainUserExp                   int  `json:"gain_user_exp"`
	IsRewardAccessoryInPresentBox bool `json:"is_reward_accessory_in_present_box"`
	ActiveEventResult             *any `json:"active_event_result"`
	LiveResultTower               *any `json:"live_result_tower"`
	LiveResultMemberGuild         *any `json:"live_result_member_guild"`
	LiveFinishStatus              int  `json:"live_finish_status"`
}

type LiveFinishLiveDifficultyInfo struct {
	LiveID                int `xorm:"'live_id'"`
	RewardUserExp         int
	ConsumedLP            int `xorm:"'consumed_lp'"`
	RewardBaseLovePoint   int
	RewardCenterLovePoint int `xorm:"-"`
	LoseAtDeath           bool
	IsCountTarget         bool
}

type LiveDifficultyMission struct {
	Position    int
	TargetValue int
	Reward      model.RewardByContent `xorm:"extends"`
}

func LiveFinish(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	req := LiveFinishReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	utils.CheckErr(err)
	exists, liveState := serverdb.LoadLiveState(UserID)
	utils.MustExist(exists)

	session := serverdb.GetSession(ctx, UserID)
	liveState.DeckID = session.UserStatus.LatestLiveDeckID
	liveState.LiveStage.LiveDifficultyID = session.UserStatus.LastLiveDifficultyID

	db := ctx.MustGet("masterdata.db").(*xorm.Engine)
	info := LiveFinishLiveDifficultyInfo{}
	exists, err = db.Table("m_live_difficulty").Where("live_difficulty_id = ?", liveState.LiveStage.LiveDifficultyID).Get(&info)
	utils.CheckErrMustExist(err, exists)

	liveMemberMappingID := 0
	db.Table("m_live").Where("live_id = ?", info.LiveID).Cols("live_member_mapping_id").Get(&liveMemberMappingID)
	centerPositions := []int{}
	err = db.Table("m_live_member_mapping").Where("mapping_id = ? AND is_center = 1", liveMemberMappingID).
		Cols("position").Find(&centerPositions)
	utils.CheckErr(err)
	info.RewardCenterLovePoint = klab.CenterBondGainBasedOnBondGain(info.RewardBaseLovePoint) / len(centerPositions)

	// record this live
	liveRecord := session.GetLiveDifficultyRecord(session.UserStatus.LastLiveDifficultyID)
	liveRecord.IsNew = false
	lastPlayDeck := session.BuildLastPlayLiveDifficultyDeck(liveState.DeckID, liveState.LiveStage.LiveDifficultyID)
	lastPlayDeck.Voltage = req.LiveScore.CurrentScore

	liveResult := LiveFinishLiveResult{
		LiveDifficultyMasterID: session.UserStatus.LastLiveDifficultyID,
		LiveDeckID:             session.UserStatus.LatestLiveDeckID,
		StandardDrops:          []any{},
		AdditionalDrops:        []any{},
		GimmickDrops:           []any{},
		Voltage:                lastPlayDeck.Voltage,
		LastBestVoltage:        liveRecord.MaxScore,
		BeforeUserExp:          session.UserStatus.Exp,
		LiveFinishStatus:       req.LiveFinishStatus}

	liveRecord.PlayCount++
	lastPlayDeck.IsCleared = req.LiveFinishStatus == enum.LiveFinishStatusCleared
	for i := 1; i <= 3; i++ {
		(*liveResult.LiveResultAchievements.AppendNew()).Position = i
	}
	(*liveResult.LiveResultAchievements.Objects[0]).IsAlreadyAchieved = liveRecord.ClearedDifficultyAchievement1 != nil
	(*liveResult.LiveResultAchievements.Objects[1]).IsAlreadyAchieved = liveRecord.ClearedDifficultyAchievement2 != nil
	(*liveResult.LiveResultAchievements.Objects[2]).IsAlreadyAchieved = liveRecord.ClearedDifficultyAchievement3 != nil
	if lastPlayDeck.IsCleared {
		// update clear record
		liveRecord.ClearCount++
		if liveRecord.MaxScore < req.LiveScore.CurrentScore { // if new high score
			liveRecord.MaxScore = req.LiveScore.CurrentScore
		}
		if liveRecord.MaxCombo < req.LiveScore.HighestComboCount {
			liveRecord.MaxCombo = req.LiveScore.HighestComboCount
		}
		if info.IsCountTarget { // counted toward target and profiles
			liveStats := session.GetUserLiveStats()
			idx := klab.LiveDifficultyTypeIndexFromLiveDifficultyID(liveState.LiveStage.LiveDifficultyID)
			liveStats.LivePlayCount[idx]++
			if liveRecord.ClearCount == 1 { // 1st clear
				liveStats.LiveClearCount[idx]++
			}
			session.UpdateUserLiveStats(liveStats)
		}

		// and award items
		missions := []LiveDifficultyMission{}

		db.Table("m_live_difficulty_mission").Where("live_difficulty_master_id = ?", session.UserStatus.LastLiveDifficultyID).
			OrderBy("position").Find(&missions)
		for i := 0; i < 3; i++ {
			if (i == 0) || (req.LiveScore.CurrentScore >= missions[i].TargetValue) {
				(*liveResult.LiveResultAchievements.Objects[i]).IsCurrentlyAchieved = true
				if !(*liveResult.LiveResultAchievements.Objects[i]).IsAlreadyAchieved { // new, add reward
					session.AddRewardContent(missions[i].Reward)
					switch i {
					case 0:
						liveRecord.ClearedDifficultyAchievement1 = new(int)
						*liveRecord.ClearedDifficultyAchievement1 = 1
					case 1:
						liveRecord.ClearedDifficultyAchievement2 = new(int)
						*liveRecord.ClearedDifficultyAchievement2 = 2
					case 2:
						liveRecord.ClearedDifficultyAchievement3 = new(int)
						*liveRecord.ClearedDifficultyAchievement3 = 3
					}
				}
			}
		}
		liveResult.GainUserExp = info.RewardUserExp
	}

	bondCardPosition := make(map[int]int)
	for i, _ := range req.LiveScore.CardStatDict.Objects {
		liveFinishCard := req.LiveScore.CardStatDict.Objects[i]

		// calculate mvp
		if liveFinishCard.GotVoltage > liveResult.MVP.GetVoltage {
			liveResult.MVP.GetVoltage = liveFinishCard.GotVoltage
			liveResult.MVP.CardMasterID = liveFinishCard.CardMasterID
			liveResult.MVP.SkillTriggeredCount = liveFinishCard.SkillTriggeredCount
			liveResult.MVP.AppealCount = liveFinishCard.AppealCount
		}

		// update card stat and member bond if cleared
		if lastPlayDeck.IsCleared {
			isCenter := (i+1 == centerPositions[0])
			isCenter = isCenter || ((len(centerPositions) > 1) && (i+1 == centerPositions[1]))
			addedBond := info.RewardBaseLovePoint
			if isCenter {
				addedBond += info.RewardCenterLovePoint
			}

			userCard := session.GetUserCard(liveFinishCard.CardMasterID)
			userCard.LiveJoinCount++
			userCard.ActiveSkillPlayCount += liveFinishCard.SkillTriggeredCount
			session.UpdateUserCard(userCard)
			// update member love point
			memberMasterID := klab.MemberMasterIDFromCardMasterID(liveFinishCard.CardMasterID)

			pos, exists := bondCardPosition[memberMasterID]
			// only use 1 card master id or an idol might be shown multiple times
			if !exists {
				memberLoveStatus := liveResult.MemberLoveStatuses.AppendNew()
				memberLoveStatus.RewardLovePoint = addedBond
				memberLoveStatus.CardMasterID = liveFinishCard.CardMasterID
				bondCardPosition[memberMasterID] = liveResult.MemberLoveStatuses.Length - 1
			} else {
				(*liveResult.MemberLoveStatuses.Objects[pos]).RewardLovePoint += addedBond
			}
		}
	}
	for memberMasterID, pos := range bondCardPosition {
		addedBond := session.AddLovePoint(memberMasterID,
			(*liveResult.MemberLoveStatuses.Objects[pos]).RewardLovePoint)
		(*liveResult.MemberLoveStatuses.Objects[pos]).RewardLovePoint = addedBond
	}

	liveResult.LiveResultAchievementStatus.ClearCount = liveRecord.ClearCount
	liveResult.LiveResultAchievementStatus.GotVoltage = req.LiveScore.CurrentScore
	liveResult.LiveResultAchievementStatus.RemainingStamina = req.LiveScore.RemainingStamina
	if liveState.PartnerUserID != 0 {
		liveResult.Partner = new(model.UserBasicInfo)
		*liveResult.Partner = session.GetOtherUserBasicProfile(liveState.PartnerUserID)
	}
	session.UpdateLiveDifficultyRecord(liveRecord)
	session.SetLastPlayLiveDifficultyDeck(lastPlayDeck)
	liveFinishResp := session.Finalize(handler.GetUserData("userModelDiff.json"), "user_model_diff")
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result", liveResult)

	resp := handler.SignResp(ctx.GetString("ep"), liveFinishResp, config.SessionKey)
	// fmt.Println(resp)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
