package test

import (
	"deribit-connector/pkg/deribit"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testInstrumentName = "testInst"
)

func TestBids(t *testing.T) {
	orderbook := deribit.NewOrderBook(testInstrumentName)
	orderbook.AddBid(1.1, 100)
	orderbook.AddBid(0.9, 101)
	orderbook.AddBid(1.2, 99)
	orderbook.AddBid(0.8, 102)

	assert.Equal(t, deribit.Level{Price: 1.2, Quantity: 99}, *orderbook.GetBestBid())

	expectedBidsList := []deribit.Level{
		{Price: 1.2, Quantity: 99}, {Price: 1.1, Quantity: 100},
		{Price: 0.9, Quantity: 101}, {Price: 0.8, Quantity: 102},
	}
	actualBidsList, _ := orderbook.ToList()
	assert.Equal(t, len(expectedBidsList), len(actualBidsList))
	for i := 0; i < len(expectedBidsList); i++ {
		assert.Equal(t, expectedBidsList[i], *actualBidsList[i])
	}

	orderbook.RemoveBid(0.9)
	expectedBidsListAfterRemoval := []deribit.Level{
		{Price: 1.2, Quantity: 99}, {Price: 1.1, Quantity: 100},
		{Price: 0.8, Quantity: 102},
	}
	actualBidsListAfterRemoval, _ := orderbook.ToList()
	assert.Equal(t, len(expectedBidsListAfterRemoval), len(actualBidsListAfterRemoval))
	for i := 0; i < len(expectedBidsListAfterRemoval); i++ {
		assert.Equal(t, expectedBidsListAfterRemoval[i], *actualBidsListAfterRemoval[i])
	}
}

func TestAsks(t *testing.T) {
	orderbook := deribit.NewOrderBook(testInstrumentName)
	orderbook.AddAsk(1.1, 100)
	orderbook.AddAsk(0.9, 101)
	orderbook.AddAsk(1.2, 99)
	orderbook.AddAsk(0.8, 102)

	assert.Equal(t, deribit.Level{Price: 0.8, Quantity: 102}, *orderbook.GetBestAsk())

	expectedAsksList := []deribit.Level{
		{Price: 0.8, Quantity: 102}, {Price: 0.9, Quantity: 101},
		{Price: 1.1, Quantity: 100}, {Price: 1.2, Quantity: 99},
	}
	_, actualAsksList := orderbook.ToList()
	assert.Equal(t, len(expectedAsksList), len(actualAsksList))
	for i := 0; i < len(expectedAsksList); i++ {
		assert.Equal(t, expectedAsksList[i], *actualAsksList[i])
	}

	orderbook.RemoveAsk(0.9)
	expectedAsksListAfterRemoval := []deribit.Level{
		{Price: 0.8, Quantity: 102}, {Price: 1.1, Quantity: 100},
		{Price: 1.2, Quantity: 99},
	}
	_, actualAsksListAfterRemoval := orderbook.ToList()
	assert.Equal(t, len(expectedAsksListAfterRemoval), len(actualAsksListAfterRemoval))
	for i := 0; i < len(expectedAsksListAfterRemoval); i++ {
		assert.Equal(t, expectedAsksListAfterRemoval[i], *actualAsksListAfterRemoval[i])
	}
}

func TestJSON(t *testing.T) {
	orderbook := deribit.NewOrderBook(testInstrumentName)
	orderbook.AddAsk(1.2, 100)
	orderbook.AddAsk(1.3, 101)
	orderbook.AddAsk(1.1, 99)
	orderbook.AddAsk(1.0, 102)
	orderbook.AddBid(0.7, 100)
	orderbook.AddBid(0.9, 101)
	orderbook.AddBid(0.6, 99)
	orderbook.AddBid(0.8, 102)

	orderbookBytes, err := json.Marshal(orderbook)
	assert.NoError(t, err)

	var newOrderbook deribit.OrderBook
	err = json.Unmarshal(orderbookBytes, &newOrderbook)
	assert.NoError(t, err)

	assert.Equal(t, orderbook.Timestamp, newOrderbook.Timestamp)
	assert.Equal(t, orderbook.InstrumentName, newOrderbook.InstrumentName)

	obBids, obAsks := orderbook.ToList()
	newObBids, newObAsks := newOrderbook.ToList()

	assert.Equal(t, len(obBids), len(newObBids))
	for i := 0; i < len(obBids); i++ {
		assert.Equal(t, *obBids[i], *newObBids[i])
	}

	assert.Equal(t, len(obAsks), len(newObAsks))
	for i := 0; i < len(obAsks); i++ {
		assert.Equal(t, *obAsks[i], *newObAsks[i])
	}

}
