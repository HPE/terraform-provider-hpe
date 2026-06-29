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

	// Get appliance settings
	req := &morpheus.Request{}
	applianceSettingsResponse, err := client.GetApplianceSettings(req)
	if err != nil {
		log.Fatal(err)
	}
	result := applianceSettingsResponse.Result.(*morpheus.GetApplianceSettingsResult)
	log.Println(result.ApplianceSettings)
}
