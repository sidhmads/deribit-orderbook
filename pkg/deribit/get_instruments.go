package deribit

import (
	"fmt"
)

type GetInstruments struct {
	*Deribit
}

type GetInstrumentsResponseResult struct {
	Jsonrpc string       `json:"jsonrpc"`
	UsIn    int64        `json:"usIn"`
	UsOut   int64        `json:"usOut"`
	UsDiff  int          `json:"usDiff"`
	Testnet bool         `json:"testnet"`
	Result  []Instrument `json:"result"`
}

type Instrument struct {
	TickSize                 float64            `json:"tick_size"`
	TakerCommission          float64            `json:"taker_commission"`
	SettlementPeriod         string             `json:"settlement_period"`
	SettlementCurrency       string             `json:"settlement_currency"`
	Rfq                      bool               `json:"rfq"`
	QuoteCurrency            string             `json:"quote_currency"`
	PriceIndex               string             `json:"price_index"`
	MinTradeAmount           float64            `json:"min_trade_amount"`
	MaxLiquidationCommission float64            `json:"max_liquidation_commission"`
	MaxLeverage              float64            `json:"max_leverage"`
	MakerCommission          float64            `json:"maker_commission"`
	Kind                     string             `json:"kind"`
	IsActive                 bool               `json:"is_active"`
	InstrumentName           string             `json:"instrument_name"`
	InstrumentID             float64            `json:"instrument_id"`
	FutureType               string             `json:"future_type"`
	ExpirationTimestamp      int64              `json:"expiration_timestamp"`
	CreationTimestamp        int64              `json:"creation_timestamp"`
	CounterCurrency          string             `json:"counter_currency"`
	ContractSize             float64            `json:"contract_size"`
	BlockTradeTickSize       float64            `json:"block_trade_tick_size"`
	BlockTradeMinTradeAmount float64            `json:"block_trade_min_trade_amount"`
	BlockTradeCommission     float64            `json:"block_trade_commission"`
	BaseCurrency             string             `json:"base_currency"`
	OptionType               string             `json:"option_type"`
	Strike                   float64            `json:"strike"`
	TickSizeSteps            []TickSizeStepType `json:"tick_size_steps"`
}

type TickSizeStepType struct {
	AbovePrice float64 `json:"above_price"`
	TickSize   float64 `json:"tick_size"`
}

func NewGetInstruments(d *Deribit) *GetInstruments {
	return &GetInstruments{
		Deribit: d,
	}
}

func (g *GetInstruments) GetInstruments(validCurrencies, validInstrumentKinds []string) ([]Instrument, error) {
	instruments := []Instrument{}
	for _, currency := range validCurrencies {
		for _, kind := range validInstrumentKinds {
			var instrumentsResult GetInstrumentsResponseResult
			url := fmt.Sprintf("%s%s?currency=%s&kind=%s", g.config.API_URL_BASE, g.config.GET_INSTRUMENTS_ENDPOINT, currency, kind)
			instrumentsRequest := createNewRequest(get, url, nil, nil, &instrumentsResult)
			err := instrumentsRequest.sendHTTPRequest()
			if err != nil {
				return nil, err
			}
			instruments = append(instruments, instrumentsResult.Result...)
		}
	}
	return instruments, nil
}
