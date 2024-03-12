package constant

import "time"

// Base URLs for Binance prod server
const (
	// This base endpoint can be used to access the following API endpoints that have NONE as security type
	API_OPEN = "https://data.binance.com"

	//
	BASE_API_PROD_0 = "https://api.binance.com"
	BASE_API_PROD_1 = "https://api1.binance.com"
	BASE_API_PROD_2 = "https://api2.binance.com"
	BASE_API_PROD_3 = "https://api3.binance.com"
	BASE_API_PROD_4 = "https://api4.binance.com"

	// Websocket Market Streams
	// User Data Streams are accessed at /ws/<listenKey> or /stream?streams=<listenKey>
	BASE_WS_PROD_1 = "wss://stream.binance.com:9443" // /ws, /stream
	BASE_WS_PROD_2 = "wss://stream.binance.com:443"
	BASE_WS_PROD_3 = "wss://ws-api.binance.com" // /ws-api/v3
	// Streams can be accessed either in a single raw stream or in a combined stream
	// Raw streams are accessed at /ws/<streamName>
	// Combined streams are accessed at /stream?streams=<streamName1>/<streamName2>/<streamName3>
	// Combined stream events are wrapped as follows: {"stream":"<streamName>","data":<rawPayload>}
	// All symbols for streams are lowercase
	// A single connection to stream.binance.com is only valid for 24 hours; expect to be disconnected at the 24 hour mark
	// The websocket server will send a ping frame every 3 minutes. If the websocket server does not receive a pong frame back from the connection within a 10 minute period, the connection will be disconnected. Unsolicited pong frames are allowed.
	// The base endpoint wss://data-stream.binance.com can be subscribed to receive market data messages. Users data stream is NOT available from this URL.
)

// Base URLs for Binance test server
const (
	BASE_API_TEST = "https://testnet.binance.vision" // /api
	BASE_WS_TEST  = "wss://testnet.binance.vision"   // /ws-api/v3, /ws, /stream
)

// API paths
const (
	PATH_PING               = "/api/v3/ping"
	PATH_TIME               = "/api/v3/time"
	PATH_EXCHANGE_INFO      = "/api/v3/exchangeInfo"
	PATH_GET_ACCOUNT_STATUS = "/sapi/v1/account/status"
	PATH_WALLET_STATUS      = "/sapi/v1/system/status"
)

// TODO: Not sure I need these, since they are already set in the paths above.
// API path types
// const (
// 	PATH_API    = "/api"
// 	PATH_SAPI   = "/sapi"
// 	PATH_WS     = "/ws"
// 	PATH_WS_API = "/ws-api/v3"
// 	PATH_STREAM = "/stream"
// )

const (
	TIMEOUT_DURATION_MILLISECOND = 10000
	TIMEOUT                      = time.Duration(TIMEOUT_DURATION_MILLISECOND) * time.Millisecond
)
