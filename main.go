package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Order struct {
	No                int     `json:"no"`
	PlatformProductId string  `json:"platformProductId"`
	Qty               int     `json:"qty"`
	UnitPrice         float64 `json:"unitPrice"`
	TotalPrice        float64 `json:"totalPrice"`
}

type CleanedOrder struct {
	No         int     `json:"no"`
	ProductId  string  `json:"productId"`
	MaterialId string  `json:"materialId,omitempty"`
	ModelId    string  `json:"modelId,omitempty"`
	Qty        int     `json:"qty"`
	UnitPrice  float64 `json:"unitPrice"`
	TotalPrice float64 `json:"totalPrice"`
}

var Orders []Order
var noItem = 1
var allItem = 0

func main() {
	var i string
	fmt.Scan(&i)
	file, err := os.Open("testdata/" + i + ".json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var info []Order
	json.Unmarshal(byteValue, &info)

	output, err := json.MarshalIndent(createOrder(info), "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(output))
}

func createOrder(input []Order) []CleanedOrder {
	allItem = 0
	noItem = 1
	var allProductId []CleanedOrder
	m := make(map[string]int)
	for _, ProductId := range input {
		allProductId = append(allProductId, parceProductCode(&ProductId, m)...)
	}
	product := &CleanedOrder{
		No:         noItem,
		ProductId:  "WIPING-CLOTH",
		Qty:        allItem,
		UnitPrice:  0,
		TotalPrice: 0,
	}
	noItem++
	allProductId = append(allProductId, *product)
	for material, qty := range m {
		product := &CleanedOrder{
			No:         noItem,
			ProductId:  material + "-CLEANNER",
			Qty:        qty,
			UnitPrice:  0,
			TotalPrice: 0,
		}
		noItem++
		allProductId = append(allProductId, *product)
	}
	return allProductId
}

func parceProductCode(order *Order, m map[string]int) []CleanedOrder {
	var output []CleanedOrder
	PPId := order.PlatformProductId
	var item []string
	for i := 0; i < len(PPId)-1; i++ {
		if PPId[i] == 'F' && PPId[i+1] == 'G' {
			temp := ""
			for j := i; j < len(PPId); j++ {
				if PPId[j] == '/' {
					break
				}
				temp += string(PPId[j])
			}
			item = append(item, temp)
			i += len(temp) - 1
		}
	}
	var sum = 0
	for _, s := range item {
		_, num := extractNumberAndTrim(s)
		sum += num
	}
	for _, s := range item {
		id, qty := extractNumberAndTrim(s)
		parts := strings.Split(id, "-")
		var modelId string
		if len(parts) > 3 {
			modelId = parts[2] + "-" + parts[3]
		} else {
			modelId = parts[2]
		}
		product := &CleanedOrder{
			No:         noItem,
			ProductId:  id,
			MaterialId: parts[0] + "-" + parts[1],
			ModelId:    modelId,
			Qty:        order.Qty * qty,
			UnitPrice:  float64(order.UnitPrice) / float64(sum),
			TotalPrice: (float64(order.UnitPrice) / float64(sum)) * float64(order.Qty*qty),
		}
		m[parts[1]] = m[parts[1]] + order.Qty*qty
		noItem++
		allItem += order.Qty * qty
		output = append(output, *product)
	}

	return output
}

func extractNumberAndTrim(s string) (string, int) {
	parts := strings.Split(s, "*")
	if len(parts) != 2 {
		return s, 1
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return parts[0], 1
	}
	return parts[0], num
}
