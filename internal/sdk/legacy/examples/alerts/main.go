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

	// Get alerts
	req := &morpheus.Request{}
	alertsResponse, err := client.ListAlerts(req)
	if err != nil {
		log.Fatal(err)
	}
	result := alertsResponse.Result.(*morpheus.ListAlertsResult)
	log.Println(result)
}
