package internal_test

import (
	"testing"

	"github.com/airgap-solution/crypto-wallet-rest/internal"
)

func TestTestFunc(t *testing.T) {
	tests := []struct {
		testName      string
		inputValue    string
		expectedValue bool
	}{
		{testName: "match test", inputValue: "test", expectedValue: true},
		//{testName: "non-match", inputValue: "hello", expectedValue: false},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result := internal.TestFunc(tt.inputValue)
			if result != tt.expectedValue {
				t.Errorf("TestFunc(%q) = %v; want %v", tt.inputValue, result, tt.expectedValue)
			}
		})
	}
}
