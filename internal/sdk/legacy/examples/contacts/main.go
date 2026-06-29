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

	// List contacts
	req := &morpheus.Request{}
	response, err := client.ListContacts(req)
	if err != nil {
		log.Fatal(err)
	}
	result := response.Result.(*morpheus.ListContactsResult)
	log.Println(result.Contacts)
}
