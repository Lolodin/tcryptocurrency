package telegram

import (
	"cryptobot/internal/cryptocurrency"
	"cryptobot/internal/cryptocurrency/client"
	"cryptobot/internal/cryptocurrency/observer"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

type Currency interface {
	GetCryptocurrency(symbol string) (*cryptocurrency.Metadata, error)
}

type Telegram struct {
	Clients         []Currency
	Symbols         *client.SymbolPool
	Bot             *tgbotapi.BotAPI
	Observer        *observer.Observer
	CentralCurrency Currency
	EventChan       chan observer.SubscriberData
}

func NewTelegram(bot *tgbotapi.BotAPI, clients []Currency, symbols *client.SymbolPool, observer2 *observer.Observer, ch chan observer.SubscriberData, cc Currency) *Telegram {

	return &Telegram{Bot: bot, Clients: clients, Symbols: symbols, Observer: observer2, EventChan: ch, CentralCurrency: cc}
}

func (t *Telegram) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	go func() {
		for {
			select {
			case d := <-t.EventChan:
				m := tgbotapi.NewMessage(d.ChactID, d.String())
				m.ParseMode = tgbotapi.ModeHTML
				if _, err := t.Bot.Send(m); err != nil {
					log.Println(err)
				}
			}
		}
	}()
	pool, ok := t.Symbols.PoolSymbol.Load().(map[string]*client.MetaData)
	if !ok {
		log.Panicln("pool error")
	}
	fmt.Println("len pool", len(pool))
	updates := t.Bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if !update.Message.IsCommand() { // ignore any non-command Messages
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = tgbotapi.ModeHTML
		end:
			switch update.Message.Command() {
			case "p":
				symbol := update.Message.CommandArguments()
				pool, ok := t.Symbols.PoolSymbol.Load().(map[string]*client.MetaData)
				if !ok {
					continue
				}
				if val, ok := pool[strings.ToLower(symbol)]; ok && val.Price != 0 {
					for _, currency := range t.Clients {
						data, err := currency.GetCryptocurrency(symbol)
						if err == nil && data != nil && data.USDT != 0 {
							val.Price = float64(data.USDT)
							break
						}

					}
					data := cryptocurrency.Metadata{Name: val.Name, USDT: cryptocurrency.Price(val.Price), Vol: val.Vol, Change: val.Change}
					msg.Text = data.String()
					break end
				}
				for _, currency := range t.Clients {
					data, err := currency.GetCryptocurrency(symbol)
					if err == nil && data != nil && data.USDT != 0 {
						if _, ok := pool[strings.ToLower(symbol)]; ok {
							data.Name = symbol
						}

						msg.Text = data.String()
						break end
					}
				}

				log.Println("currency not found")
				continue
			case "sub":
				str := update.Message.CommandArguments()
				params := strings.Split(str, " ")
				if len(params) < 2 {
					continue
				}
				symbol := strings.ToLower(params[0])
				factor := params[1]
				pool, ok := t.Symbols.PoolSymbol.Load().(map[string]*client.MetaData)
				if !ok {
					continue
				}
				if val, ok := pool[symbol]; ok && val.Price != 0 {
					for _, currency := range t.Clients {
						data, err := currency.GetCryptocurrency(symbol)
						if err == nil && data != nil && data.USDT != 0 {
							val.Price = float64(data.USDT)
							break
						}

					}

					f, err := strconv.ParseFloat(factor, 64)
					if err != nil {
						continue
					}
					d := &observer.SubscriberData{
						OldValue: val.Price,
						ChactID:  update.Message.Chat.ID,
						Symbols:  symbol,
						Factor:   f,
					}
					t.Observer.Subscribe(d)
					msg.Text = "You subscribed: " + val.Name + "|" + strconv.FormatFloat(d.OldValue, 'f', 10, 64) + "| Factor: " + strconv.FormatFloat(f*100, 'f', 2, 64) + "%"
				}
			case "unsub":
				symbol := update.Message.CommandArguments()
				d := &observer.SubscriberData{
					ChactID: update.Message.Chat.ID,
					Symbols: symbol,
				}

				t.Observer.Unsubscribe(d)
				msg.Text = "You unsubscribed"
			case "cur":
				symbol := update.Message.CommandArguments()

				c, err := t.CentralCurrency.GetCryptocurrency(symbol)
				if err != nil {
					msg.Text = err.Error()
					break
				}
				msg.Text = c.String()
			default:
				continue
			}

			if _, err := t.Bot.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}

	return errors.New("stop bot")
}
