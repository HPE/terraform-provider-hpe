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

	// Get archives
	req := &morpheus.Request{}
	archivesResponse, err := client.ListArchives(req)
	if err != nil {
		log.Fatal(err)
	}
	result := archivesResponse.Result.(*morpheus.ListArchivesResult)
	log.Println(result)
}
