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

	// List load balancer pools
	req := &morpheus.Request{}
	loadBalancerPoolResponse, err := client.ListLoadBalancerPools(1, req)
	if err != nil {
		log.Fatal(err)
	}
	result := loadBalancerPoolResponse.Result.(*morpheus.ListLoadBalancerPoolsResult)
	log.Println(&result.LoadBalancerPools)
}
