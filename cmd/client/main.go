package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
	"github.com/michelemendel/binance/client"
	c "github.com/michelemendel/binance/constant"
)

// https://github.com/binance/binance-connector-go

func init() {
	envFile := filepath.Join("", ".env")
	err := godotenv.Load(envFile)
	if err != nil {
		slog.Error("error loading file ", "file", envFile, "error", err)
	}
}

func main() {
	var baseAPI string
	var baseWS string
	var apiKey string
	var secretKey string

	env := os.Getenv("ENV")
	if env == "test" {
		baseAPI = c.BASE_API_TEST
		baseWS = c.BASE_WS_TEST
		apiKey = os.Getenv("API_KEY_TEST")
		secretKey = os.Getenv("SECRET_KEY_TEST")
	} else {
		baseAPI = c.BASE_API_PROD_0
		baseWS = c.BASE_WS_PROD_1
		apiKey = os.Getenv("API_KEY")
		secretKey = os.Getenv("SECRET_KEY")
	}

	fmt.Println("baseAPI:", baseAPI)
	conn := binance_connector.NewClient(apiKey, secretKey, baseAPI)
	client := client.NewClient(env, conn, apiKey, secretKey, baseAPI, baseWS)

	fmt.Printf("env:%s\nbaseAPI:%s\nbaseWS:%s\n", client.Env, client.BaseAPI, client.BaseWS)

	// Buy
	// quoteAssetAmount := 100.0
	symbol := "BTCFDUSD"
	// client.Buy(symbol, quoteAssetAmount, 0, "MARKET")
	// order, err := client.Buy(symbol, quoteAssetAmount, 0)
	// if err != nil {
	// fmt.Println(err)
	// return
	// }
	// qty := util.String2Float(order.ExecutedQty)
	// fmt.Printf("You bought %s for %v, and got amount:%v\n", symbol, quoteAssetAmount, qty)

	// Sell
	qty := 0.001
	total, commission := client.Sell(symbol, 0, qty)
	fmt.Printf("You sold %s, amount:%v, commission:%v\n", symbol, total, commission)

	// client.ExchangeInfo("BTCFDUSD")
	// client.AccountStatus()
	// client.Time()
}
