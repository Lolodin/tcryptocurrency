package client

import (
	"cryptobot/internal/cryptocurrency"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type BinanceClient struct {
	Client *http.Client
}

func NewBinanceClient() *BinanceClient {
	return &BinanceClient{Client: &http.Client{Timeout: 200 * time.Millisecond}}
}

func (p BinanceClient) GetCryptocurrency(symbol string) (*cryptocurrency.Metadata, error) {
	b, err := p.Client.Get("https://www.binance.com/api/v3/ticker/24hr?symbol=" + symbol + "USDT")
	if err != nil || b.StatusCode != 200 {
		return nil, err
	}
	data, err := ioutil.ReadAll(b.Body)
	if err != nil {
		return nil, err
	}
	m := &ModelBinance{}
	err = json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}

	price, _ := strconv.ParseFloat(m.LastPrice, 64)
	if price == 0 {
		return nil, errors.New("not data")
	}
	change, _ := strconv.ParseFloat(m.PriceChangePercent, 64)
	return &cryptocurrency.Metadata{
		cryptocurrency.Price(price),
		symbol,
		change,
		0,
	}, nil
}

type ModelBinance struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstId            int    `json:"firstId"`
	LastId             int    `json:"lastId"`
	Count              int    `json:"count"`
}
