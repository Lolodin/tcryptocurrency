package client

import (
	"fmt"
	"testing"
)

func TestCoinMarketCupClient_GetCryptocurrency(t *testing.T) {
	client := NewCoinMarketCupClient()
	data, err := client.GetCryptocurrency("mex")
	fmt.Println(data, err, data.USDT.String())
}
func TestSymbolPool_Update(t *testing.T) {
	pool := NewSymbolPool()
	fmt.Println(len(pool.PoolSymbol))
}
