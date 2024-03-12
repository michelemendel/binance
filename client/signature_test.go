package client

import (
	"testing"
)

// Some made up secret key
const secretKey = "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j"

type args struct {
	queryString string
	body        string
	timestamp   int64
}

func TestSignature(t *testing.T) {
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{name: "RequestBody",
			args: args{
				queryString: "",
				body:        "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000",
				timestamp:   1499827319559,
			},
			expected: "c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71",
		},
		{name: "QueryString",
			args: args{
				queryString: "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000",
				body:        "",
				timestamp:   1499827319559,
			},
			expected: "c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71",
		},
		{name: "RequestBodyAndQueryString",
			args: args{
				queryString: "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC",
				body:        "quantity=1&price=0.1&recvWindow=5000",
				timestamp:   1499827319559,
			},
			expected: "0fd168b8ddb4876a0358a8d14d0c9f3da0e9b20c5d52b2a00fcf7d1c602f9a77",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Signature(secretKey, tt.args.queryString, tt.args.body, tt.args.timestamp)
			expected := tt.expected
			if actual != expected {
				t.Errorf("Signature() = %v, want %v", actual, expected)
			}
		})
	}
}
