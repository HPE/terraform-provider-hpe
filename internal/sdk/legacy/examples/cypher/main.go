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

	createReq := &morpheus.Request{
		Body: map[string]interface{}{
			"value": "test secret",
		},
		QueryParams: map[string]string{
			"ttl":  "15m",
			"type": "string",
		},
	}
	createResp, err := client.CreateCypher("secret/gotest", createReq)
	if err != nil {
		fmt.Println(err)
	}
	result := createResp.Result.(*morpheus.CreateCypherResult)
	fmt.Println(result)

	// List cyphers
	req := &morpheus.Request{}
	cyphersResponse, err := client.ListCyphers(req)
	if err != nil {
		log.Fatal(err)
	}
	listResult := cyphersResponse.Result.(*morpheus.ListCypherResult)
	log.Println(listResult.Cyphers)

	// Get cypher
	getReq := &morpheus.Request{}
	getResp, deleteErr := client.GetCypher("secret/gotest", getReq)
	if deleteErr != nil {
		fmt.Println(deleteErr)
	}
	getResult := getResp.Result.(*morpheus.GetCypherResult)
	fmt.Println(getResult)

	// Delete cypher
	deleteReq := &morpheus.Request{}
	deleteResp, deleteErr := client.DeleteCypher("secret/gotest", deleteReq)
	if deleteErr != nil {
		fmt.Println(deleteErr)
	}
	deleteResult := deleteResp.Result.(*morpheus.DeleteCypherResult)
	fmt.Println(deleteResult)
}
