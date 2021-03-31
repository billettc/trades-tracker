package price

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/billettc/trades-tracker/models"

	"github.com/billettc/trades-tracker/db"

	"github.com/stretchr/testify/require"
)

func TestGetter_PriceAtDate(t *testing.T) {
	database, err := db.NewMongoDB("mongodb://localhost:27017")
	require.NoError(t, err)

	g := NewPriceGetter(API_KEY, SECRET_KEY, database)

	price, err := g.PriceAtDate(context.Background(), "EOSUSDT", models.NewMinuteRoundedTime(time.Date(2019, 2, 7, 15, 04, 0, 0, time.UTC)))
	require.NoError(t, err)
	fmt.Println("price:", price)
	price, err = g.PriceAtDate(context.Background(), "EOSUSDT", models.NewMinuteRoundedTime(time.Date(2018, 12, 2, 14, 34, 0, 0, time.UTC)))
	require.NoError(t, err)
	fmt.Println("price:", price)
	price, err = g.PriceAtDate(context.Background(), "EOSUSDT", models.NewMinuteRoundedTime(time.Date(2018, 12, 3, 14, 34, 0, 0, time.UTC)))
	require.NoError(t, err)
	fmt.Println("price:", price)
}
