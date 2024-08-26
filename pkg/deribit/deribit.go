package deribit

const (
	SpotInstrumentKind   string = "spot"
	FutureInstrumentKind string = "future"
	OptionInstrumentKind string = "option"
	FutureCombo          string = "future_combo"
	OptionCombo          string = "option_combo"

	BtcCurrency  string = "BTC"
	EthCurrency  string = "ETH"
	UsdcCurrency string = "USDC"
	UsdtCurrency string = "USDT"
	EurrCurrency string = "EURR"
	AnyCurrency  string = "ANY"

	orderbookTopic string = "deribit-orderbook"
)

var (
	validInstrumentKindMappings = map[string]string{SpotInstrumentKind: SpotInstrumentKind, FutureInstrumentKind: FutureInstrumentKind,
		OptionInstrumentKind: OptionInstrumentKind, FutureCombo: FutureCombo, OptionCombo: OptionCombo}

	validCurrencyMappings = map[string]string{BtcCurrency: BtcCurrency, EthCurrency: EthCurrency, UsdcCurrency: UsdcCurrency,
		UsdtCurrency: UsdtCurrency, EurrCurrency: EurrCurrency, AnyCurrency: "any"}
)

type Deribit struct {
	config *DeribitConfig
}

func NewDeribit() (*Deribit, error) {
	config := &DeribitConfig{}
	err := config.readFromEnv()
	if err != nil {
		return nil, err
	}

	return &Deribit{
		config: config,
	}, nil
}

func (d *Deribit) GetValidCurrenciesFromUser(currencies string) []string {
	parsedCurrencies := splitAndTrim(currencies, ToUpper)
	for _, val := range parsedCurrencies {
		if val == AnyCurrency {
			validOptions := []string{}
			for key, currency := range validCurrencyMappings {
				if key != AnyCurrency {
					validOptions = append(validOptions, currency)
				}
			}
			return validOptions
		}
	}
	return d.getValidOptions(parsedCurrencies, validCurrencyMappings)
}

func (d *Deribit) GetValidInstrumentKindFromUser(instrumentTypes string) []string {
	parsedInstrumentTypes := splitAndTrim(instrumentTypes, ToLower)
	return d.getValidOptions(parsedInstrumentTypes, validInstrumentKindMappings)
}

func (d *Deribit) getValidOptions(arr []string, mappings map[string]string) []string {
	validOptions := []string{}
	for _, val := range arr {
		if option, ok := mappings[val]; ok {
			validOptions = append(validOptions, option)
		}
	}
	return validOptions
}
