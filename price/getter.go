package price

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2"
)

type Getter struct {
	client *binance.Client
}

func NewPriceGetter(apiKey string, secretKey string) *Getter {
	client := binance.NewClient(apiKey, secretKey)

	return &Getter{
		client: client,
	}
}

func (g Getter) PriceAtDate(ctx context.Context, symbol string, time *time.Time) float64 {
	g.client.NewKlinesService().
		Symbol(symbol).
		StartTime(time.Unix()).
		EndTime(time.Unix()).
		Do(ctx)

	return 1
}
