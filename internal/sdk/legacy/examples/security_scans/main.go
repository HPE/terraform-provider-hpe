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

	// Get activity
	req := &morpheus.Request{}
	securityScanResponse, err := client.GetSecurityScan(1, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(securityScanResponse.JsonData)
	result := securityScanResponse.Result.(*morpheus.GetSecurityScanResult)
	log.Println(result.SecurityScan)
}
