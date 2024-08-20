package main

import (
	"context"
	"deribit-connector/pkg/deribit"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	execute()
}

func execute() {
	var rootCMD = &cobra.Command{
		Use:   "connector",
		Short: "Command line Interface for connector features",
	}

	rootCMD.AddCommand(&cobra.Command{
		Use:   "orderbook",
		Short: "Starts orderbook feed",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			connector, err := deribit.NewDeribit()
			if err != nil {
				panic(err)
			}
			validCurrencies := connector.GetValidCurrenciesFromUser(args[0])
			validInstrumentKinds := connector.GetValidInstrumentKindFromUser(args[1])

			if len(validCurrencies) == 0 || len(validInstrumentKinds) == 0 {
				panic(fmt.Errorf("invalid user inputs for either or both valid currencies and valid instrument kinds"))
			}
			getInstruments := deribit.NewGetInstruments(connector)
			instruments, err := getInstruments.GetInstruments(validCurrencies, validInstrumentKinds)
			if err != nil {
				panic(err)
			}

			orderbook, err := deribit.NewOrderbook(connector, instruments)
			if err != nil {
				panic(err)
			}
			return orderbook.StreamOrderbooks(context.Background())
		},
	})

	rootCMD.AddCommand(&cobra.Command{
		Use:   "orderbook-consumer",
		Short: "Starts orderbook consumer",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			connector, err := deribit.NewDeribit()
			if err != nil {
				panic(err)
			}
			validCurrencies := connector.GetValidCurrenciesFromUser(args[0])
			validInstrumentKinds := connector.GetValidInstrumentKindFromUser(args[1])

			if len(validCurrencies) == 0 || len(validInstrumentKinds) == 0 {
				panic(fmt.Errorf("invalid user inputs for either or both valid currencies and valid instrument kinds"))
			}
			getConsumer := deribit.NewOrderbookConsumer(connector)
			getConsumer.StartConsuming(context.Background(), validCurrencies, validInstrumentKinds)
			return nil
		},
	})

	if err := rootCMD.Execute(); err != nil {
		zap.Error(err)
		os.Exit(1)
	}
}
