package deribit

import (
	"fmt"
	"strings"
)

type CaseType int

const (
	ToUpper CaseType = iota
	ToLower
	Unchanged
)

func splitAndTrim(input string, caseType CaseType) []string {
	if caseType == ToUpper {
		input = strings.ToUpper(input)
	} else if caseType == ToLower {
		input = strings.ToLower(input)
	}
	parts := strings.Split(input, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	return parts
}

func SplitToBatches[K any](arr []K, sizePerPartition int) [][]K {
	var partitioned [][]K

	temp := []K{}
	for _, val := range arr {
		temp = append(temp, val)
		if len(temp) == sizePerPartition {
			partitioned = append(partitioned, temp)
			temp = make([]K, 0)
		}
	}
	if len(temp) > 0 {
		partitioned = append(partitioned, temp)
	}
	return partitioned
}

func createOrderbookTopic(currency, instrumentKind string) string {
	if strings.ToLower(instrumentKind) == "spot" {
		return "spot-orderbook"
	}
	return fmt.Sprintf("%s-%s-orderbook", strings.ToLower(currency), strings.ToLower(instrumentKind))
}