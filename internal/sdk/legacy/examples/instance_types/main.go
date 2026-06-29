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

	// List instance types
	req := &morpheus.Request{}
	instanceTypesResponse, err := client.ListInstanceTypes(req)
	if err != nil {
		log.Fatal(err)
	}
	result := instanceTypesResponse.Result.(*morpheus.ListInstanceTypesResult)
	log.Println(result.InstanceTypes)
}
