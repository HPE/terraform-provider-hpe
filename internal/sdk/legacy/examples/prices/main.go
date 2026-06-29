package main

import (
	"fmt"
	"log"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func main() {
	client := morpheus.NewClient("https://yourmorpheus.com")
	client.SetUsernameAndPassword("username", "password")
	resp, err := client.Login()
	if err != nil {
		fmt.Println("LOGIN ERROR: ", err)
	}
	fmt.Println("LOGIN RESPONSE:", resp)

	// List prices
	req := &morpheus.Request{}
	priceResponse, err := client.ListPrices(req)
	if err != nil {
		log.Fatal(err)
	}
	result := priceResponse.Result.(*morpheus.ListPricesResult)
	log.Println(result.Prices)
}
