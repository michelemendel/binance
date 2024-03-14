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

//--------------------------------------------------------------------------------
// Streams

// Individual Symbol Mini Ticker Stream
// https://binance-docs.github.io/apidocs/spot/en/#individual-symbol-mini-ticker-stream
func (client Client) StreamMiniTicker(symbols []string, quoteOrderQuantity, qty float64) {
	fmt.Println("StreamMiniTicker", client.BaseWS, symbols)
	isCombined := true
	wsStreamClient := binance_connector.NewWebsocketStreamClient(isCombined, client.BaseWS)

	wsHandler := func(e *binance_connector.WsMarketTickerStatEvent) {
		// fmt.Println(binance_connector.PrettyPrint(e))
		value := util.String2Float(e.LastPrice) * qty
		profit := value - quoteOrderQuantity
		// profitInPercents := (profit / quoteOrderQuantity) * 100
		fmt.Printf("%s : %s : %v : %v : %v\n", e.Symbol, e.LastPrice, quoteOrderQuantity, value, profit)
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	doneCh, stopCh, err := wsStreamClient.WsCombinedMarketTickersStatServe(symbols, wsHandler, errHandler)
	if err != nil {
		fmt.Println("StreamMiniTicker", err)
		return
	}

	go func() {
		time.Sleep(20 * time.Second)
		fmt.Println("stopping stream...")
		close(stopCh)
		close(doneCh)
	}()

	done := <-doneCh
	fmt.Println("done", done)
}
