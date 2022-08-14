package main

import (
	"cryptobot/internal/cryptocurrency/client"
	"cryptobot/internal/cryptocurrency/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

//scp -P 3333 main root@185.237.218.45:/
func main() {
	bot, err := tgbotapi.NewBotAPI("5566392860:AAGHHzZNoVuCSBrPc1LdXch2U8ZSX-iDwxU")
	if err != nil {
		log.Panicln(err)
	}
	bot.Debug = true
	sl := make([]telegram.Currency, 0)
	pool := client.NewSymbolPool()
	sl = append(sl, &client.CoinGecko{Client: &http.Client{}, Pool: pool}, client.NewBinanceClient(), client.NewPoliniex(), client.NewCoinMarketCupClient())

	app := telegram.NewTelegram(bot, sl, pool)
l:
	if err := app.Run(); err != nil {
		log.Println(err)
		goto l
	}
}
