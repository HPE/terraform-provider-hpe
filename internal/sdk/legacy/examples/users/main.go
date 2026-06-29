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

	// List users
	req := &morpheus.Request{}
	listUsersResponse, err := client.ListUsers(req)
	if err != nil {
		log.Fatal(err)
	}
	result := listUsersResponse.Result.(*morpheus.ListUsersResult)
	log.Println(result)
}
