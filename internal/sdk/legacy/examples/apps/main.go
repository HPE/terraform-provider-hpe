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

	// Get apps
	req := &morpheus.Request{}
	appsResponse, err := client.ListApps(req)
	if err != nil {
		log.Fatal(err)
	}
	result := appsResponse.Result.(*morpheus.ListAppsResult)
	log.Println(result.Apps)
}
