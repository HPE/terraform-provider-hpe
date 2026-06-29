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

	// List network routers
	req := &morpheus.Request{}
	networkRouterResponse, err := client.ListNetworkRouters(req)
	if err != nil {
		log.Fatal(err)
	}
	result := networkRouterResponse.Result.(*morpheus.ListNetworkRoutersResult)
	log.Println((*result.NetworkRouters))

	// List network router types
	networkRouterTypeResponse, err := client.ListNetworkRouterTypes(req)
	if err != nil {
		log.Fatal(err)
	}
	routeTypes := networkRouterTypeResponse.Result.(*morpheus.ListNetworkRouterTypesResult)
	log.Println((*routeTypes.NetworkRouterTypes))
}
