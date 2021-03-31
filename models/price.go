package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Price struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Symbol   string
	DateTime MinuteRoundedTime
	Price    float64
}

func NewPrice(symbol string, dateTime MinuteRoundedTime, price float64) *Price {

	return &Price{
		Symbol:   symbol,
		DateTime: dateTime,
		Price:    price,
	}
}

func (p Price) String() string {
	return fmt.Sprintf("symbol: %s price: %f: date_time: %s", p.Symbol, p.Price, p.DateTime)
}
