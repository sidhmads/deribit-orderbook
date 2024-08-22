package test

import (
	"deribit-connector/pkg/deribit"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitToBatches(t *testing.T) {
	testOne := []int{}
	expectedOutput := [][]int{}
	temp := []int{}
	for i := 1; i < 500; i++ {
		testOne = append(testOne, i)
		temp = append(temp, i)
		if i%10 == 0 {
			expectedOutput = append(expectedOutput, temp)
			temp = make([]int, 0)
		}
	}
	if len(temp) != 0 {
		expectedOutput = append(expectedOutput, temp)
	}
	assert.Equal(t, expectedOutput, deribit.SplitToBatches(testOne, 10))
}
