package main

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestCreateOrder_MultipleCases(t *testing.T) {
	for i := 1; i <= 7; i++ {
		caseName := "case" + strconv.Itoa(i)
		t.Run(caseName, func(t *testing.T) {
			inputFile := filepath.Join("testdata", caseName+".json")
			expectedFile := filepath.Join("testdata", caseName+".json.expected.json")

			inputData, err := os.ReadFile(inputFile)
			if err != nil {
				t.Fatalf("Failed to read input file: %v", err)
			}

			expectedData, err := os.ReadFile(expectedFile)
			if err != nil {
				t.Fatalf("Failed to read expected output file: %v", err)
			}

			var inputOrders []Order
			err = json.Unmarshal(inputData, &inputOrders)
			if err != nil {
				t.Fatalf("Failed to unmarshal input JSON: %v", err)
			}

			actualResult := createOrder(inputOrders)
			actualOutput, err := json.MarshalIndent(actualResult, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal actual output: %v", err)
			}

			var expectedOutput, actualParsedOutput []map[string]interface{}
			if err := json.Unmarshal(expectedData, &expectedOutput); err != nil {
				t.Fatalf("Failed to unmarshal expected output: %v", err)
			}
			if err := json.Unmarshal(actualOutput, &actualParsedOutput); err != nil {
				t.Fatalf("Failed to unmarshal actual output: %v", err)
			}

			if !equalJSON(expectedOutput, actualParsedOutput) {
				expectedPretty, _ := json.MarshalIndent(expectedOutput, "", "  ")
				actualPretty, _ := json.MarshalIndent(actualParsedOutput, "", "  ")
				t.Errorf("Output mismatch\nExpected: %s\nActual: %s", expectedPretty, actualPretty)
			}
		})
	}
}

func equalJSON(a, b []map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for k, va := range a[i] {
			vb, ok := b[i][k]
			if !ok || !equalValue(va, vb) {
				return false
			}
		}
	}
	return true
}

func equalValue(a, b interface{}) bool {
	aNum, aOk := toFloat64(a)
	bNum, bOk := toFloat64(b)
	if aOk && bOk {
		return math.Abs(aNum-bNum) < 0.00001
	}
	return a == b
}

func toFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case json.Number:
		f, err := v.Float64()
		if err == nil {
			return f, true
		}
	}
	return 0, false
}
