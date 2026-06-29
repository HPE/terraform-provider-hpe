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

	// List network proxies
	req := &morpheus.Request{}
	networkProxyResponse, err := client.ListNetworkProxies(req)
	if err != nil {
		log.Fatal(err)
	}
	result := networkProxyResponse.Result.(*morpheus.ListNetworkProxiesResult)
	test := result.NetworkProxies
	fmt.Println(test)
	log.Println(&result.NetworkProxies)
}
