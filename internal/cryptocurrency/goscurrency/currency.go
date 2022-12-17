package goscurrency

import (
	"cryptobot/internal/cryptocurrency"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
)

const key = "20NY8z0bFl5Md8rmJ1vBjPvV0rEw4vQ0"
const url = "https://api.apilayer.com/currency_data/list"

type T struct {
	Success   bool               `json:"success"`
	Timestamp int                `json:"timestamp"`
	Source    string             `json:"source"`
	Quotes    map[string]float64 `json:"quotes"`
}
type Currency struct {
	MapWithCurrency *atomic.Value
	Client          *http.Client
	Url             string
}

func NewCurrency() *Currency {
	return &Currency{Url: url, Client: http.DefaultClient, MapWithCurrency: &atomic.Value{}}
}

func ByteTomap(d []byte) (map[string]float64, error) {
	t := &T{}
	err := json.Unmarshal(d, t)
	if err != nil || len(t.Quotes) == 0 {
		return nil, err
	}
	m := map[string]float64{}
	for s, f := range t.Quotes {
		key := strings.ToLower(s)
		key = strings.TrimPrefix(key, "usd")
		m[key] = f
	}

	return m, nil
}
func (p *Currency) Update() error {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	r.Header.Set("apikey", key)

	b, err := p.Client.Do(r)
	if err != nil {
		return err
	}

	bs, _ := ioutil.ReadAll(b.Body)
	m, err := ByteTomap(bs)
	if err != nil {
		return err
	}
	p.MapWithCurrency.Store(m)
	return nil
}

func (p *Currency) GetMap() (map[string]float64, bool) {
	val, ok := p.MapWithCurrency.Load().(map[string]float64)
	return val, ok
}

func (p *Currency) GetCryptocurrency(symbol string) (*cryptocurrency.Metadata, error) {
	c := &cryptocurrency.Metadata{}
	mPrice, ok := p.GetMap()
	if !ok {
		return nil, fmt.Errorf("value is empty")
	}
	price, ok := mPrice[symbol]
	if !ok {
		return nil, fmt.Errorf("can't find currency at map")
	}
	c.Name = symbol
	c.USDT = cryptocurrency.Price(price)
	return c, nil
}
