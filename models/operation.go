package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Operation struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" text:"-"`
	DateTime  MinuteRoundedTime
	Account   string
	Type      string
	Coin      string
	Change    float64
	ChangeUSD float64 `bson:"change_usd"`
	Price     float64
}

func NewOperationFromRecord(record []string) (*Operation, error) {
	change, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing change: %w", err)
	}

	t, err := time.Parse("2006-01-02 15:04:05", record[0])
	if err != nil {
		return nil, fmt.Errorf("parsing date time: %w", err)
	}

	return &Operation{
		DateTime: NewMinuteRoundedTime(t),
		Account:  record[1],
		Type:     record[2],
		Coin:     record[3],
		Change:   change,
	}, nil
}

func (o *Operation) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}

	return string(data)
}

func (o *Operation) ToRecord() []string {
	//err = w.Write([]string{"UTC_Time", "Account", "Operation", "Coin", "USD_Price", "Change", "Change_USD"})
	return []string{o.DateTime.String(), o.Account, o.Type, o.Coin, fmt.Sprintf("%f", o.Price), fmt.Sprintf("%f", o.Change), fmt.Sprintf("%f", o.ChangeUSD)}
}
