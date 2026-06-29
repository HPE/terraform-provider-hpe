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

	// List power schedules
	req := &morpheus.Request{}
	response, err := client.ListPowerSchedules(req)
	if err != nil {
		log.Fatal(err)
	}
	result := response.Result.(*morpheus.ListPowerSchedulesResult)
	log.Println(result.PowerSchedules)
}
