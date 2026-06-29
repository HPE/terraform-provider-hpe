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

	// List oauth clients
	req := &morpheus.Request{}
	oauthClientResponse, err := client.ListOauthClients(req)
	if err != nil {
		log.Fatal(err)
	}
	result := oauthClientResponse.Result.(*morpheus.ListOauthClientsResult)
	oauthClients := result.OauthClients
	log.Println(oauthClients)
}
