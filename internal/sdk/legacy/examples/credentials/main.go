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

	// List credentials
	req := &morpheus.Request{}
	credentialsResponse, err := client.ListCredentials(req)
	if err != nil {
		log.Fatal(err)
	}
	result := credentialsResponse.Result.(*morpheus.ListCredentialsResult)
	log.Println(result)
}
