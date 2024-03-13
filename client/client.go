package client

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
	c "github.com/michelemendel/binance/constant"
	"github.com/michelemendel/binance/entity"
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

//--------------------------------------------------------------------------------
// Order

// https://binance-docs.github.io/apidocs/spot/en/#new-order-trade
// symbol-BTCFDUSD, type-MARKET, quantity-0.001, orderType-Market
func (c Client) Buy(pair string, quoteOrderQuantity, quantity float64) (*binance_connector.CreateOrderResponseFULL, error) {
	order, err := c.Order("BUY", pair, quoteOrderQuantity, quantity)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (c Client) Sell(pair string, quoteOrderQuantity, quantity float64) (*binance_connector.CreateOrderResponseFULL, error) {
	order, err := c.Order("SELL", pair, quoteOrderQuantity, quantity)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (c Client) Order(side, pair string, quoteOrderQuantity, quantity float64) (*binance_connector.CreateOrderResponseFULL, error) {
	fmt.Printf("side:%s, pair:%s, quoteOrderQuantity:%v, quantity:%v\n", side, pair, quoteOrderQuantity, quantity)

	orderType := "MARKET"
	newOrder := c.Conn.
		NewCreateOrderService().
		Symbol(pair).
		Side(side).
		Type(orderType)

	if quoteOrderQuantity > 0 {
		newOrder = newOrder.QuoteOrderQty(quoteOrderQuantity)
	} else if quantity > 0 {
		newOrder = newOrder.Quantity(quantity)
	} else {
		return nil, fmt.Errorf("quoteOrderQuantity or quantity must be greater than 0")
	}

	resp, err := newOrder.Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error creating order: %v", err)
	}

	order := resp.(*binance_connector.CreateOrderResponseFULL)
	// fmt.Println(binance_connector.PrettyPrint(resp))
	util.PP(order)
	return order, nil
}

// --------------------------------------------------------------------------------
// User
func (client Client) AccountStatus() {
	resp := client.Get(c.PATH_GET_ACCOUNT_STATUS, "")
	var decData entity.AccountStatusResp
	decode(resp, &decData)

	util.PP(decData)
}

// Trade Fee (USER_DATA)
// GET /sapi/v1/asset/tradeFee

// Query User Wallet Balance (USER_DATA)
// GET /sapi/v1/asset/wallet/balance

// --------------------------------------------------------------------------------
// System

// Test Connectivity
// https://binance-docs.github.io/apidocs/spot/en/#test-connectivity
func (c Client) Ping() {
	err := c.Conn.NewPingService().Do(context.Background())
	if err != nil {
		code := regexp.MustCompile(`code=(\d+)`).FindStringSubmatch(err.Error())[1]
		if code == "0" {
			fmt.Printf("%s, no connection, quitting\n", c.BaseAPI)
			os.Exit(0)
		}
	}
	fmt.Printf("%s, connection OK\n", c.BaseAPI)
}

// Check Server Time
// https://binance-docs.github.io/apidocs/spot/en/#check-server-time
func (client Client) Time() {
	resp := client.Get(c.PATH_TIME, "")
	fmt.Println("resp:", string(resp))

	var decData entity.TimeResp
	decode(resp, &decData)

	util.PP(decData)
}

// Symbol Price Ticker
// https://binance-docs.github.io/apidocs/spot/en/#symbol-price-ticker
// GET /api/v3/ticker/price
func (c Client) SymbolPriceTicker(pair string) {
	priceTicker, err := c.Conn.
		NewTickerPriceService().
		Symbol(pair).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(binance_connector.PrettyPrint(priceTicker))
}

func (client Client) ExchangeInfo(pair string) {
	query := "symbol=" + pair
	resp := client.Get(c.PATH_EXCHANGE_INFO, query)
	var decData entity.ExchangeInfoRespX
	decode(resp, &decData)
	decData.ServerTimeStr = util.Time2String(decData.ServerTime)

	util.PP(decData)
}
