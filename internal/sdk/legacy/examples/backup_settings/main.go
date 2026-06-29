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

	// Get backup settings
	req := &morpheus.Request{}
	backupSettingsResponse, err := client.GetBackupSettings(req)
	if err != nil {
		log.Fatal(err)
	}
	result := backupSettingsResponse.Result.(*morpheus.GetBackupSettingsResult)
	log.Println(result)
}
