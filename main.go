package main

import (
	"context"
	"cryptobot/internal/cryptocurrency/client"
	"cryptobot/internal/cryptocurrency/observer"
	"cryptobot/internal/cryptocurrency/telegram"
	"cryptobot/internal/scheduler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

//scp -P 3333 main root@185.237.218.45:/
func main() {
	bot, err := tgbotapi.NewBotAPI("5566392860:AAGHHzZNoVuCSBrPc1LdXch2U8ZSX-iDwxU")
	if err != nil {
		log.Panicln(err)
	}
	bot.Debug = true
	ch := make(chan observer.SubscriberData)

	sl := make([]telegram.Currency, 0)
	pool := client.NewSymbolPool()

	observer := observer.NewObserver(pool, ch)
	s := scheduler.Scheduler{}
	s.SetupJob(time.Minute*10, pool.Update)
	s.SetupJob(time.Minute*15, observer.Service)
	err = s.Run(context.Background())
	if err != nil {
		log.Panicln(err)
	}

	sl = append(sl, client.NewBinanceClient(), client.NewPoliniex(), client.NewCoinMarketCupClient())

	app := telegram.NewTelegram(bot, sl, pool, observer, ch)
l:
	if err := app.Run(); err != nil {
		log.Println(err)
		goto l
	}
}
