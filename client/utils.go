package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"

	"github.com/michelemendel/binance/util"
)

func (client *Client) Get(path, query string) []uint8 {
	var url string
	isSAPI, _ := regexp.MatchString("/sapi/", path)

	if isSAPI {
		ts := util.TimeNowInMillis()
		signature := Signature(client.SecretKey, query, "", ts)
		url = client.SAPIEndpoint(path, query, signature, ts)
	} else {
		url = client.APIEndpoint(path, query)
	}

	slog.Info("connection to server", "url", url)

	httpClient := &http.Client{Timeout: client.Timeout}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("error creating request", "error", err)
		// return
	}

	req.Header.Set("X-MBX-APIKEY", client.APIKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		slog.Error("error making request", "url", url, "error", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response", "error", err)
		return nil
	}
	return body
}

func (c *Client) APIEndpoint(path, query string) string {
	base := fmt.Sprintf("%s%s", c.BaseAPI, path)

	if query != "" {
		base = base + fmt.Sprintf("?%s", query)
	}

	return base
}

func (c *Client) SAPIEndpoint(path, query, signature string, timestamp int64) string {
	base := fmt.Sprintf("%s%s?timestamp=%v", c.BaseAPI, path, timestamp)

	if signature != "" {
		base = base + fmt.Sprintf("&signature=%s", signature)
	}

	if query != "" {
		base = base + fmt.Sprintf("&%s", query)
	}

	return base
}

func Signature(secretKey, query, body string, timestamp int64) string {
	var payload string

	if query != "" {
		payload = query + body + fmt.Sprintf("&timestamp=%d", timestamp)
	} else {
		payload = fmt.Sprintf("timestamp=%d", timestamp)
	}

	signed := HMACSign(secretKey, payload)
	urlEncoded := urlEncode(signed)
	return urlEncoded
}

// todo: Need this, do we?
// func base64Encode(str string) string {
// 	return base64.StdEncoding.EncodeToString([]byte(str))
// }

func HMACSign(secret, data string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func urlEncode(str string) string {
	return url.QueryEscape(str)
}

// func makePayload(params any) string {
// 	payload, err := json.Marshal(params)
// 	if err != nil {
// 		slog.Error("error marshalling payload", "error", err)
// 		return ""
// 	}
// 	return string(payload)
// }

func decode(resp []uint8, respInstance any) error {
	err := json.Unmarshal(resp, &respInstance)
	if err != nil {
		slog.Error("error decoding response", "error", err)
		return err
	}
	return nil
}
