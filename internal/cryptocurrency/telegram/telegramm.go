package telegram

import (
	"cryptobot/internal/cryptocurrency"
	"cryptobot/internal/cryptocurrency/client"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

type Currency interface {
	GetCryptocurrency(symbol string) (*cryptocurrency.Metadata, error)
}

type Telegram struct {
	Clients []Currency
	Symbols *client.SymbolPool
	Bot     *tgbotapi.BotAPI
}

func NewTelegram(bot *tgbotapi.BotAPI, clients []Currency, symbols *client.SymbolPool) *Telegram {

	return &Telegram{Bot: bot, Clients: clients, Symbols: symbols}
}

func (t *Telegram) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

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
				if val, ok := t.Symbols.PoolSymbol[strings.ToLower(symbol)]; ok && val.Price != 0 {
					data := cryptocurrency.Metadata{Name: val.Name, USDT: cryptocurrency.Price(val.Price), Vol: val.Vol, Change: val.Change}
					msg.Text = data.String()
					break end
				}
				for _, currency := range t.Clients {
					data, err := currency.GetCryptocurrency(symbol)
					if err == nil && data != nil && data.USDT != 0 {
						if val, ok := t.Symbols.PoolSymbol[strings.ToLower(symbol)]; ok {
							data.Name = val.Name
						}

						msg.Text = data.String()
						break end
					}
				}

				log.Println("currency not found")
				continue
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
