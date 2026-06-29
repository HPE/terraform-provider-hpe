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

	// List integrations
	req := &morpheus.Request{
		QueryParams: map[string]string{
			"max": "50",
		},
	}
	integrationsResponse, err := client.ListIntegrations(req)
	if err != nil {
		log.Fatal(err)
	}
	result := integrationsResponse.Result.(*morpheus.ListIntegrationsResult)
	for _, integration := range *result.Integrations {
		log.Println(integration)
	}
}
