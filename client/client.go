package client

import (
	"fmt"
	"os"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
	c "github.com/michelemendel/binance/constant"
	"github.com/michelemendel/binance/util"
)

type Client struct {
	Env       string
	Conn      *binance_connector.Client
	APIKey    string
	SecretKey string
	Timeout   time.Duration
	BaseAPI   string
	BaseWS    string
}

func NewClient(env string, conn *binance_connector.Client, apiKey, secretKey, baseAPI, baseWS string) *Client {
	return &Client{
		Env:       env,
		Conn:      conn,
		APIKey:    apiKey,
		SecretKey: secretKey,
		Timeout:   c.TIMEOUT,
		BaseAPI:   baseAPI,
		BaseWS:    baseWS,
	}
}
func Run() {
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

	conn := binance_connector.NewClient(apiKey, secretKey, baseAPI)
	client := NewClient(env, conn, apiKey, secretKey, baseAPI, baseWS)
	client.Ping()
	fmt.Printf("env:%s\nbaseAPI:%s\nbaseWS:%s\n", client.Env, client.BaseAPI, client.BaseWS)

	// Buy/Sell
	// qty := buy(client)
	// fmt.Println("qty:", qty)
	// sell(client, qty)

	// client.ExchangeInfo("BTCFDUSD")
	// client.AccountStatus()

	// Streams
	qty := 0.00136
	quoteOrderQuantity := 100.0
	// symbols := []string{"BTCFDUSD", "ETHFDUSD"}
	symbols := []string{"BTCFDUSD"}
	client.StreamMiniTicker(symbols, quoteOrderQuantity, qty)

}

func buy(client *Client) float64 {
	symbol := "BTCFDUSD"
	quoteOrderQuantity := 100.0
	order, err := client.Buy(symbol, quoteOrderQuantity, 0)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	qty := util.String2Float(order.ExecutedQty)
	price := util.String2Float(order.Fills[0].Price)
	fmt.Printf("bought:%s, price:%v for %v, received amount:%v\n", symbol, price, quoteOrderQuantity, qty)
	return qty
}

func sell(client *Client, qty float64) {
	symbol := "BTCFDUSD"
	order, err := client.Sell(symbol, 0, qty)
	if err != nil {
		fmt.Println(err)
		return
	}
	exQty := util.String2Float(order.ExecutedQty)
	price := util.String2Float(order.Fills[0].Price)
	total := exQty * price
	commission := util.String2Float(order.Fills[0].Commission)
	fmt.Printf("sold:%s, price:%v for qty %v received amount %v, total %v (commission:%v)\n", symbol, price, qty, exQty, total, commission)
}
