package main

import (
	"cryptobot/internal/cryptocurrency/client"
	"cryptobot/internal/cryptocurrency/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

//scp -P 3333 main root@185.237.218.45:/
func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panicln(err)
	}
	bot.Debug = true
	sl := make([]telegram.Currency, 0)
	pool := client.NewSymbolPool()
	sl = append(sl, client.NewBinanceClient(), client.NewPoliniex(), client.NewCoinMarketCupClient())

	app := telegram.NewTelegram(bot, sl, pool)
l:
	if err := app.Run(); err != nil {
		log.Println(err)
		goto l
	}
}
