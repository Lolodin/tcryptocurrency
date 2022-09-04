package client

import (
	"cryptobot/internal/cryptocurrency"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type CoinGecko struct {
	Client *http.Client
	Pool   *SymbolPool
}

func (p *CoinGecko) GetCryptocurrency(symbol string) (*cryptocurrency.Metadata, error) {
	if s, ok := p.Pool.PoolSymbol.Load().(map[string]MetaData)[symbol]; ok {
		symbol = strings.ToLower(s.Id)
	}
	b, err := p.Client.Get("https://api.coingecko.com/api/v3/coins/" + symbol + "?localization=false&tickers=true&market_data=false&community_data=false&developer_data=false&sparkline=false")
	log.Println(err)
	if err != nil || b.StatusCode != 200 {
		log.Println(err, 200)
		return nil, err
	}
	data, err := ioutil.ReadAll(b.Body)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	m := &T2{}
	err = json.Unmarshal(data, m)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	for _, ticker := range m.Tickers {
		if ticker.Target != "USDT" {
			continue
		}
		price := ticker.Last
		value := 0.0
		for _, s := range m.Tickers {
			value += s.ConvertedVolume.Usd
		}
		log.Println("Ok CoinGecko")
		return &cryptocurrency.Metadata{
			cryptocurrency.Price(price),
			symbol,
			0,
			value,
		}, nil
	}
	log.Println(err)
	return nil, errors.New("not found")
}

type T2 struct {
	Id              string      `json:"id"`
	Symbol          string      `json:"symbol"`
	Name            string      `json:"name"`
	AssetPlatformId interface{} `json:"asset_platform_id"`
	Platforms       struct {
		Field1 string `json:""`
	} `json:"platforms"`
	BlockTimeInMinutes int           `json:"block_time_in_minutes"`
	HashingAlgorithm   string        `json:"hashing_algorithm"`
	Categories         []string      `json:"categories"`
	PublicNotice       interface{}   `json:"public_notice"`
	AdditionalNotices  []interface{} `json:"additional_notices"`
	Description        struct {
		En string `json:"en"`
	} `json:"description"`
	Links struct {
		Homepage                    []string    `json:"homepage"`
		BlockchainSite              []string    `json:"blockchain_site"`
		OfficialForumUrl            []string    `json:"official_forum_url"`
		ChatUrl                     []string    `json:"chat_url"`
		AnnouncementUrl             []string    `json:"announcement_url"`
		TwitterScreenName           string      `json:"twitter_screen_name"`
		FacebookUsername            string      `json:"facebook_username"`
		BitcointalkThreadIdentifier interface{} `json:"bitcointalk_thread_identifier"`
		TelegramChannelIdentifier   string      `json:"telegram_channel_identifier"`
		SubredditUrl                string      `json:"subreddit_url"`
		ReposUrl                    struct {
			Github    []string      `json:"github"`
			Bitbucket []interface{} `json:"bitbucket"`
		} `json:"repos_url"`
	} `json:"links"`
	Image struct {
		Thumb string `json:"thumb"`
		Small string `json:"small"`
		Large string `json:"large"`
	} `json:"image"`
	CountryOrigin                string  `json:"country_origin"`
	GenesisDate                  string  `json:"genesis_date"`
	SentimentVotesUpPercentage   float64 `json:"sentiment_votes_up_percentage"`
	SentimentVotesDownPercentage float64 `json:"sentiment_votes_down_percentage"`
	MarketCapRank                int     `json:"market_cap_rank"`
	CoingeckoRank                int     `json:"coingecko_rank"`
	CoingeckoScore               float64 `json:"coingecko_score"`
	DeveloperScore               float64 `json:"developer_score"`
	CommunityScore               float64 `json:"community_score"`
	LiquidityScore               float64 `json:"liquidity_score"`
	PublicInterestScore          float64 `json:"public_interest_score"`
	PublicInterestStats          struct {
		AlexaRank   int         `json:"alexa_rank"`
		BingMatches interface{} `json:"bing_matches"`
	} `json:"public_interest_stats"`
	StatusUpdates []interface{} `json:"status_updates"`
	LastUpdated   time.Time     `json:"last_updated"`
	Tickers       []struct {
		Base   string `json:"base"`
		Target string `json:"target"`
		Market struct {
			Name                string `json:"name"`
			Identifier          string `json:"identifier"`
			HasTradingIncentive bool   `json:"has_trading_incentive"`
		} `json:"market"`
		Last            float64 `json:"last"`
		Volume          float64 `json:"volume"`
		ConvertedVolume struct {
			Btc float64 `json:"btc"`
			Eth float64 `json:"eth"`
			Usd float64 `json:"usd"`
		} `json:"converted_volume"`
		TrustScore             string      `json:"trust_score"`
		BidAskSpreadPercentage float64     `json:"bid_ask_spread_percentage"`
		Timestamp              time.Time   `json:"timestamp"`
		LastTradedAt           time.Time   `json:"last_traded_at"`
		LastFetchAt            time.Time   `json:"last_fetch_at"`
		IsAnomaly              bool        `json:"is_anomaly"`
		IsStale                bool        `json:"is_stale"`
		TradeUrl               *string     `json:"trade_url"`
		TokenInfoUrl           interface{} `json:"token_info_url"`
		CoinId                 string      `json:"coin_id"`
		TargetCoinId           string      `json:"target_coin_id,omitempty"`
	} `json:"tickers"`
}
