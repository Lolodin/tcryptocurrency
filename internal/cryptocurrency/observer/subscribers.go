package observer

import (
	"cryptobot/internal/cryptocurrency/client"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
)

type SubscriberData struct {
	ChactID  int64
	Factor   float64
	Symbols  string
	NewValue float64
	OldValue float64
}

type key struct {
	ChactID int64
	Factor  float64
	Symbols string
}

func (s *SubscriberData) String() string {
	sb := strings.Builder{}
	sb.WriteString("Movement <b>")
	sb.WriteString(s.Symbols)
	sb.WriteString("</b>\n")
	sb.WriteString("Old Price: ")
	sb.WriteString(strconv.FormatFloat(s.OldValue, 'f', 12, 64))
	sb.WriteString("\n")
	sb.WriteString(`New Price: `)
	sb.WriteString(strconv.FormatFloat(s.NewValue, 'f', 12, 64))
	sb.WriteString("\n")
	sb.WriteString(`Factor: `)
	sb.WriteString(strconv.FormatFloat(s.Factor, 'f', 2, 64))
	sb.WriteString("%")

	return sb.String()
}

type Observer struct {
	Subscribers map[key]*SubscriberData
	*client.SymbolPool
	EventChan chan SubscriberData
	M         sync.Mutex
}

func NewObserver(symbolPool *client.SymbolPool, eventChan chan SubscriberData) *Observer {
	return &Observer{Subscribers: map[key]*SubscriberData{}, SymbolPool: symbolPool, EventChan: eventChan}
}

func (o *Observer) Subscribe(data *SubscriberData) {
	key := key{ChactID: data.ChactID, Symbols: data.Symbols}
	o.M.Lock()
	defer o.M.Unlock()
	o.Subscribers[key] = data
}

func (o *Observer) Unsubscribe(data *SubscriberData) {
	key := key{ChactID: data.ChactID, Symbols: data.Symbols}
	o.M.Lock()
	defer o.M.Unlock()
	delete(o.Subscribers, key)
}

func (o *Observer) Service() error {
	for _, subscribers := range o.Subscribers {
		v, ok := o.PoolSymbol.Load().(map[string]*client.MetaData)[subscribers.Symbols]
		if ok {
			if v.Price == 0 {
				log.Println("price empty")
				continue
			}
			subscribers.NewValue = v.Price
			diff := math.Abs(subscribers.OldValue - subscribers.NewValue)
			f := subscribers.OldValue * subscribers.Factor
			log.Println("diff", diff, "factor", f)
			//Если движение выше процента установленного в конфиге, то оповещаем канал
			if diff > f {
				o.EventChan <- *subscribers
			}

		} else {
			o.Unsubscribe(subscribers)
			continue
		}
		subscribers.OldValue = subscribers.NewValue
		subscribers.NewValue = 0
	}

	return nil

}
