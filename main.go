package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/michelemendel/binance/utils"
	"go.uber.org/zap"
)

const (
	// Production base URLs

	// The base endpoint can be used to access the following API endpoints that have NONE as security type
	API_OPEN = "https://data.binance.com"

	API_0 = "https://api.binance.com"
	API_1 = "https://api1.binance.com"
	API_2 = "https://api2.binance.com"
	API_3 = "https://api3.binance.com"
	API_4 = "https://api4.binance.com"

	// Websocket Market Streams
	// User Data Streams are accessed at /ws/<listenKey> or /stream?streams=<listenKey>
	MARKET_WS_1 = "wss://stream.binance.com:9443" // /ws, /stream
	MARKET_WS_2 = "wss://stream.binance.com:443"
	MARKET_WS_3 = "wss://ws-api.binance.com" // /ws-api/v3
	// Streams can be accessed either in a single raw stream or in a combined stream
	// Raw streams are accessed at /ws/<streamName>
	// Combined streams are accessed at /stream?streams=<streamName1>/<streamName2>/<streamName3>
	// Combined stream events are wrapped as follows: {"stream":"<streamName>","data":<rawPayload>}
	// All symbols for streams are lowercase
	// A single connection to stream.binance.com is only valid for 24 hours; expect to be disconnected at the 24 hour mark
	// The websocket server will send a ping frame every 3 minutes. If the websocket server does not receive a pong frame back from the connection within a 10 minute period, the connection will be disconnected. Unsolicited pong frames are allowed.
	// The base endpoint wss://data-stream.binance.com can be subscribed to receive market data messages. Users data stream is NOT available from this URL.

	// Test base URLs
	API_0_TEST     = "https://testnet.binance.vision" // /api
	MARKET_WS_TEST = "wss://testnet.binance.vision"   // /ws-api/v3, /ws, /stream

	// path type
	// API    = "/api"
	// SAPI   = "/sapi"
	// WS     = "/ws"
	// WS_API = "/ws-api/v3"
	// STREAM = "/stream"
)

const (
	TIMEOUT_DURATION_MILLISECOND = 10000
	TIMEOUT                      = time.Duration(TIMEOUT_DURATION_MILLISECOND) * time.Millisecond
)

type Server struct {
	Timeout time.Duration
	BaseAPI string
	BaseWS  string
}

func createServer(baseAPI, baseWS string) *Server {
	return &Server{
		Timeout: TIMEOUT,
		BaseAPI: baseAPI,
		BaseWS:  baseWS,
	}
}

var lg *zap.SugaredLogger

var SecretKey string

func init() {
	lg = utils.Log()
	envFile := filepath.Join("", ".env")
	err := godotenv.Load(envFile)
	if err != nil {
		lg.Panic("[main] Error loading file ", envFile)
	}
	SecretKey = os.Getenv("SECRET_KEY")
}

func main() {
	// fmt.Println("timeMillis", utils.ToTimeMillis())
	// fmt.Println("signature", utils.Signature())
	// testAccountStatus()
	// testExchangeInfo()
	// testTime()
}

func testAccountStatus() {
	accountStatusParams := AccountStatusParams{
		Query: AccountStatusQuery{
			RecvWindow: 5000,
			Timestamp:  utils.ToTimeMillis(),
			Signature:  "",
		},
	}

	fmt.Println("accountStatusParams", accountStatusParams)

	server := createServer(API_0, MARKET_WS_1)
	// server := createServer(API_0_TEST, MARKET_WS_TEST)

	resp := server.get(ACCOUNT_STATUS_PATH)
	var decData AccountStatusResp
	decode(resp, &decData)
	fmt.Printf("RESPONSE: %v, %T\n", decData.Data, decData)
}

func makePayload(params any) string {
	payload, err := json.Marshal(params)
	if err != nil {
		lg.Errorf("Error marshalling payload %s", err)
		return ""
	}
	return string(payload)
}

func testExchangeInfo(pair string) {
	// server := createServer(API_0, MARKET_WS_1)
	server := createServer(API_0_TEST, MARKET_WS_TEST)

	// BTCBUSD
	query := "symbol=" + pair
	resp := server.get(EXCHANGE_INFO_PATH + "?" + query)
	var decData ExchangeInfoResp
	decode(resp, &decData)
	decData.ServerTimeStr = utils.ToTime(decData.ServerTime)
	// fmt.Printf("RESPONSE: %+v,\n", utils.ToTime(decData.ServerTime))
	utils.PP(decData)
}

func testTime() {
	// server := createServer(API_0, MARKET_WS_1)
	server := createServer(API_0_TEST, MARKET_WS_TEST)

	resp := server.get(TIME_PATH)
	var decData TimeResp
	decode(resp, &decData)
	fmt.Printf("RESPONSE: %v, %T\n", decData.ServerTime, decData)
}

func decode(resp []uint8, respInstance any) error {
	dec := json.NewDecoder(bytes.NewReader(resp))
	for {
		err := dec.Decode(&respInstance)
		if err == io.EOF {
			break
		} else if err != nil {
			lg.Errorf("Error decoding response %s", err)
			return err
		}
	}
	return nil
}

func genericDecoder(data []uint8) map[string]interface{} {
	var decData map[string]interface{}
	err := json.Unmarshal(data, &decData)
	if err != nil {
		lg.Errorf("Error decoding response %s", err)
		return nil
	}
	return decData
}

func (s *Server) get(path string) []uint8 {
	c := http.Client{Timeout: s.Timeout}
	url := s.APIEndpoint(path)
	lg.Infof("URL:%s", url)
	resp, err := c.Get(url)
	if err != nil {
		lg.Errorf("Error connecting to server %s", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		lg.Errorf("Error reading response %s", err)
		return nil
	}
	return body
}

func (s *Server) APIEndpoint(path string) string {
	return fmt.Sprintf("%s%s", s.BaseAPI, path)
}
