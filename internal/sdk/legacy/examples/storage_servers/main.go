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

	// List storage servers
	req := &morpheus.Request{}
	response, err := client.ListStorageServers(req)
	if err != nil {
		log.Fatal(err)
	}
	result := response.Result.(*morpheus.ListStorageServersResult)
	log.Println(result.StorageServers)
}
