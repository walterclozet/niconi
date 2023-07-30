package serverdb

import (
	"elichika/model"

	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

// fetch a card, use the value in diff is present, otherwise fetch from db
func (session *Session) GetCard(cardMasterID int) model.CardInfo {
	card, exist := session.CardDiffs[cardMasterID]
	if exist {
		return card
	}
	card = model.CardInfo{}
	exists, err := Engine.Table("s_user_card").
	Where("user_id = ? AND card_master_id = ?", session.UserInfo.UserID, cardMasterID).Get(&card)
	if err != nil {
		panic(err)
	}

	// if not in db then fetch from json
	if !exists {
		// insert this card, from json for now
		cardData := DbGetUserData("userCard.json")
		gjson.Parse(cardData).Get("user_card_by_card_id").ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				if err := json.Unmarshal([]byte(value.String()), &card); err != nil {
					panic(err)
				}
				if card.CardMasterID == cardMasterID {
					exists = true
					return false
				}
			}
			return true
		})
		if !exists {
			panic("Card doesn't exist")
		}
		card.UserID = session.UserInfo.UserID
		fmt.Println("Insert new card, user_id: ", card.UserID, ", card_master_id: ", cardMasterID)
		_, err := Engine.Table("s_user_card").Insert(&card)
		if err != nil {
			panic(err)
		}
	}
	return card
}

func (session *Session) GetAllCards() []model.CardInfo {
	var cards []model.CardInfo
	err := Engine.Table("s_user_card").Where("user_id = ?", session.UserInfo.UserID).Find(&cards)
	if err != nil {
		panic(err)
	}
	return cards
}

func (session *Session) UpdateCard(card model.CardInfo) {
	session.CardDiffs[card.CardMasterID] = card
}

func (session *Session) FinalizeCardDiffs() []any {
	userCardByCardID := []any{}
	for cardMasterID, card := range session.CardDiffs {
		userCardByCardID = append(userCardByCardID, cardMasterID)
		userCardByCardID = append(userCardByCardID, card)
		// .AllCols() is necessary to all field
		affected, err := Engine.Table("s_user_card").
			Where("user_id = ? AND card_master_id = ?", session.UserInfo.UserID, cardMasterID).AllCols().Update(card)
		if (err != nil) || (affected != 1) {
			panic(err)
		}
	}
	return userCardByCardID
}

// insert all the cards
func (session *Session) InsertCards(cards []model.CardInfo) {
	affected, err := Engine.Table("s_user_card").Insert(&cards)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", affected, " cards")
}
