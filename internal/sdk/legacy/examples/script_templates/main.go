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

	// List script templates
	req := &morpheus.Request{}
	response, err := client.ListScriptTemplates(req)
	if err != nil {
		log.Fatal(err)
	}
	result := response.Result.(*morpheus.ListScriptTemplatesResult)
	log.Println(result.ScriptTemplates)
}
