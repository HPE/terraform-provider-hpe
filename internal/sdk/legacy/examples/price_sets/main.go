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

	// List price sets
	req := &morpheus.Request{}
	priceSetResponse, err := client.ListPriceSets(req)
	if err != nil {
		log.Fatal(err)
	}
	result := priceSetResponse.Result.(*morpheus.ListPriceSetsResult)
	log.Println(result.PriceSets)
}
