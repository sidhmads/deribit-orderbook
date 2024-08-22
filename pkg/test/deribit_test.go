package test

import (
	"deribit-connector/pkg/deribit"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValidCurrenciesFromUser(t *testing.T) {
	deribitConnector, err := deribit.NewDeribit()
	if err != nil {
		t.Fatal(err)
	}

	inputOne := "btC,usDT, EurR"
	expectedListOne := []string{deribit.BtcCurrency, deribit.UsdtCurrency, deribit.EurrCurrency}
	slices.Sort(expectedListOne)
	actualListOne := deribitConnector.GetValidCurrenciesFromUser(inputOne)
	slices.Sort(actualListOne)

	assert.Equal(t, expectedListOne, actualListOne)

	inputTwo := "aNy"
	expectedListTwo := []string{deribit.BtcCurrency, deribit.EthCurrency,
		deribit.UsdcCurrency, deribit.UsdtCurrency, deribit.EurrCurrency}
	slices.Sort(expectedListTwo)
	actualListTwo := deribitConnector.GetValidCurrenciesFromUser(inputTwo)
	slices.Sort(actualListTwo)
	assert.Equal(t, expectedListTwo, actualListTwo)
}

func TestGetValidInstrumentKindFromUser(t *testing.T) {
	inputOne := "sPot, future, OPTion, future_combo,option_COMBO"
	deribitConnector, err := deribit.NewDeribit()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{deribit.SpotInstrumentKind, deribit.FutureInstrumentKind, deribit.OptionInstrumentKind, deribit.FutureCombo, deribit.OptionCombo}, deribitConnector.GetValidInstrumentKindFromUser(inputOne))
}
