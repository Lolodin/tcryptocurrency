package client

import (
	"cryptobot/internal/cryptocurrency"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

const api = "https://pro-api.coinmarketcap.com/v1/tools/price-conversion"
const key = "eb2b296e-6962-4656-9f6c-d5a2d432145a"
const header = "X-CMC_PRO_API_KEY"

type CoinMarketcup struct {
	Client *http.Client
}

func NewCoinMarketCupClient() *CoinMarketcup {
	client := &http.Client{Timeout: 200 * time.Millisecond}
	return &CoinMarketcup{Client: client}
}

func (p *CoinMarketcup) GetCryptocurrency(symbol string) (*cryptocurrency.Metadata, error) {
	q := url.Values{}
	q.Add("symbol", symbol)
	q.Add("amount", "1")

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}
	fillReq(req, q)
	res, err := p.Client.Do(req)
	if err != nil || res.StatusCode != 200 {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	m := &Model{}
	err = json.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}

	data := &cryptocurrency.Metadata{
		USDT: cryptocurrency.Price(m.Data.Quote.USD.Price),
		Name: m.Data.Name,
	}

	return data, nil
}

func fillReq(req *http.Request, q url.Values) {
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accepts", "application/json")
	req.Header.Add(header, key)
}

type Model struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
		Notice       interface{} `json:"notice"`
	} `json:"status"`
	Data struct {
		Id          int       `json:"id"`
		Symbol      string    `json:"symbol"`
		Name        string    `json:"name"`
		Amount      int       `json:"amount"`
		LastUpdated time.Time `json:"last_updated"`
		Quote       struct {
			USD struct {
				Price       float64   `json:"price"`
				LastUpdated time.Time `json:"last_updated"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

type SymbolPool struct {
	Client     *http.Client
	PoolSymbol atomic.Value
}

func NewSymbolPool() *SymbolPool {
	pool := &SymbolPool{Client: &http.Client{}, PoolSymbol: atomic.Value{}}
	err := pool.Update()
	if err != nil {
		log.Panicln(err)
	}
	go func(pool *SymbolPool) {
		ticker := time.NewTicker(30 * time.Minute)
		for {
			select {
			case <-ticker.C:
				err := pool.Update()
				if err != nil {
					log.Println("error update symbol pool")
				}
			}
		}
	}(pool)

	return pool
}

func (p *SymbolPool) Update() error {
	log.Println("update pool")
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/list", nil)
	if err != nil {
		return err
	}

	res, err := p.Client.Do(req)
	if err != nil || res.StatusCode != 200 {
		return errors.New("bad request")
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return p.update(err, b)
}

func (p *SymbolPool) update(err error, b []byte) error {
	m := []ParseDataID{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	idlist := make([]string, 0, len(m))
	temp := make(map[string]*MetaData, len(m))
	for _, datum := range m {
		if _, ok := temp[strings.ToLower(datum.Symbol)]; !ok {
			data := &MetaData{
				Id:   datum.Id,
				Name: datum.Name,
			}

			idlist = append(idlist, datum.Id)
			temp[strings.ToLower(datum.Symbol)] = data
			temp[strings.ToLower(datum.Id)] = data

		} else {
			if datum.Symbol == datum.Id {
				data := &MetaData{
					Id:   datum.Id,
					Name: datum.Name,
				}
				idlist = append(idlist, datum.Id)
				temp[strings.ToLower(datum.Symbol)] = data
				temp[strings.ToLower(datum.Id)] = data
			}
		}
	}
	for i := 0; i < len(idlist); i += 400 {
		if i+400 > len(idlist) {
			temp, err = p.funcName(idlist[i:], temp)
			if err != nil {
				time.Sleep(200 * time.Second)
				temp, _ = p.funcName(idlist[i:], temp)
			}
			break
		}
		temp, err = p.funcName(idlist[i:i+400], temp)
		if err != nil {
			time.Sleep(200 * time.Second)
			temp, _ = p.funcName(idlist[i:], temp)
		}

	}

	p.PoolSymbol.Store(temp)

	return nil
}

func (p *SymbolPool) funcName(idlist []string, temp map[string]*MetaData) (map[string]*MetaData, error) {
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price", nil)
	if err != nil {
		return temp, err
	}
	u := url.Values{}
	u.Add("ids", strings.TrimSpace(strings.Join(idlist, ",")))
	u.Add("vs_currencies", "usd")
	u.Add("include_24hr_vol", "true")
	u.Add("include_24hr_change", "true")
	req.URL.RawQuery = u.Encode()
	fmt.Println(req.URL.String())
	res, err := p.Client.Do(req)
	if err != nil || res.StatusCode != 200 {
		return temp, errors.New("bad request")
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return temp, err
	}

	info := make(ParseData)
	err = json.Unmarshal(b, &info)
	if err != nil {
		return temp, err
	}
	log.Println("info len", len(info))
	for k, s2 := range info {
		v, ok := temp[k]
		if ok {
			v.Vol = s2.Usd24HVol
			v.Change = s2.Usd24HChange
			v.Price = s2.Usd
		} else {
			el := &MetaData{
				Id:     k,
				Name:   k,
				Vol:    s2.Usd24HVol,
				Change: s2.Usd24HChange,
				Price:  s2.Usd,
			}
			temp[k] = el
		}
	}
	return temp, nil
}

type MetaData struct {
	Id     string
	Name   string
	Vol    float64
	Change float64
	Price  float64
}

type ParseDataID struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type ParseData map[string]struct {
	Usd          float64 `json:"usd"`
	Usd24HVol    float64 `json:"usd_24h_vol"`
	Usd24HChange float64 `json:"usd_24h_change"`
}
