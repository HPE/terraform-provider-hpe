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

	// List hosts
	req := &morpheus.Request{}
	hostsResp, err := client.ListHosts(req)
	if err != nil {
		log.Fatal(err)
	}
	result := hostsResp.Result.(*morpheus.ListHostsResult)
	log.Println(result.Hosts)

	// List host types
	hostTypeReq := &morpheus.Request{}
	hostTypesResp, err := client.ListHostTypes(hostTypeReq)
	if err != nil {
		log.Fatal(err)
	}
	hostTypes := hostTypesResp.Result.(*morpheus.ListHostTypesResult)
	log.Println(hostTypes.HostTypes)
}
