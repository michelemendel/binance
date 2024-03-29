package entity

// --------------------------------------------------------------------------------
type PingResp struct{}

// --------------------------------------------------------------------------------
type TimeResp struct {
	ServerTime uint64 `json:"serverTime"`
}

// --------------------------------------------------------------------------------
// System
type ExchangeInfoResp struct {
	Timezone   string `json:"timezone"`
	ServerTime uint64 `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
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
		AllowTrailingStop          bool     `json:"allowTrailingStop"`
		CancelReplaceAllowed       bool     `json:"cancelReplaceAllowed"`
		IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
		IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
		Filters                    []struct {
			FilterType            string `json:"filterType"`
			MinPrice              string `json:"minPrice,omitempty"`
			MaxPrice              string `json:"maxPrice,omitempty"`
			TickSize              string `json:"tickSize,omitempty"`
			MinQty                string `json:"minQty,omitempty"`
			MaxQty                string `json:"maxQty,omitempty"`
			StepSize              string `json:"stepSize,omitempty"`
			Limit                 int    `json:"limit,omitempty"`
			MinTrailingAboveDelta int    `json:"minTrailingAboveDelta,omitempty"`
			MaxTrailingAboveDelta int    `json:"maxTrailingAboveDelta,omitempty"`
			MinTrailingBelowDelta int    `json:"minTrailingBelowDelta,omitempty"`
			MaxTrailingBelowDelta int    `json:"maxTrailingBelowDelta,omitempty"`
			BidMultiplierUp       string `json:"bidMultiplierUp,omitempty"`
			BidMultiplierDown     string `json:"bidMultiplierDown,omitempty"`
			AskMultiplierUp       string `json:"askMultiplierUp,omitempty"`
			AskMultiplierDown     string `json:"askMultiplierDown,omitempty"`
			AvgPriceMins          int    `json:"avgPriceMins,omitempty"`
			MinNotional           string `json:"minNotional,omitempty"`
			ApplyMinToMarket      bool   `json:"applyMinToMarket,omitempty"`
			MaxNotional           string `json:"maxNotional,omitempty"`
			ApplyMaxToMarket      bool   `json:"applyMaxToMarket,omitempty"`
			MaxNumOrders          int    `json:"maxNumOrders,omitempty"`
			MaxNumAlgoOrders      int    `json:"maxNumAlgoOrders,omitempty"`
		} `json:"filters"`
		Permissions                     []string `json:"permissions"`
		DefaultSelfTradePreventionMode  string   `json:"defaultSelfTradePreventionMode"`
		AllowedSelfTradePreventionModes []string `json:"allowedSelfTradePreventionModes"`
	} `json:"symbols"`
}

type ExchangeInfoRespX struct {
	Timezone      string `json:"timezone"`
	ServerTime    uint64 `json:"serverTime"`
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
type AccountStatusResp struct {
	Data string `json:"data"`
}

// --------------------------------------------------------------------------------
type WalletStatusResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
