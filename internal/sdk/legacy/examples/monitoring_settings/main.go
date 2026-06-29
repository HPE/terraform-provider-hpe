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

	// Get guidance settings
	req := &morpheus.Request{}
	monitoringSettingsResponse, err := client.GetMonitoringSettings(req)
	if err != nil {
		log.Fatal(err)
	}
	result := monitoringSettingsResponse.Result.(*morpheus.GetMonitoringSettingsResult)
	log.Println(result.MonitoringSettings)
}
