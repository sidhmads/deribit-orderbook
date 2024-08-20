package deribit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

const (
	pingDuration                        = time.Second * 10
	reconnectInterval                   = time.Second * 2
	subscriptionMethod                  = "subscription"
	maxInstrumentPerSubscriptionRequest = 500
)

type GetOrderbook struct {
	*Deribit
	*baseWS
	producer                    sarama.SyncProducer
	instrumentNames             []string
	instrumentKafkaTopicMapping map[string]string
	orderbookMapping            map[string]*OrderBook
}

type OrderBookFeedRequestMsg struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		Channels []string `json:"channels"`
	} `json:"params"`
}

type OrderBookResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int64  `json:"id"`
	Testnet bool   `json:"testnet"`
	Method  string `json:"method"`
	Params  struct {
		Data struct {
			Type           string          `json:"type"`
			Timestamp      int64           `json:"timestamp"`
			InstrumentName string          `json:"instrument_name"`
			ChangeID       int64           `json:"change_id"`
			PrevChangeID   int64           `json:"prev_change_id"`
			Bids           [][]interface{} `json:"bids"`
			Asks           [][]interface{} `json:"asks"`
		} `json:"data"`
		Channel string `json:"channel"`
	} `json:"params"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result []string `json:"result"`
	UsIn   int64    `json:"usIn"`
	UsOut  int64    `json:"usOut"`
	UsDiff int      `json:"usDiff"`
}

func NewOrderbook(d *Deribit, instruments []Instrument) (*GetOrderbook, error) {
	instrumentKafkaTopicMapping := make(map[string]string)
	instrumentNames := []string{}
	for _, inst := range instruments {
		instrumentNames = append(instrumentNames, inst.InstrumentName)
		instrumentKafkaTopicMapping[inst.InstrumentName] = createOrderbookTopic(inst.SettlementCurrency, inst.Kind)
	}
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{d.config.KAFKA_SERVER_ADDRESS}, producerConfig)
	if err != nil {
		return nil, err
	}
	return &GetOrderbook{
		Deribit:                     d,
		baseWS:                      nil,
		instrumentNames:             instrumentNames,
		orderbookMapping:            make(map[string]*OrderBook),
		instrumentKafkaTopicMapping: instrumentKafkaTopicMapping,
		producer:                    producer,
	}, nil
}

func (o *GetOrderbook) StreamOrderbooks(ctx context.Context) error {
	if len(o.instrumentNames) == 0 {
		log.Println("There are no instruments to subscribe")
		return nil
	}
	o.baseWS = newWebsocket(wssSchema, o.config.WS_HOST, o.config.WS_PATH, o.subscribeOrderbooks, o.processOrderbookEvent, pingDuration, nil, reconnectInterval)
	o.baseWS.startStreaming(ctx)
	return nil
}

func (o *GetOrderbook) subscribeOrderbooks() error {
	partitionedInstruments := splitToBatches(o.instrumentNames, maxInstrumentPerSubscriptionRequest)

	var err error
	for _, instrumentPartition := range partitionedInstruments {
		var channels []string
		for _, inst := range instrumentPartition {
			channels = append(channels, o.createChannelName(inst))
		}

		requestMsg := &OrderBookFeedRequestMsg{}
		requestMsg.Jsonrpc = "2.0"
		requestMsg.Method = "public/subscribe"
		requestMsg.Params.Channels = channels
		err := o.baseWS.conn.WriteJSON(requestMsg)
		if err != nil {
			return err
		}
	}
	return err
}

func (o *GetOrderbook) createChannelName(instrumentName string) string {
	return "book." + instrumentName + "." + o.config.ORDERBOOK_INTERVAL
}

func (o *GetOrderbook) processOrderbookEvent(msg interface{}) error {
	var orderbookResponse OrderBookResponse
	err := json.Unmarshal(msg.([]byte), &orderbookResponse)
	if err != nil {
		log.Printf("Error while trying to unmarshal orderbook response from server, err: %s\n", err.Error())
		return err
	}
	if orderbookResponse.Method == subscriptionMethod {
		o.processOrderbook(&orderbookResponse)
	} else if len(orderbookResponse.Result) != 0 {
		log.Printf("Received successful subscription message, subscribed to %d channels\n", len(orderbookResponse.Result))
	} else if orderbookResponse.Error.Message != "" {
		errMessage := fmt.Errorf("Error received from deribit websocket subscription, err: code: %d, message: %s", orderbookResponse.Error.Code, orderbookResponse.Error.Message)
		log.Println(errMessage.Error())
		return errMessage
	} else {
		log.Println("Unknown orderbook response received")
	}
	return nil
}

func (o *GetOrderbook) processOrderbook(orderbookResponse *OrderBookResponse) {
	if orderbookResponse.Params.Data.InstrumentName == "" {
		return
	}

	if _, ok := o.orderbookMapping[orderbookResponse.Params.Data.InstrumentName]; !ok {
		o.orderbookMapping[orderbookResponse.Params.Data.InstrumentName] = NewOrderBook(orderbookResponse.Params.Data.InstrumentName)
	}

	orderbook := o.orderbookMapping[orderbookResponse.Params.Data.InstrumentName]
	for _, bidData := range orderbookResponse.Params.Data.Bids {
		bidAction := bidData[0].(string)
		bidPrice := bidData[1].(float64)
		bidAmount := bidData[2].(float64)

		if bidAction == "new" {
			orderbook.AddBid(bidPrice, bidAmount)
		} else if bidAction == "change" {
			orderbook.RemoveBid(bidPrice)
			orderbook.AddBid(bidPrice, bidAmount)
		} else if bidAction == "delete" {
			orderbook.RemoveBid(bidPrice)
		}
	}

	for _, askData := range orderbookResponse.Params.Data.Asks {
		askAction := askData[0].(string)
		askPrice := askData[1].(float64)
		askAmount := askData[2].(float64)

		if askAction == "new" {
			orderbook.AddAsk(askPrice, askAmount)
		} else if askAction == "change" {
			orderbook.RemoveAsk(askPrice)
			orderbook.AddAsk(askPrice, askAmount)
		} else if askAction == "delete" {
			orderbook.RemoveAsk(askPrice)
		}
	}
	orderbook.Timestamp = time.UnixMilli(orderbookResponse.Params.Data.Timestamp).UTC()
	o.produceKafkaMessage(orderbook)
}

func (o *GetOrderbook) produceKafkaMessage(ob *OrderBook) {
	orderbookBytes, err := json.Marshal(ob)
	if err != nil {
		log.Printf("Unable to convert Orderbook to bytes, err: %v\n", err.Error())
		return
	}
	log.Println(string(orderbookBytes))
	if topic, ok := o.instrumentKafkaTopicMapping[ob.InstrumentName]; ok {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(ob.InstrumentName),
			Value: sarama.ByteEncoder(orderbookBytes),
		}

		partition, offset, err := o.producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message: %s\n", err)
		} else {
			log.Printf("Message sent to partition %d with offset %d\n", partition, offset)
		}
	}

}
