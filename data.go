package main

// --------------------------------------------------------------------------------
const PING_PATH = "/api/v3/ping"

type PingResp struct{}

// --------------------------------------------------------------------------------
const TIME_PATH = "/api/v3/time"

type TimeResp struct {
	ServerTime uint64 `json:"serverTime"`
}

// --------------------------------------------------------------------------------
const EXCHANGE_INFO_PATH = "/api/v3/exchangeInfo"

type ExchangeInfoResp struct {
	Timezone      string `json:"timezone"`
	ServerTime    int64  `json:"serverTime"`
	ServerTimeStr string `json:"serverTimeStr"` //Not part of API
	RateLimits    []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol                     string   `json:"symbol"`
		Status                     string   `json:"status"`
		BaseAsset                  string   `json:"baseAsset"`
		BaseAssetPrecision         int      `json:"baseAssetPrecision"`
		QuoteAsset                 string   `json:"quoteAsset"`
		QuotePrecision             int      `json:"quotePrecision"`
		QuoteAssetPrecision        int      `json:"quoteAssetPrecision"`
		BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
		QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
		OrderTypes                 []string `json:"orderTypes"`
		IcebergAllowed             bool     `json:"icebergAllowed"`
		OcoAllowed                 bool     `json:"ocoAllowed"`
		QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
		IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
		IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
		Filters                    []struct {
			FilterType       string `json:"filterType"`
			MinPrice         string `json:"minPrice,omitempty"`
			MaxPrice         string `json:"maxPrice,omitempty"`
			TickSize         string `json:"tickSize,omitempty"`
			MinQty           string `json:"minQty,omitempty"`
			MaxQty           string `json:"maxQty,omitempty"`
			StepSize         string `json:"stepSize,omitempty"`
			MinNotional      string `json:"minNotional,omitempty"`
			Limit            int    `json:"limit,omitempty"`
			MaxNumAlgoOrders int    `json:"maxNumAlgoOrders,omitempty"`
		} `json:"filters"`
		Permissions []string `json:"permissions"`
	} `json:"symbols"`
}

// --------------------------------------------------------------------------------
const ACCOUNT_STATUS_PATH = "/sapi/v1/account/status"

type AccountStatusQuery struct {
	RecvWindow int    `json:"recvWindow"`
	Timestamp  int64  `json:"timestamp"`
	Signature  string `json:"signature"`
}

type AccountStatusParams struct {
	Query AccountStatusQuery `json:"query"`
}

type AccountStatusResp struct {
	Data string `json:"data"`
}

// --------------------------------------------------------------------------------
// System status

const WALLET_STATUS_PATH = "/sapi/v1/system/status"

type WalletStatusResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
