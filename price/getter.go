package price

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/billettc/trades-tracker/db"
	"github.com/billettc/trades-tracker/models"
)

type Getter struct {
	client   *binance.Client
	database *db.MongoDB
}

func NewPriceGetter(apiKey string, secretKey string, database *db.MongoDB) *Getter {
	client := binance.NewClient(apiKey, secretKey)
	return &Getter{
		client:   client,
		database: database,
	}
}

func (g Getter) PriceAtDate(ctx context.Context, symbol string, dateTime models.MinuteRoundedTime) (*models.Price, error) {
	price, err := g.database.GetPrice(symbol, dateTime)
	if err != nil {
		return nil, fmt.Errorf("get price from database: %w", err)
	}

	if price != nil {
		return price, nil
	}

	start := dateTime
	end := dateTime
	//start := models.NewMinuteRoundedTime(dateTime.Add(-1 * time.Hour))
	//end := models.NewMinuteRoundedTime(dateTime.Add(1 * time.Hour))
	fmt.Println("Calling binance: symbol:", symbol, "time", dateTime, "micro:", start.MicroSecond())

	//Tuesday, August 11, 2020 6:00:00 AM

	result, err := g.client.NewKlinesService().
		Symbol(symbol).
		StartTime(start.MicroSecond()).
		EndTime(end.MicroSecond()).
		Interval("1m").
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("kline service do: %w", err)
	}

	kline := result[0]
	h, err := strconv.ParseFloat(kline.High, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing high price: %w", err)
	}
	l, err := strconv.ParseFloat(kline.Low, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing high price: %w", err)
	}

	p := (h + l) / 2.0
	price = models.NewPrice(symbol, dateTime, p)

	err = g.database.SavePrice(price)
	if err != nil {
		return nil, fmt.Errorf("saving price to db: %w", err)
	}

	return price, nil
}
