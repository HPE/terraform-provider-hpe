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

	// List report types
	req := &morpheus.Request{}
	reportTypeResponse, err := client.ListReportTypes(req)
	if err != nil {
		log.Fatal(err)
	}
	result := reportTypeResponse.Result.(*morpheus.ListReportTypesResult)
	log.Println(result.ReportTypes)

	// Find report type by name
	findByNameResponse, err := client.FindReportTypeByName("Amazon Reservation Coverage")
	if err != nil {
		log.Fatal(err)
	}
	result = findByNameResponse.Result.(*morpheus.ListReportTypesResult)
	log.Println(result.ReportTypes)
}
