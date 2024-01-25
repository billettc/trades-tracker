package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/billettc/trades-tracker/db"
)

func main() {

	out, err := os.Create("/Users/cbillett/t/all-ops.csv")
	check(err)

	database, err := db.NewMongoDB("mongodb://localhost:27017")
	check(err)

	w := csv.NewWriter(out)
	err = w.Write([]string{"UTC_Time", "Account", "Operation", "Coin", "USD_Price", "Change", "Change_USD"})
	check(err)

	operations, err := database.ListOperations()

	for _, o := range operations {
		switch o.Type {
		case "Sell", "Buy", "Fee":
		case "Withdraw":
			//continue
		case "Deposit":
			//continue
		case "Transaction Related":
		case "Large OTC trading":
		case "Funding Fee":
		case "Realize profit and loss":
		case "NFT transaction":
		case "Small assets exchange BNB":
			continue
		default:
			panic(fmt.Sprintf("wtf: ops type: %s", o.Type))
		}
		err := w.Write(o.ToRecord())
		w.Flush()
		check(err)
	}

	err = out.Close()
	check(err)

	fmt.Println("All done!")

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
