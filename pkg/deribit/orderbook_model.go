package deribit

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

type Level struct {
	Price    float64
	Quantity float64
}

type OrderBook struct {
	Bids           *treemap.Map
	Asks           *treemap.Map
	Timestamp      time.Time
	InstrumentName string
}

func NewOrderBook(instrumentName string) *OrderBook {
	return &OrderBook{
		Bids:           treemap.NewWith(bidComparator),
		Asks:           treemap.NewWith(utils.Float64Comparator),
		InstrumentName: instrumentName,
		Timestamp:      time.Now().UTC(),
	}
}

func bidComparator(a, b interface{}) int {
	return -1 * utils.Float64Comparator(a, b)
}

func (ob *OrderBook) AddBid(price, quantity float64) {
	ob.Bids.Put(price, &Level{Price: price, Quantity: quantity})
}

func (ob *OrderBook) RemoveBid(price float64) {
	ob.Bids.Remove(price)
}

func (ob *OrderBook) GetBestBid() *Level {
	if ob.Bids.Empty() {
		return nil
	}
	_, value := ob.Bids.Min()
	return value.(*Level)
}

func (ob *OrderBook) AddAsk(price, quantity float64) {
	ob.Asks.Put(price, &Level{Price: price, Quantity: quantity})
}

func (ob *OrderBook) RemoveAsk(price float64) {
	ob.Asks.Remove(price)
}

func (ob *OrderBook) GetBestAsk() *Level {
	if ob.Asks.Empty() {
		return nil
	}
	_, value := ob.Asks.Min()
	return value.(*Level)
}

func (ob *OrderBook) ToList() ([]*Level, []*Level) {
	Bids := []*Level{}
	Asks := []*Level{}

	ob.Bids.Each(func(key, value interface{}) {
		Bids = append(Bids, value.(*Level))
	})

	ob.Asks.Each(func(key, value interface{}) {
		Asks = append(Asks, value.(*Level))
	})

	return Bids, Asks
}

func (ob *OrderBook) MarshalJSON() ([]byte, error) {
	Bids, Asks := ob.ToList()
	type Alias OrderBook
	return json.Marshal(&struct {
		Bids []*Level `json:"Bids"`
		Asks []*Level `json:"Asks"`
		*Alias
	}{
		Bids:  Bids,
		Asks:  Asks,
		Alias: (*Alias)(ob),
	})
}

func (ob *OrderBook) UnmarshalJSON(body []byte) error {
	type Alias OrderBook
	temp := &struct {
		Bids []*Level `json:"Bids"`
		Asks []*Level `json:"Asks"`
		*Alias
	}{
		Alias: (*Alias)(ob),
	}

	err := json.Unmarshal(body, temp)
	if err != nil {
		return fmt.Errorf("unable to unmarshal orderbook object, err: %w", err)
	}

	ob.Bids = treemap.NewWith(bidComparator)
	for _, level := range temp.Bids {
		ob.Bids.Put(level.Price, level)
	}

	ob.Asks = treemap.NewWith(utils.Float64Comparator)
	for _, level := range temp.Asks {
		ob.Asks.Put(level.Price, level)
	}
	return nil
}
