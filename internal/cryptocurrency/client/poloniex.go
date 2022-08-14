package client

import (
	"cryptobot/internal/cryptocurrency"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Poloniex struct {
	Client *http.Client
}

func NewPoliniex() *Poloniex {
	return &Poloniex{Client: &http.Client{}}
}

func (p *Poloniex) GetCryptocurrency(symbol string) (*cryptocurrency.Metadata, error) {
	b, err := p.Client.Get("https://api.poloniex.com/markets/" + symbol + "_usdt/price")
	if err != nil || b.StatusCode != 200 {
		return nil, err
	}
	data, err := ioutil.ReadAll(b.Body)
	if err != nil {
		return nil, err
	}
	m := &ModelPolo{}
	err = json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}

	price, _ := strconv.ParseFloat(m.Price, 64)
	if price == 0 {
		return nil, errors.New("not data")
	}
	change, _ := strconv.ParseFloat(m.DailyChange, 64)
	return &cryptocurrency.Metadata{
		cryptocurrency.Price(price),
		symbol,
		change * 100,
		0,
	}, nil
}

type ModelPolo struct {
	Symbol      string `json:"symbol"`
	Price       string `json:"price"`
	Time        int64  `json:"time"`
	DailyChange string `json:"dailyChange"`
	Ts          int64  `json:"ts"`
}
