package userdata

import (
	"elichika/model"
	"elichika/utils"
)

func (session *Session) GetTradeProductUser(productID int) int {
	result := 0
	exists, err := session.Db.Table("u_trade_product").
		Where("user_id = ? AND product_id = ?", session.UserStatus.UserID, productID).
		Cols("traded_count").Get(&result)
	utils.CheckErr(err)
	if !exists {
		result = 0
	}
	return result
}

func (session *Session) SetTradeProductUser(productID, newTradedCount int) {
	record := model.TradeProductUser{
		UserID:      session.UserStatus.UserID,
		ProductID:   productID,
		TradedCount: newTradedCount,
	}

	exists, err := session.Db.Table("u_trade_product").
		Where("user_id = ? AND product_id = ?", session.UserStatus.UserID, productID).
		Update(record)
	utils.CheckErr(err)
	if exists == 0 {
		_, err := session.Db.Table("u_trade_product").Insert(record)
		utils.CheckErr(err)
	}
}

func (session *Session) GetTrades(tradeType int) []model.Trade {
	trades := []model.Trade{}
	for _, trade := range session.Gamedata.Trade.TradesByType[tradeType] {
		for j, product := range trade.Products {
			product.TradedCount = session.GetTradeProductUser(product.ProductID)
			trade.Products[j] = product
		}
		trades = append(trades, trade)
	}
	return trades
}

// return whether the item is added to present box
func (session *Session) ExecuteTrade(productID int, tradeCount int) bool {
	// update count
	tradedCount := session.GetTradeProductUser(productID)
	tradedCount += tradeCount
	session.SetTradeProductUser(productID, tradedCount)

	// award items and take away source item
	product := session.Gamedata.Trade.Products[productID]
	trade := session.Gamedata.Trade.Trades[product.TradeID]
	content := product.ActualContent
	content.ContentAmount *= int64(tradeCount)
	session.AddResource(content)
	session.RemoveResource(model.Content{
		ContentType:   trade.SourceContentType,
		ContentID:     trade.SourceContentID,
		ContentAmount: int64(product.SourceAmount) * int64(tradeCount),
	})

	return true
}
