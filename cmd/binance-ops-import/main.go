package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/billettc/trades-tracker/db"
	"github.com/billettc/trades-tracker/models"
	"github.com/billettc/trades-tracker/price"
)

func main() {
	//f, err := os.Open("/Users/cbillett/devel/perso/go/trades-tracker/cmd/binance-ops-import/binance-operations-2018.csv")
	f, err := os.Open("/Users/cbillett/Downloads/binance-2021-operations.csv")
	check(err)

	database, err := db.NewMongoDB("mongodb://localhost:27017")
	check(err)

	priceGetter := price.NewPriceGetter(price.API_KEY, price.SECRET_KEY, database)

	accountState := &AccountState{
		priceGetter: priceGetter,
		database:    database,
		ctx:         context.Background(),
	}

	r := csv.NewReader(f)
	_, err = r.Read() //skipping header
	check(err)

	var operations []*models.Operation

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		o, err := models.NewOperationFromRecord(record)
		check(err)

		if o.Type == "Distribution" {
			continue
		}
		o, err = accountState.HandleOperation(o)
		check(err)
		operations = append(operations, o)
	}

	for _, o := range operations {
		err = database.SaveOperation(o)
		check(err)
	}

	//out, err := os.Create("/Users/cbillett/t/ops2.csv")
	//check(err)
	//
	//w := csv.NewWriter(out)
	//err = w.Write([]string{"UTC_Time", "Account", "Operation", "Coin", "USD_Price", "Change", "Change_USD"})
	//check(err)
	//
	//for _, o := range operations {
	//	switch o.Type {
	//	case "Sell", "Buy", "Fee":
	//	case "Deposit":
	//	case "Withdraw":
	//		continue
	//	default:
	//		panic(fmt.Sprintf("wtf: ops type: %s", o.Type))
	//	}
	//	err := w.Write(o.ToRecord())
	//	w.Flush()
	//	check(err)
	//}
	//
	//out.Close()
	fmt.Println("All done!")
}

type AccountState struct {
	priceGetter *price.Getter
	database    *db.MongoDB
	ctx         context.Context
	balance     float64
	gain        float64
}

func (s *AccountState) HandleOperation(o *models.Operation) (*models.Operation, error) {
	o.ChangeUSD = o.Change
	if o.Coin != "USDT" {

		suffix := ""
		switch o.Coin {
		case "SOL", "TWT":
			suffix = "BUSD"
		default:
			suffix = "USDT"

		}

		price, err := s.priceGetter.PriceAtDate(s.ctx, o.Coin+suffix, o.DateTime)
		if err != nil {
			return nil, fmt.Errorf("getting usd price for operation : %w", err)
		}
		o.Price = price.Price
		o.ChangeUSD = o.Change * price.Price
	}

	//prevBalance := s.balance
	//s.balance += o.ChangeUSD
	//var gain float64
	//switch o.Type {
	//case "Sell", "Buy", "Fee":
	//	gain = s.balance - prevBalance
	//	s.gain += s.balance - prevBalance
	//}
	//
	//fmt.Println("-----")
	//fmt.Println(o.Type, "Coin:", o.Coin, "Change:", o.Change, "USD Change:", o.ChangeUSD)
	//fmt.Println("Gain:", s.gain, "Ops Gain:", gain, "Balance:", s.balance, "Prev", prevBalance)

	return o, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
