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

	// List network pool servers
	req := &morpheus.Request{}
	networkPoolServersResponse, err := client.ListNetworkPoolServers(req)
	if err != nil {
		log.Fatal(err)
	}
	result := networkPoolServersResponse.Result.(*morpheus.ListNetworkPoolServersResult)
	poolServers := result.NetworkPoolServers
	log.Println(poolServers)
}
